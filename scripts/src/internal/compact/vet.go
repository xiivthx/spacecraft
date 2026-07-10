package compact

import "strings"

// FilterGoVet compacts `go vet` output.
// Strips per-package header lines (# package/path), keeps only vet error lines.
type FilterGoVet struct{}

func (FilterGoVet) Apply(stdout string) string {
	if stdout == "" {
		return "ok"
	}
	lines := strings.Split(stdout, "\n")
	var errLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		// Strip # package header lines.
		if strings.HasPrefix(trimmed, "# ") {
			continue
		}
		// Keep vet error lines (vet: path/file.go:line:col: message).
		errLines = append(errLines, trimmed)
	}
	if len(errLines) == 0 {
		return "ok"
	}
	return strings.Join(errLines, "\n")
}
