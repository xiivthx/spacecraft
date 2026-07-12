// Package eval provides agent eval framework types and engines.
package eval

// EvalRubric defines scoring criteria for each of the 5 rubric dimensions.
type EvalRubric struct {
	Dimensions []RubricDimension `json:"dimensions"`
}

// RubricDimension describes a single scoring axis (0-4).
type RubricDimension struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// EvalDatasetEntry is a single labelled eval example.
type EvalDatasetEntry struct {
	ID            string   `json:"id"`
	Label         string   `json:"label"`
	EvidenceRefs  []string `json:"evidenceRefs"`
	ExpectedPass  bool     `json:"expectedPass"`
	ExpectedScore *int     `json:"expectedScore,omitempty"`
}

// EvalDataset is a collection of labelled eval examples.
type EvalDataset struct {
	Examples []EvalDatasetEntry `json:"examples"`
}

// EvalConfig holds per-mission eval configuration.
type EvalConfig struct {
	CoverageThreshold float64 `json:"coverageThreshold"`
}

// DefaultCoverageThreshold is the threshold used when no config.json exists.
const DefaultCoverageThreshold = 0.8

// DeterministicCheckResult is the outcome of a single deterministic check.
type DeterministicCheckResult struct {
	Name    string `json:"name"`
	Passed  bool   `json:"passed"`
	Details string `json:"details,omitempty"`
	ExitCode int   `json:"exitCode"`
}

// DeterministicResult aggregates all deterministic check results.
type DeterministicResult struct {
	Checks         []DeterministicCheckResult `json:"checks"`
	AllPassed      bool                       `json:"allPassed"`
	TrajectorySummary string                  `json:"trajectorySummary,omitempty"`
}

// DimensionScore is a single rubric dimension scored 0-4.
type DimensionScore struct {
	Dimension string  `json:"dimension"`
	Score     int     `json:"score"`
	Weight    float64 `json:"weight"`
	Rationale string  `json:"rationale,omitempty"`
}

// EvalScorecard aggregates rubric scores across all dimensions.
type EvalScorecard struct {
	Scores  []DimensionScore `json:"scores"`
	Total   float64          `json:"total"`
	Max     float64          `json:"max"`
	Average float64          `json:"average"`
	Summary string           `json:"summary,omitempty"`
}

// LMJudgeResult is the output of a secondary model evaluation.
type LMJudgeResult struct {
	Score     int    `json:"score"`
	MaxScore  int    `json:"maxScore"`
	Reasoning string `json:"reasoning,omitempty"`
	Model     string `json:"model,omitempty"`
	Fallback  bool   `json:"fallback"`
	FallbackReason string `json:"fallbackReason,omitempty"`
}

// EvalResult is the top-level eval result written to evidence.jsonl.
type EvalResult struct {
	MissionID         string              `json:"missionId"`
	Deterministic     DeterministicResult `json:"deterministic"`
	Scorecard         EvalScorecard       `json:"scorecard"`
	LMJudge           LMJudgeResult       `json:"lmJudge"`
	Coverage          float64             `json:"coverage"`
	CoveredChecks     int                 `json:"coveredChecks"`
	TotalChecks       int                 `json:"totalChecks"`
	CoverageSatisfied bool                `json:"coverageSatisfied"`
}

// StdDimensions returns the standard 5 rubric dimensions.
func StdDimensions() []RubricDimension {
	return []RubricDimension{
		{Name: "task_success", Description: "Did the agent accomplish stated task objectives?"},
		{Name: "tool_use_quality", Description: "Were tool calls correct, efficient, and free of errors?"},
		{Name: "trajectory_compliance", Description: "Did the agent follow the expected workflow order?"},
		{Name: "hallucination", Description: "Rate of fabricated tool names, file paths, dependencies, or content"},
		{Name: "response_quality", Description: "Clarity, correctness, and completeness of final output"},
	}
}

// DefaultEvalConfig returns an EvalConfig with reasonable defaults.
func DefaultEvalConfig() EvalConfig {
	return EvalConfig{
		CoverageThreshold: DefaultCoverageThreshold,
	}
}
