package app

import (
	"fmt"
	"sort"
)

func FormatProfileVars(name string, vars map[string]string) string {
	if vars == nil {
		return "(empty)\n"
	}
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := fmt.Sprintf("%s:\n", name)
	for _, k := range keys {
		out += fmt.Sprintf("  %s=%q\n", k, vars[k])
	}
	return out
}
