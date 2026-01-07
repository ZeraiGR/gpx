package shell

import (
	"strings"
	"testing"
)

func TestUpsertBlock_InsertIntoEmpty(t *testing.T) {
	block := RenderBlock([]string{"export GOPROXY='x'"})
	got := UpsertBlock("", block)
	if got != block {
		t.Fatalf("got:\n%q\nwant:\n%q", got, block)
	}
}

func TestUpsertBlock_AppendWhenNoMarkers(t *testing.T) {
	block := RenderBlock([]string{"export GOPROXY='x'"})
	rc := "export PATH=$PATH\n"
	got := UpsertBlock(rc, block)
	if got == rc {
		t.Fatalf("expected updated content")
	}
	// Ensure markers exist
	if !containsAll(got, BeginMarker, EndMarker) {
		t.Fatalf("expected markers in content, got:\n%s", got)
	}
}

func TestUpsertBlock_ReplaceExisting(t *testing.T) {
	old := "export PATH=$PATH\n\n" +
		BeginMarker + "\n" +
		"export GOPROXY='old'\n" +
		EndMarker + "\n" +
		"alias ll='ls -la'\n"

	block := RenderBlock([]string{"export GOPROXY='new'"})
	got := UpsertBlock(old, block)

	if !containsAll(got, "export GOPROXY='new'") {
		t.Fatalf("expected new value, got:\n%s", got)
	}
	if containsAll(got, "export GOPROXY='old'") {
		t.Fatalf("old value should be removed, got:\n%s", got)
	}
}

func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}
