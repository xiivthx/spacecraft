package compact

import (
	"strings"
)

// FilterGoTest compacts `go test` output.
// Keeps FAIL/PASS + failure output, strips OK/--- PASS lines,
// collapses stack traces into single-line summaries.
type FilterGoTest struct{}

func (FilterGoTest) Apply(stdout string) string {
	lines := strings.Split(stdout, "\n")

	type failure struct {
		name    string
		pkg     string
		elapsed string
		output  []string
	}

	var failures []failure
	var currentFail *failure
	var inFailOutput bool
	var passCount int
	var failCount int
	finalStatus := ""
	keepLine := true

	// Also capture stderr-like content that appears in stdout (race detector, panics).
	hasErrorOutput := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		keepLine = true

		// Detect test result lines.
		if strings.HasPrefix(trimmed, "--- PASS:") {
			passCount++
			inFailOutput = false
			currentFail = nil
			continue // strip PASS lines
		}
		if strings.HasPrefix(trimmed, "--- FAIL:") {
			failCount++
			inFailOutput = true
			currentFail = &failure{name: extractTestName(trimmed, "--- FAIL:")}
			failures = append(failures, *currentFail)
			continue
		}

		// Skip OK lines for individual packages.
		if strings.HasPrefix(trimmed, "ok ") {
			continue
		}

		// Capture summary lines.
		if trimmed == "FAIL" || trimmed == "PASS" {
			finalStatus = trimmed
			continue
		}
		if strings.HasPrefix(trimmed, "FAIL\t") {
			finalStatus = trimmed
			continue
		}

		// Skip === RUN lines for passing tests.
		if strings.HasPrefix(trimmed, "=== RUN ") {
			// Keep only if we're in a failure context.
			if inFailOutput && currentFail != nil {
				keepLine = false // don't add separately
			} else {
				continue
			}
		}

		// Skip coverage output.
		if strings.HasPrefix(trimmed, "coverage:") {
			continue
		}

		// Capture failure output lines.
		if inFailOutput && currentFail != nil && keepLine {
			// Find the failure in the slice and append.
			for i := range failures {
				if failures[i].name == currentFail.name {
					if trimmed != "" {
						failures[i].output = append(failures[i].output, trimmed)
					}
					break
				}
			}
		}

		// Capture panics, race detector output.
		if strings.HasPrefix(trimmed, "panic:") ||
			strings.HasPrefix(trimmed, "fatal error:") ||
			strings.Contains(trimmed, "WARNING: DATA RACE") {
			hasErrorOutput = true
		}
	}

	// Build output.
	var result strings.Builder

	// Failures first.
	for _, f := range failures {
		result.WriteString("FAIL: ")
		result.WriteString(f.name)
		result.WriteByte('\n')
		for _, out := range f.output {
			result.WriteString("  ")
			result.WriteString(out)
			result.WriteByte('\n')
		}
	}

	// Error output.
	if hasErrorOutput {
		result.WriteString("ERRORS DETECTED\n")
	}

	// Summary.
	if failCount > 0 || passCount > 0 {
		if result.Len() > 0 {
			result.WriteByte('\n')
		}
		result.WriteString("---\n")
		if finalStatus != "" {
			result.WriteString(finalStatus)
		} else if failCount > 0 {
			result.WriteString("FAIL: ")
			result.WriteString(intToStr(failCount))
			result.WriteString("/")
			result.WriteString(intToStr(passCount + failCount))
			result.WriteString(" tests")
		} else {
			result.WriteString("ok")
		}
	}

	out := result.String()
	if out == "" {
		return "ok"
	}
	return strings.TrimRight(out, "\n")
}

// FilterGoBuild compacts `go build` output.
// Strips progress lines, keeps only error output.
type FilterGoBuild struct{}

func (FilterGoBuild) Apply(stdout string) string {
	var result strings.Builder
	lines := strings.Split(stdout, "\n")
	hasContent := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		// Error lines typically contain a colon-separated path.
		// Go error format: path/file.go:line:col: message
		if isGoErrorLine(trimmed) {
			result.WriteString(trimmed)
			result.WriteByte('\n')
			hasContent = true
		}
	}

	if !hasContent {
		return "ok"
	}
	return strings.TrimRight(result.String(), "\n")
}

// isGoErrorLine returns true if the line looks like a Go compiler error.
// Format: path/file.go:line:col: message
func isGoErrorLine(line string) bool {
	// Must contain at least one colon followed by a number (line number).
	colonCount := 0
	for _, ch := range line {
		if ch == ':' {
			colonCount++
		}
	}
	if colonCount < 2 {
		return false
	}
	// Contains ".go:" — strong signal of a Go error line.
	if strings.Contains(line, ".go:") {
		return true
	}
	return false
}

func extractTestName(line, prefix string) string {
	name := strings.TrimSpace(strings.TrimPrefix(line, prefix))
	// Strip elapsed time in parentheses.
	if idx := strings.LastIndex(name, " ("); idx >= 0 {
		name = strings.TrimSpace(name[:idx])
	}
	return name
}

func intToStr(n int) string {
	if n < 0 {
		return "0"
	}
	const digits = "0123456789"
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = digits[n%10]
		n /= 10
	}
	return string(buf[i:])
}
