// Package mission defines the core data model types for the Spacecraft mission system.
package mission

import "encoding/json"

// GitBlock stores git state at mission creation/update time.
type GitBlock struct {
	IsRepo            bool    `json:"isRepo"`
	Root              *string `json:"root"`
	Branch            *string `json:"branch"`
	WorkBranch        *string `json:"workBranch,omitempty"`
	WorkBranchBoundAt *string `json:"workBranchBoundAt,omitempty"`
	BaseSha           *string `json:"baseSha"`
	DirtyAtStart      *bool   `json:"dirtyAtStart"`
	DirtyFilesAtStart int     `json:"dirtyFilesAtStart"`
}

// ArtifactsBlock defines the expected file paths for mission artifacts.
type ArtifactsBlock struct {
	Spec       string `json:"spec"`
	Plan       string `json:"plan"`
	Evidence   string `json:"evidence"`
	Review     string `json:"review"`
	ReviewJson string `json:"reviewJson"`
	Questions  string `json:"questions"`
	Decisions  string `json:"decisions"`
	Design     string `json:"design"`
}

// ClarificationBlock tracks the clarification state for a mission.
type ClarificationBlock struct {
	Status            string  `json:"status"`
	BlockingQuestions int     `json:"blockingQuestions"`
	LastQuestion      *string `json:"lastQuestion"`
}

// Mission is the root data structure for a Spacecraft mission.
type Mission struct {
	ID            string             `json:"id"`
	Title         string             `json:"title"`
	Description   string             `json:"description,omitempty"`
	State         string             `json:"state"`
	CreatedAt     string             `json:"createdAt"`
	UpdatedAt     string             `json:"updatedAt"`
	BaseSha       *string            `json:"baseSha"`
	HeadSha       *string            `json:"headSha"`
	Branch        *string            `json:"branch,omitempty"`
	WorkBranch    *string            `json:"workBranch,omitempty"`
	Git           GitBlock           `json:"git"`
	Artifacts     ArtifactsBlock     `json:"artifacts"`
	Clarification ClarificationBlock `json:"clarification"`
}

// MissionRecord is an in-memory record combining mission data with filesystem state.
type MissionRecord struct {
	ID       string
	Mission  *Mission
	Dir      string
	Active   bool
	Branches []string
}

// Task represents a single task inside a Plan.
type Task struct {
	ID     *string `json:"id"`
	Title  *string `json:"title"`
	Status *string `json:"status"`
}

// Plan is the work plan for a mission, containing a list of tasks.
type Plan struct {
	MissionId string `json:"missionId"`
	Tasks     []Task `json:"tasks"`
}

// EvidenceEntry records the result of running a verification command.
type EvidenceEntry struct {
	ID        string  `json:"id"`
	Type      *string `json:"type,omitempty"`
	Label     string  `json:"label"`
	Command   string  `json:"command"`
	ExitCode  int     `json:"exitCode"`
	Stdout    string  `json:"stdout"`
	Stderr    string  `json:"stderr"`
	Compact   *string `json:"compact,omitempty"`
	CreatedAt string  `json:"createdAt"`
}

// GitInfoData summarizes git state for the current directory.
type GitInfoData struct {
	Available  bool   `json:"available"`
	IsRepo     bool   `json:"isRepo"`
	Root       string `json:"root"`
	Branch     string `json:"branch"`
	Sha        string `json:"sha"`
	Dirty      bool   `json:"dirty"`
	DirtyFiles int    `json:"dirtyFiles"`
}

// MissionInfo is a summary of a mission used in resolution output.
type MissionInfo struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	State    string   `json:"state"`
	Active   bool     `json:"active"`
	Branches []string `json:"branches"`
	Signal   *string  `json:"signal"`
}

// SignalInfo describes a single resolution signal.
type SignalInfo struct {
	Source            string   `json:"source"`
	Value             string   `json:"value"`
	ExpectedMissionId string   `json:"expectedMissionId,omitempty"`
	MissionId         *string  `json:"missionId"`
	MissionIds        []string `json:"missionIds,omitempty"`
}

// ConflictInfo describes a resolution conflict between signals.
type ConflictInfo struct {
	Type       string       `json:"type"`
	Source     string       `json:"source,omitempty"`
	MissionId  string       `json:"missionId,omitempty"`
	Value      string       `json:"value,omitempty"`
	MissionIds []string     `json:"missionIds,omitempty"`
	Signals    []SignalInfo `json:"signals,omitempty"`
}

// CandidateInfo wraps MissionInfo with a selection number.
type CandidateInfo struct {
	MissionInfo
	Number *int `json:"number"`
}

