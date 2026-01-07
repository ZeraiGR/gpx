package envx

import (
	"reflect"
	"testing"
)

func TestQuoteForShell(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", "''"},
		{"abc", "'abc'"},
		{"a b", "'a b'"},
		{"a'b", `'a'"'"'b'`},
		{"$HOME", "'$HOME'"},
		{"x*y", "'x*y'"},
	}

	for _, tt := range tests {
		if got := QuoteForShell(tt.in); got != tt.want {
			t.Fatalf("QuoteForShell(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestParseAssignments(t *testing.T) {
	got, err := ParseAssignments([]string{
		"goprivate=github.com/acme/*",
		"GONOSUMDB=",
	})
	if err != nil {
		t.Fatalf("ParseAssignments error: %v", err)
	}
	want := Vars{
		"GOPRIVATE": "github.com/acme/*",
		"GONOSUMDB": "",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseAssignments got %+v, want %+v", got, want)
	}
}

func TestExportLines_Sorted(t *testing.T) {
	vars := Vars{
		"GOTOOLCHAIN": "auto",
		"GOPROXY":     "https://proxy.golang.org,direct",
	}
	lines, err := vars.ExportLines()
	if err != nil {
		t.Fatalf("ExportLines error: %v", err)
	}
	// sorted by key: GOPROXY, GOTOOLCHAIN
	want := []string{
		"export GOPROXY='https://proxy.golang.org,direct'",
		"export GOTOOLCHAIN='auto'",
	}
	if !reflect.DeepEqual(lines, want) {
		t.Fatalf("ExportLines got %#v, want %#v", lines, want)
	}
}
