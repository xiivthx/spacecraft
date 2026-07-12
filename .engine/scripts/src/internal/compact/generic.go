package compact

import (
	"fmt"
	"strings"
)

// FilterGeneric is the fallback filter for unknown commands.
// Deduplicates consecutive identical lines and truncates output >500 lines.
// Error lines are never deduplicated or hidden.
type FilterGeneric struct{}

func (FilterGeneric) Apply(stdout string) string {
	if stdout == "" {
		return ""
	}

	// Phase 1: dedup consecutive identical lines.
	deduped := dedupLines(stdout)

	// Phase 2: truncate if needed (after dedup — count deduped lines).
	lines := strings.Split(deduped, "\n")
	if len(lines) > 500 {
		return truncateLines(lines)
	}

	return deduped
}

// dedupLines collapses consecutive identical lines.
// Lines containing error/fatal/warning patterns are never deduplicated.
func dedupLines(input string) string {
	lines := strings.Split(input, "\n")
	if len(lines) <= 1 {
		return input
	}

	var result strings.Builder
	var prev string
	count := 0

	flush := func() {
		if count == 1 {
			result.WriteString(prev)
			result.WriteByte('\n')
		} else if count > 1 {
			result.WriteString(prev)
			fmt.Fprintf(&result, " [x%d]\n", count)
		}
	}

	for _, line := range lines {
		// Never deduplicate error/likely-important lines.
		if isErrorLine(line) {
			flush()
			prev = line
			count = 1
			flush()
			prev = ""
			count = 0
			continue
		}

		if line == prev {
			count++
		} else {
			flush()
			prev = line
			count = 1
		}
	}
	flush()

	return strings.TrimRight(result.String(), "\n")
}

// truncateLines applies head+summary+tail truncation for >500 lines.
func truncateLines(lines []string) string {
	const headTail = 250

	var result strings.Builder

	// Head: first 250 lines.
	for i := 0; i < headTail && i < len(lines); i++ {
		result.WriteString(lines[i])
		result.WriteByte('\n')
	}

	// Summary separator.
	skipped := len(lines) - 2*headTail
	fmt.Fprintf(&result, "--- %d lines skipped (total: %d) ---\n", skipped, len(lines))

	// Tail: last 250 lines.
	start := len(lines) - headTail
	if start < headTail {
		start = headTail
	}
	for i := start; i < len(lines); i++ {
		result.WriteString(lines[i])
		result.WriteByte('\n')
	}

	return strings.TrimRight(result.String(), "\n")
}

// isErrorLine checks if a line contains error/fatal/panic/warning patterns.
// These lines should never be deduplicated or hidden.
func isErrorLine(line string) bool {
	lower := strings.ToLower(line)
	// High-confidence error patterns — match only as word prefixes or standalone.
	for _, pattern := range []string{
		"error:", "fatal:", "panic:", "fail:", "failed:",
		"traceback", "exception", "warning:", "critical:",
		"segmentation fault", "stack overflow", "out of memory",
	} {
		idx := strings.Index(lower, pattern)
		if idx >= 0 {
			// Verify this is a word boundary: preceded by whitespace, start of line, or start.
			if idx == 0 || lower[idx-1] == ' ' || lower[idx-1] == '\t' {
				return true
			}
			// Also match at start of parenthesized note: "(error:".
			if idx > 0 && lower[idx-1] == '(' {
				return true
			}
		}
	}
	// Lines starting with common error indicators.
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "E ") || strings.HasPrefix(trimmed, "W ") {
		return true
	}
	// Look for typical error format: file.go:line:col: message
	if strings.Contains(trimmed, ".go:") &&
		(strings.Contains(trimmed, "error") || strings.Contains(trimmed, "undefined") ||
		 strings.Contains(trimmed, "cannot") || strings.Contains(trimmed, "not enough") ||
		 strings.Contains(trimmed, "too many")) {
		return true
	}
	return false
}
