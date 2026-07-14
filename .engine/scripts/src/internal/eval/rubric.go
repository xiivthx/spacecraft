package eval

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"spacecraft/internal/mission"
)

func resolveContent(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	if strings.Contains(s, "\n") {
		return s
	}
	if len(s) > 500 {
		return s
	}
	data, err := os.ReadFile(s)
	if err != nil {
		return s
	}
	return string(data)
}

// ScoreRubric evaluates all evidence entries against the given rubric and returns a scorecard.
func ScoreRubric(rubric *EvalRubric, entries []mission.EvidenceEntry) (*EvalScorecard, error) {
	if rubric == nil {
		return nil, fmt.Errorf("eval: rubric is nil")
	}
	if err := validateRubricDimensions(rubric); err != nil {
		return nil, err
	}

	scores := []DimensionScore{
		scoreTaskSuccess(entries),
		scoreToolUseQuality(entries),
		scoreTrajectoryCompliance(entries),
		scoreHallucination(entries),
		scoreResponseQuality(entries),
	}

	var total float64
	maxScore := 4.0 * float64(len(scores))
	for _, s := range scores {
		total += float64(s.Score) * s.Weight
	}

	avg := 0.0
	if len(scores) > 0 {
		avg = total / float64(len(scores))
	}

	summary := fmt.Sprintf("Average score: %.1f/4.0 across %d dimensions", avg, len(scores))

	return &EvalScorecard{
		Scores:  scores,
		Total:   math.Round(total*10) / 10,
		Max:     maxScore,
		Average: math.Round(avg*10) / 10,
		Summary: summary,
	}, nil
}

