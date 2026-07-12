package eval

import (
	"fmt"
	"strings"

	"spacecraft/internal/mission"
)

// CheckResult from compileEvidence.
type compiledEntry struct {
	entry  mission.EvidenceEntry
	passed bool
	reason string
}

// RunDeterministic runs all deterministic checks against evidence entries.
func RunDeterministic(entries []mission.EvidenceEntry) DeterministicResult {
	checks := []DeterministicCheckResult{
		checkCompile(entries),
		checkTests(entries),
		checkLint(entries),
	}

	allPassed := true
	for _, c := range checks {
		if !c.Passed {
			allPassed = false
		}
	}

	traj := summarizeTrajectory(entries)

	return DeterministicResult{
		Checks:            checks,
		AllPassed:         allPassed,
		TrajectorySummary: traj,
	}
}

func checkCompile(entries []mission.EvidenceEntry) DeterministicCheckResult {
	var compileEntries []compiledEntry
	for _, e := range entries {
		if isCompileEvidence(e) {
			compileEntries = append(compileEntries, compiledEntry{
				entry:  e,
				passed: e.ExitCode == 0,
				reason: fmt.Sprintf("exit code %d", e.ExitCode),
			})
		}
	}
	if len(compileEntries) == 0 {
		return DeterministicCheckResult{Name: "compile", Passed: true, Details: "no compile evidence found", ExitCode: 0}
	}
	result := DeterministicCheckResult{Name: "compile", Passed: true}
	for _, ce := range compileEntries {
		if !ce.passed {
			result.Passed = false
			result.ExitCode = ce.entry.ExitCode
			result.Details = fmt.Sprintf("%s (%s): %s", ce.entry.Label, ce.entry.Command, ce.reason)
		}
	}
	if result.Passed {
		result.Details = fmt.Sprintf("%d compile checks passed", len(compileEntries))
	}
	return result
}

func checkTests(entries []mission.EvidenceEntry) DeterministicCheckResult {
	var testEntries []compiledEntry
	for _, e := range entries {
		if isTestEvidence(e) {
			testEntries = append(testEntries, compiledEntry{
				entry:  e,
				passed: e.ExitCode == 0,
				reason: fmt.Sprintf("exit code %d", e.ExitCode),
			})
		}
	}
	if len(testEntries) == 0 {
		return DeterministicCheckResult{Name: "test", Passed: true, Details: "no test evidence found", ExitCode: 0}
	}
	result := DeterministicCheckResult{Name: "test", Passed: true}
	for _, te := range testEntries {
		if !te.passed {
			result.Passed = false
			result.ExitCode = te.entry.ExitCode
			result.Details = fmt.Sprintf("%s (%s): %s", te.entry.Label, te.entry.Command, te.reason)
		}
	}
	if result.Passed {
		result.Details = fmt.Sprintf("%d test checks passed", len(testEntries))
	}
	return result
}

func checkLint(entries []mission.EvidenceEntry) DeterministicCheckResult {
	var lintEntries []compiledEntry
	for _, e := range entries {
		if isLintEvidence(e) {
			lintEntries = append(lintEntries, compiledEntry{
				entry:  e,
				passed: e.ExitCode == 0,
				reason: fmt.Sprintf("exit code %d", e.ExitCode),
			})
		}
	}
	if len(lintEntries) == 0 {
		return DeterministicCheckResult{Name: "lint", Passed: true, Details: "no lint evidence found", ExitCode: 0}
	}
	result := DeterministicCheckResult{Name: "lint", Passed: true}
	for _, le := range lintEntries {
		if !le.passed {
			result.Passed = false
			result.ExitCode = le.entry.ExitCode
			result.Details = fmt.Sprintf("%s (%s): %s", le.entry.Label, le.entry.Command, le.reason)
		}
	}
	if result.Passed {
		result.Details = fmt.Sprintf("%d lint checks passed", len(lintEntries))
	}
	return result
}

func summarizeTrajectory(entries []mission.EvidenceEntry) string {
	if len(entries) == 0 {
		return "no evidence to analyze"
	}

	var cmds []string
	totalPass := 0
	totalFail := 0
	for _, e := range entries {
		cmds = append(cmds, e.Command)
		if e.ExitCode == 0 {
			totalPass++
		} else {
			totalFail++
		}
	}

	summary := fmt.Sprintf("%d evidence entries analyzed: %d passed, %d failed. Tool sequence: %s",
		len(entries), totalPass, totalFail, strings.Join(cmds, " → "))

	return summary
}

func isCompileEvidence(e mission.EvidenceEntry) bool {
	cmd := strings.ToLower(e.Command)
	label := strings.ToLower(e.Label)
	return strings.Contains(cmd, "build") ||
		strings.Contains(cmd, "compile") ||
		strings.Contains(label, "compile") ||
		strings.Contains(label, "build")
}

func isTestEvidence(e mission.EvidenceEntry) bool {
	cmd := strings.ToLower(e.Command)
	label := strings.ToLower(e.Label)
	return strings.Contains(cmd, "test") ||
		strings.Contains(label, "test")
}

func isLintEvidence(e mission.EvidenceEntry) bool {
	cmd := strings.ToLower(e.Command)
	label := strings.ToLower(e.Label)
	return strings.Contains(cmd, "lint") ||
		strings.Contains(cmd, "vet") ||
		strings.Contains(label, "lint")
}
