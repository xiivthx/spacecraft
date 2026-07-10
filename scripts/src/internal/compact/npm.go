package compact

import "strings"

// FilterNpmTest compacts `npm test` output (Jest-style).
// Collapses passing test suites, keeps failing tests with output, preserves summary.
type FilterNpmTest struct{}

func (FilterNpmTest) Apply(stdout string) string {
	if stdout == "" {
		return "ok"
	}
	lines := strings.Split(stdout, "\n")

	var result strings.Builder
	var failLines []string
	inFail := false
	failName := ""
	hasFailures := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Detect failure suite header: "FAIL path/to/test.js"
		if strings.HasPrefix(trimmed, "FAIL ") {
			inFail = true
			hasFailures = true
			continue
		}

		// Detect passing suite: "PASS path/to/test.js" — skip entirely.
		if strings.HasPrefix(trimmed, "PASS ") {
			inFail = false
			continue
		}

		// Inside a passing suite, skip ✓ test lines.
		if !inFail && strings.HasPrefix(trimmed, "✓") {
			continue
		}

		// Inside a failing suite: capture ✕ test name.
		if inFail && strings.HasPrefix(trimmed, "✕") {
			failName = strings.TrimSpace(strings.TrimPrefix(trimmed, "✕"))
			continue
		}

		// Inside a failing suite: capture indented output lines.
		if inFail && isIndented(line) {
			failLines = append(failLines, trimmed)
			continue
		}

		// Blank line inside failure — keep tracking (don't end failure).
		if inFail && trimmed == "" {
			continue
		}

		// Non-indented, non-blank line while in fail — flush and leave fail mode.
		if inFail && trimmed != "" {
			if failName != "" {
				failLines = append(failLines, "FAIL: "+failName)
			}
			inFail = false
			failName = ""
			// Re-check this line for summary or next suite.
		}

		// Capture summary lines.
		if isSummaryLine(trimmed) {
			result.WriteString(trimmed)
			result.WriteByte('\n')
			continue
		}
	}

	// Flush any remaining failure.
	if failName != "" {
		failLines = append(failLines, "FAIL: "+failName)
	}

	// Write failure output.
	for _, fl := range failLines {
		result.WriteString(fl)
		result.WriteByte('\n')
	}

	out := strings.TrimRight(result.String(), "\n")
	if out == "" {
		if hasFailures {
			return "FAIL"
		}
		return "ok"
	}
	// If output is only summary lines (no failures), return "ok" for all-pass.
	if !hasFailures {
		summaryOnly := true
		for _, ln := range strings.Split(out, "\n") {
			if !isSummaryLine(ln) && ln != "" {
				summaryOnly = false
				break
			}
		}
		if summaryOnly {
			return "ok"
		}
	}
	return out
}

// isIndented checks if a line is indented (starts with spaces or tabs).
func isIndented(line string) bool {
	return strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "\t")
}

func isSummaryLine(line string) bool {
	return strings.HasPrefix(line, "Test Suites:") ||
		strings.HasPrefix(line, "Tests:") ||
		strings.HasPrefix(line, "Snapshots:") ||
		strings.HasPrefix(line, "Time:")
}
