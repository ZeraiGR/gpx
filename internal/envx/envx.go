package envx

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

type Vars map[string]string

// ValidateKey ensures env var name looks like POSIX-ish: [A-Z_][A-Z0-9_]*
func ValidateKey(key string) error {
	if key == "" {
		return fmt.Errorf("empty key")
	}
	for i, r := range key {
		if i == 0 {
			if !(r == '_' || ('A' <= r && r <= 'Z')) {
				return fmt.Errorf("invalid key %q: must start with A-Z or _", key)
			}
		} else {
			if !(r == '_' || ('A' <= r && r <= 'Z') || ('0' <= r && r <= '9')) {
				return fmt.Errorf("invalid key %q: only A-Z, 0-9, _ allowed", key)
			}
		}
	}
	return nil
}

func (v Vars) KeysSorted() []string {
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ExportLines returns shell-safe export lines, sorted by key.
// Empty values are exported as empty string: export KEY=""
func (v Vars) ExportLines() ([]string, error) {
	keys := v.KeysSorted()
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		if err := ValidateKey(k); err != nil {
			return nil, err
		}
		out = append(out, fmt.Sprintf(`export %s=%s`, k, QuoteForShell(v[k])))
	}
	return out, nil
}

// QuoteForShell quotes value for POSIX shell using single quotes, safely.
// abc -> 'abc'
// a'b -> 'a'"'"'b'
func QuoteForShell(s string) string {
	// Single-quote strategy is robust and predictable.
	return "'" + strings.ReplaceAll(s, "'", `'"'"'`) + "'"
}

// ParseAssignments parses KEY=VALUE tokens into Vars.
// Example: ["GOPRIVATE=github.com/x/*", "GONOSUMDB="]
func ParseAssignments(tokens []string) (Vars, error) {
	vars := Vars{}
	for _, t := range tokens {
		eq := strings.IndexByte(t, '=')
		if eq <= 0 { // no '=' or empty key
			return nil, fmt.Errorf("invalid assignment %q (expected KEY=VALUE)", t)
		}
		key := t[:eq]
		val := t[eq+1:]
		// Normalize keys to uppercase; we can be strict
		key = strings.ToUpper(key)
		// Reject weird whitespace in key
		for _, r := range key {
			if unicode.IsSpace(r) {
				return nil, fmt.Errorf("invalid key %q: contains whitespace", key)
			}
		}
		if err := ValidateKey(key); err != nil {
			return nil, err
		}
		vars[key] = val
	}
	return vars, nil
}

func NormalizeKey(key string) string {
	return strings.ToUpper(strings.TrimSpace(key))
}