func validateRubricDimensions(rubric *EvalRubric) error {
	stdDims := StdDimensions()
	stdNames := make(map[string]bool)
	for _, d := range stdDims {
		stdNames[d.Name] = true
	}

	var actualNames []string
	for _, d := range rubric.Dimensions {
		actualNames = append(actualNames, d.Name)
	}

	for _, expected := range stdDims {
		found := false
		for _, d := range rubric.Dimensions {
			if d.Name == expected.Name {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("eval: rubric missing dimension %q", expected.Name)
		}
	}

	if len(rubric.Dimensions) != len(stdDims) {
		return fmt.Errorf("eval: rubric has %d dimensions, expected %d: %v", len(rubric.Dimensions), len(stdDims), actualNames)
	}

	return nil
}

func scoreTaskSuccess(entries []mission.EvidenceEntry) DimensionScore {
	if len(entries) == 0 {
		return DimensionScore{Dimension: "task_success", Score: 0, Weight: 1.0, Rationale: "no evidence to evaluate"}
	}

	passed := 0
	for _, e := range entries {
		if e.ExitCode == 0 {
			passed++
		}
	}

	ratio := float64(passed) / float64(len(entries))
	score := int(math.Round(ratio * 4))

	return DimensionScore{
		Dimension: "task_success",
		Score:     score,
		Weight:    1.0,
		Rationale: fmt.Sprintf("%d/%d evidence entries passed (exit code 0)", passed, len(entries)),
	}
}

func scoreToolUseQuality(entries []mission.EvidenceEntry) DimensionScore {
	if len(entries) == 0 {
		return DimensionScore{Dimension: "tool_use_quality", Score: 0, Weight: 1.0, Rationale: "no evidence to evaluate"}
	}

	validCmds := 0
	for _, e := range entries {
		if e.Command != "" && isRecognizedCommand(e.Command) {
			validCmds++
		}
	}

	ratio := float64(validCmds) / float64(len(entries))
	score := int(math.Round(ratio * 4))
	if score < 1 && validCmds > 0 {
		score = 1
	}

	return DimensionScore{
		Dimension: "tool_use_quality",
		Score:     score,
		Weight:    1.0,
		Rationale: fmt.Sprintf("%d/%d entries use recognized commands", validCmds, len(entries)),
	}
}

func scoreTrajectoryCompliance(entries []mission.EvidenceEntry) DimensionScore {
	if len(entries) == 0 {
		return DimensionScore{Dimension: "trajectory_compliance", Score: 0, Weight: 1.0, Rationale: "no evidence to evaluate"}
	}

	phases := extractPhases(entries)
	if len(phases) == 0 {
		return DimensionScore{Dimension: "trajectory_compliance", Score: 1, Weight: 1.0, Rationale: "could not determine workflow phases from evidence"}
	}

	score := 0
	expectedOrder := []string{"build", "test", "lint", "validate", "evidence"}
	matchCount := countPhaseMatches(phases, expectedOrder)
	if matchCount > 0 {
		ratio := float64(matchCount) / float64(len(expectedOrder))
		score = int(math.Round(ratio * 4))
		if score < 1 {
			score = 1
		}
	}

	return DimensionScore{
		Dimension: "trajectory_compliance",
		Score:     score,
		Weight:    1.0,
		Rationale: fmt.Sprintf("Phases detected: %v; %d/%d expected phases matched", phases, matchCount, len(expectedOrder)),
	}
}

func scoreHallucination(entries []mission.EvidenceEntry) DimensionScore {
	if len(entries) == 0 {
		return DimensionScore{Dimension: "hallucination", Score: 4, Weight: 0.5, Rationale: "no evidence to evaluate — default clean score"}
	}

	suspiciousKeywords := []string{"unknown", "not found", "does not exist", "invalid", "undefined", "fabricated"}
	warnCount := 0
	for _, e := range entries {
		stdout := strings.ToLower(resolveContent(e.Stdout))
		stderr := strings.ToLower(resolveContent(e.Stderr))
		for _, kw := range suspiciousKeywords {
			if strings.Contains(stdout, kw) || strings.Contains(stderr, kw) {
				warnCount++
				break
			}
		}
	}

	score := 4
	if warnCount > 0 {
		score = 4 - int(math.Min(float64(warnCount), 4))
		if score < 0 {
			score = 0
		}
	}

	rationale := "no suspicious patterns detected"
	if warnCount > 0 {
		rationale = fmt.Sprintf("%d entries contain suspicious keywords (possible hallucination indicators)", warnCount)
	}

	return DimensionScore{
		Dimension: "hallucination",
		Score:     score,
		Weight:    0.5,
		Rationale: rationale,
	}
}

func scoreResponseQuality(entries []mission.EvidenceEntry) DimensionScore {
	if len(entries) == 0 {
		return DimensionScore{Dimension: "response_quality", Score: 0, Weight: 1.0, Rationale: "no evidence to evaluate"}
	}

	minContentLen := 10
	adequateEntries := 0
	for _, e := range entries {
		stdout := resolveContent(e.Stdout)
		stderr := resolveContent(e.Stderr)
		if len(strings.TrimSpace(stdout)) >= minContentLen || len(strings.TrimSpace(stderr)) >= minContentLen {
			adequateEntries++
		}
	}

	ratio := float64(adequateEntries) / float64(len(entries))
	score := int(math.Round(ratio * 4))
	if score < 1 && adequateEntries > 0 {
		score = 1
	}

	return DimensionScore{
		Dimension: "response_quality",
		Score:     score,
		Weight:    1.0,
		Rationale: fmt.Sprintf("%d/%d entries have adequate output content (>%d chars)", adequateEntries, len(entries), minContentLen),
	}
}

func isRecognizedCommand(cmd string) bool {
	recognized := []string{"go ", "npm ", "make ", "python", "node ", "scripts/", "./spacecraft", "spacecraft "}
	for _, r := range recognized {
		if strings.Contains(strings.ToLower(cmd), strings.ToLower(r)) {
			return true
		}
	}
	return false
}

func extractPhases(entries []mission.EvidenceEntry) []string {
	seen := make(map[string]bool)
	var phases []string
	phaseMap := map[string]string{
		"build":    "build",
		"compile":  "build",
		"test":     "test",
		"lint":     "lint",
		"vet":      "lint",
		"validate": "validate",
		"evidence": "evidence",
	}

	for _, e := range entries {
		label := strings.ToLower(e.Label)
		cmd := strings.ToLower(e.Command)
		for keyword, phase := range phaseMap {
			if (strings.Contains(label, keyword) || strings.Contains(cmd, keyword)) && !seen[phase] {
				phases = append(phases, phase)
				seen[phase] = true
			}
		}
	}
	sort.Strings(phases)
	return phases
}

func countPhaseMatches(phases, expectedOrder []string) int {
	matchCount := 0
	for _, expected := range expectedOrder {
		for _, p := range phases {
			if p == expected {
				matchCount++
				break
			}
		}
	}
	return matchCount
}
