package app

import (
	"fmt"
	"os"
	"sort"
)

type StatusRow struct {
	Key   string
	Value string
	Set   bool
}

func (a App) Status() ([]StatusRow, error) {
	cfg, err := a.LoadConfig()
	if err != nil {
		return nil, err
	}

	keysSet := map[string]struct{}{}
	for _, prof := range cfg.Profiles {
		for k := range prof {
			keysSet[k] = struct{}{}
		}
	}

	keys := make([]string, 0, len(keysSet))
	for k := range keysSet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	rows := make([]StatusRow, 0, len(keys))
	for _, k := range keys {
		v, ok := os.LookupEnv(k)
		rows = append(rows, StatusRow{Key: k, Value: v, Set: ok})
	}
	return rows, nil
}

func FormatStatus(rows []StatusRow) string {
	if len(rows) == 0 {
		return "(no variables found in profiles)"
	}
	out := ""
	for _, r := range rows {
		if !r.Set {
			out += fmt.Sprintf("%s: (not set)\n", r.Key)
		} else {
			out += fmt.Sprintf("%s: %q\n", r.Key, r.Value)
		}
	}
	return out
}
