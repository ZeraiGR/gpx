package envx

import (
	"fmt"
	"sort"
	"strings"
)

func UnsetLines(keys []string) ([]string, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("no keys provided")
	}

	normalized := make([]string, 0, len(keys))
	for _, k := range keys {
		k = strings.ToUpper(strings.TrimSpace(k))
		if err := ValidateKey(k); err != nil {
			return nil, err
		}
		normalized = append(normalized, k)
	}

	sort.Strings(normalized)

	out := make([]string, 0, len(normalized))
	for _, k := range normalized {
		out = append(out, fmt.Sprintf("unset %s", k))
	}
	return out, nil
}
