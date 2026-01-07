package shell

import (
	"strings"
)

const (
	BeginMarker = "# GPX_BEGIN"
	EndMarker   = "# GPX_END"
)

// RenderBlock builds the managed block content.
// It ALWAYS ends with a trailing newline.
func RenderBlock(lines []string) string {
	var b strings.Builder
	b.WriteString(BeginMarker)
	b.WriteString("\n")
	for _, ln := range lines {
		b.WriteString(ln)
		b.WriteString("\n")
	}
	b.WriteString(EndMarker)
	b.WriteString("\n")
	return b.String()
}

// UpsertBlock inserts or replaces the GPX block inside rc file content.
// - If both markers exist: replace everything between them (inclusive).
// - If no markers: append block at the end, separated by a newline if needed.
func UpsertBlock(rcContent string, block string) string {
	begin := strings.Index(rcContent, BeginMarker)
	end := strings.Index(rcContent, EndMarker)

	if begin != -1 && end != -1 && end >= begin {
		// include EndMarker line
		endLine := end + len(EndMarker)
		// extend to end-of-line if present
		if endLine < len(rcContent) && rcContent[endLine] == '\r' {
			endLine++
		}
		if endLine < len(rcContent) && rcContent[endLine] == '\n' {
			endLine++
		}
		return rcContent[:begin] + block + rcContent[endLine:]
	}

	trimmed := strings.TrimRight(rcContent, "\r\n")
	if trimmed == "" {
		return block
	}
	return trimmed + "\n\n" + block
}
