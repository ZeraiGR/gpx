package app

import (
	"fmt"
	"os"
	"sort"
)

type DiffRow struct {
	Key     string
	Current string
	Target  string
	HasCurr bool
	Changed bool
}

func (a App) DiffProfile(profile string) ([]DiffRow, error) {
	cfg, err := a.LoadConfig()
	if err != nil {
		return nil, err
	}
	p, ok := cfg.Profiles[profile]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrProfileNotFound, profile)
	}

	keys := make([]string, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	rows := make([]DiffRow, 0, len(keys))
	for _, k := range keys {
		cur, has := os.LookupEnv(k)
		tgt := p[k]
		rows = append(rows, DiffRow{
			Key:     k,
			Current: cur,
			Target:  tgt,
			HasCurr: has,
			Changed: (!has && tgt != "") || (has && cur != tgt),
		})
	}
	return rows, nil
}

func FormatDiff(rows []DiffRow) string {
	if len(rows) == 0 {
		return "(no variables in profile)"
	}
	out := ""
	for _, r := range rows {
		cur := "(not set)"
		if r.HasCurr {
			cur = fmt.Sprintf("%q", r.Current)
		}
		flag := " "
		if r.Changed {
			flag = "*"
		}
		out += fmt.Sprintf("%s %s: %s -> %q\n", flag, r.Key, cur, r.Target)
	}
	out += "\nLegend: * = would change\n"
	return out
}
