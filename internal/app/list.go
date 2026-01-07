package app

import (
	"fmt"
	"sort"

	"github.com/ZeraiGR/gpx/internal/state"
)

type ProfileItem struct {
	Name   string
	Active bool
}

func (a App) ListProfiles() ([]ProfileItem, error) {
	cfg, err := a.LoadConfig()
	if err != nil {
		return nil, err
	}

	st, _ := state.Load() // best-effort
	active := ""
	if st != nil {
		active = st.ActiveProfile
	}

	names := make([]string, 0, len(cfg.Profiles))
	for name := range cfg.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)

	out := make([]ProfileItem, 0, len(names))
	for _, n := range names {
		out = append(out, ProfileItem{Name: n, Active: n == active})
	}
	return out, nil
}

func FormatProfiles(items []ProfileItem) string {
	if len(items) == 0 {
		return "(no profiles)\n"
	}
	out := ""
	for _, it := range items {
		marker := " "
		if it.Active {
			marker = "*"
		}
		out += fmt.Sprintf("%s %s\n", marker, it.Name)
	}
	out += "\nLegend: * = active (last used/applied)\n"
	return out
}
