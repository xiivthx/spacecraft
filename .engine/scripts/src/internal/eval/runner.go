package eval

import (
	"fmt"
	"os"
	"time"

	"spacecraft/internal/mission"
)

// Runner executes the full eval suite against mission evidence.
type Runner struct {
	Store     mission.MissionStore
	EvalsDir  string
}

// NewRunner creates a new eval runner.
func NewRunner(store mission.MissionStore, evalsDir string) *Runner {
	return &Runner{Store: store, EvalsDir: evalsDir}
}

// RunResult is the result of a full eval run.
type RunResult struct {
	EvalResult *EvalResult
	Entry      *mission.EvidenceEntry
}

// Run executes the full eval suite for a mission.
func (r *Runner) Run(missionID string) (*RunResult, error) {
	entries, err := r.Store.ReadEvidenceEntries(missionID)
	if err != nil {
		return nil, fmt.Errorf("eval: read evidence: %w", err)
	}

	filtered := filterEvalEntries(entries)

	deterministic := RunDeterministic(filtered)

	rubric, err := LoadRubric(r.EvalsDir, missionID)
	if err != nil {
		rubric = &EvalRubric{Dimensions: StdDimensions()}
	}

	scorecard, err := ScoreRubric(rubric, filtered)
	if err != nil {
		return nil, fmt.Errorf("eval: rubric scoring: %w", err)
	}

	lmJudge := RunLMJudge(filtered)

	dataset, err := LoadDataset(r.EvalsDir, missionID)
	coveredChecks := 0
	if err == nil && dataset != nil {
		coveredChecks = len(dataset.Examples)
	}

	plan, err := r.Store.LoadPlan(missionID)
	totalChecks := 0
	if err == nil && plan != nil {
		totalChecks = len(plan.Tasks)
	}

	cfg := LoadConfig(r.EvalsDir, missionID)
	coverage := Coverage(coveredChecks, totalChecks)

	result := &EvalResult{
		MissionID:         missionID,
		Deterministic:     deterministic,
		Scorecard:         *scorecard,
		LMJudge:           lmJudge,
		Coverage:          coverage,
		CoveredChecks:     coveredChecks,
		TotalChecks:       totalChecks,
		CoverageSatisfied: coverage >= cfg.CoverageThreshold,
	}

	// Marshal eval result bytes for storage.
	resultBytes, err := MarshalJSON(result)
	if err != nil {
		return nil, fmt.Errorf("eval: marshal result: %w", err)
	}

	eid, stdoutP, stderrP, err := r.Store.ReserveEvidencePath(missionID)
	if err != nil {
		return nil, fmt.Errorf("eval: reserve evidence path: %w", err)
	}

	os.WriteFile(stdoutP, resultBytes, 0644)
	_ = ensureFileExists(stderrP)

	evalType := "eval"
	entry := &mission.EvidenceEntry{
		ID:        eid,
		Type:      &evalType,
		Label:     fmt.Sprintf("eval-%s", missionID),
		Command:   "spacecraft eval " + missionID,
		ExitCode:  0,
		Stdout:    string(resultBytes),
		Stderr:    "",
		CreatedAt: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
	}

	if err := r.Store.AppendEvidence(missionID, entry); err != nil {
		return nil, fmt.Errorf("eval: append evidence: %w", err)
	}

	return &RunResult{EvalResult: result, Entry: entry}, nil
}

func ensureFileExists(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	return f.Close()
}

func filterEvalEntries(entries []mission.EvidenceEntry) []mission.EvidenceEntry {
	var filtered []mission.EvidenceEntry
	for _, e := range entries {
		if e.Type == nil || *e.Type != "eval" {
			filtered = append(filtered, e)
		}
	}
	return filtered
}