// ResolveOutput is the complete result of mission resolution.
type ResolveOutput struct {
	Selected         *MissionInfo    `json:"selected"`
	Source           *string         `json:"source"`
	Safety           string          `json:"safety"`
	Signals          []SignalInfo    `json:"signals"`
	Conflicts        []ConflictInfo  `json:"conflicts"`
	Candidates       []CandidateInfo `json:"candidates"`
	CurrentMissionId *string         `json:"currentMissionId"`
	Git              GitInfoData     `json:"git"`
}

// TasksSummary summarizes plan task completion.
type TasksSummary struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
}

// WorkflowSnapshot is the output of the workflow/status command.
type WorkflowSnapshot struct {
	MissionID        string       `json:"missionId"`
	Title            string       `json:"title"`
	State            string       `json:"state"`
	Safety           string       `json:"safety"`
	Source           string       `json:"source"`
	Next             string       `json:"next"`
	NextTask         *Task        `json:"nextTask"`
	Tasks            TasksSummary `json:"tasks"`
	EvidenceCount    int          `json:"evidenceCount"`
	Blockers         []string     `json:"blockers"`
	CheckpointPolicy string       `json:"checkpointPolicy"`
}

// ReleaseGate represents a single gate in release readiness checks.
type ReleaseGate struct {
	Status    *string `json:"status"`
	Rationale *string `json:"rationale"`
}

// UnmarshalJSON accepts both a plain string (shorthand for status) and a full
// ReleaseGate object, so that agents writing review.json can use either form.
func (rg *ReleaseGate) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		rg.Status = &s
		return nil
	}
	type alias ReleaseGate
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	*rg = ReleaseGate(a)
	return nil
}

// ReleaseReadiness aggregates all release readiness gates.
type ReleaseReadiness struct {
	Version                  *ReleaseGate `json:"version"`
	Changelog                *ReleaseGate `json:"changelog"`
	SpecNote                 *ReleaseGate `json:"specNote"`
	TagPlan                  *ReleaseGate `json:"tagPlan"`
	PostRebaseVerification   *ReleaseGate `json:"postRebaseVerification"`
	EvalCoverage             *ReleaseGate `json:"evalCoverage"`
}

// Finding is a code review finding.
type Finding struct {
	ID         *string `json:"id"`
	Summary    *string `json:"summary"`
	Severity   *string `json:"severity"`
	BlocksShip *bool   `json:"blocksShip"`
	Criterion  *string `json:"criterion,omitempty"`
}

// Review is the code review output stored in review.json.
type Review struct {
	Status           *string          `json:"status"`
	Findings         []Finding        `json:"findings"`
	ReleaseReadiness ReleaseReadiness `json:"releaseReadiness"`
}

// CompactEvidenceEntry is a slimmed evidence entry for archive.
type CompactEvidenceEntry struct {
	ID        string `json:"id"`
	Label     string `json:"label,omitempty"`
	Command   string `json:"command,omitempty"`
	ExitCode  int    `json:"exitCode"`
	CreatedAt string `json:"createdAt"`
}

// CompactMission is a slimmed mission for archive.
type CompactMission struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	State      string   `json:"state"`
	CreatedAt  string   `json:"createdAt"`
	UpdatedAt  string   `json:"updatedAt"`
	ArchivedAt string   `json:"archivedAt"`
	BaseSha    *string  `json:"baseSha"`
	HeadSha    *string  `json:"headSha"`
	Git        GitBlock `json:"git"`
}

// CompactTask is a task with evidence references for archive.
type CompactTask struct {
	ID       *string  `json:"id"`
	Title    *string  `json:"title"`
	Status   *string  `json:"status"`
	Evidence []string `json:"evidence"`
}

// CompactPlan is the archive version of a Plan.
type CompactPlan struct {
	MissionID string        `json:"missionId"`
	Tasks     []CompactTask `json:"tasks"`
}

// TaskIsComplete returns true if the task status is a terminal/closed state.
func TaskIsComplete(status *string) bool {
	if status == nil {
		return false
	}
	switch *status {
	case "done", "cancelled":
		return true
	}
	return false
}

// BlockingFindings returns review findings that block shipping.
func BlockingFindings(review *Review) []Finding {
	if review == nil {
		return nil
	}
	var blocking []Finding
	for _, f := range review.Findings {
		blocks := f.BlocksShip != nil && *f.BlocksShip
		critical := f.Severity != nil && *f.Severity == "critical"
		if blocks || critical {
			blocking = append(blocking, f)
		}
	}
	return blocking
}
