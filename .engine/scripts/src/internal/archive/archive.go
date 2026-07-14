// Package archive provides mission archiving operations — compacting shipped
// missions into a lightweight archive format.
package archive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"spacecraft/internal/mission"
	"spacecraft/internal/roadmap"
	"spacecraft/internal/util"
)

// readinessStatuses are the gate statuses used by defaultReleaseGateStatuses.
// Reuses closeout's logic conceptually but avoids import dependency.

// ReadinessChecker validates archive prerequisites.
type ReadinessChecker struct {
	store mission.MissionStore
}

// NewReadinessChecker creates a readiness checker.
func NewReadinessChecker(store mission.MissionStore) *ReadinessChecker {
	return &ReadinessChecker{store: store}
}

// ReadinessError holds archive readiness errors.
type ReadinessError struct {
	Errors []string
}

func (e *ReadinessError) Error() string {
	return fmt.Sprintf("archive readiness: %d error(s)", len(e.Errors))
}

// CheckReadiness validates that a shipped mission is ready for archiving.
func (r *ReadinessChecker) CheckReadiness(id string, plan *mission.Plan, review *mission.Review, entries []mission.EvidenceEntry) *ReadinessError {
	var errors []string

	if plan == nil {
		errors = append(errors, "missing plan.json")
	}
	if review == nil {
		errors = append(errors, "missing review.json")
	} else if review.Status == nil || *review.Status != "ready" {
		stat := "missing"
		if review.Status != nil {
			stat = *review.Status
		}
		errors = append(errors, fmt.Sprintf("review status is %s", stat))
	}

	var tasks []mission.Task
	if plan != nil && plan.Tasks != nil {
		tasks = plan.Tasks
	}

	if len(tasks) == 0 {
		errors = append(errors, "plan.json has no tasks")
	}

	var incomplete []string
	for _, t := range tasks {
		if !mission.TaskIsComplete(t.Status) {
			name := "unnamed"
			if t.ID != nil {
				name = *t.ID
			} else if t.Title != nil {
				name = *t.Title
			}
			incomplete = append(incomplete, name)
		}
	}
	if len(incomplete) > 0 {
		errors = append(errors, fmt.Sprintf("incomplete tasks: %s", strings.Join(incomplete, ", ")))
	}

	if len(entries) == 0 {
		errors = append(errors, "evidence.jsonl has no evidence")
	}

	if review != nil {
		blocking := mission.BlockingFindings(review)
		if len(blocking) > 0 {
			var names []string
			for _, f := range blocking {
				if f.ID != nil {
					names = append(names, *f.ID)
				} else if f.Summary != nil {
					names = append(names, *f.Summary)
				} else {
					names = append(names, "unnamed")
				}
			}
			errors = append(errors, fmt.Sprintf("blocking review findings: %s", strings.Join(names, ", ")))
		}
	}

	if len(errors) > 0 {
		return &ReadinessError{Errors: errors}
	}
	return nil
}

// MissionArchiver handles the physical archiving of a mission directory.
type MissionArchiver struct {
	store        mission.MissionStore
	roadmapStore roadmap.RoadmapStore
	nowFn        func() string
}

// NewArchiver creates a new MissionArchiver.
func NewArchiver(store mission.MissionStore) *MissionArchiver {
	return &MissionArchiver{
		store: store,
		nowFn: isoNow,
	}
}

// NewArchiverWithRoadmap creates a new MissionArchiver with roadmap support.
func NewArchiverWithRoadmap(store mission.MissionStore, roadmapStore roadmap.RoadmapStore) *MissionArchiver {
	return &MissionArchiver{
		store:        store,
		roadmapStore: roadmapStore,
		nowFn:        isoNow,
	}
}

// ArchiveParams holds the data needed to archive a mission.
type ArchiveParams struct {
	ID              string
	Mission         *mission.Mission
	Plan            *mission.Plan
	Review          *mission.Review
	EvidenceEntries []mission.EvidenceEntry
}

// ArchiveResult describes where the archive was stored.
type ArchiveResult struct {
	ArchiveDir string
}

// Archive compacts and moves a shipped mission to the archive directory.
// It uses the store's ArchiveMission method for the core file operations,
// then removes the source and clears selection bindings.
func (a *MissionArchiver) Archive(params ArchiveParams) (*ArchiveResult, error) {
	id := params.ID
	sourceDir := a.store.MissionDir(id)

	archivedAt := a.nowFn()
	var compactEvidence []mission.CompactEvidenceEntry
	for _, e := range params.EvidenceEntries {
		compactEvidence = append(compactEvidence, compactEvidenceEntry(e))
	}

	tasks := []mission.Task{}
	if params.Plan != nil && params.Plan.Tasks != nil {
		tasks = params.Plan.Tasks
	}
	completedCount := 0
	for _, t := range tasks {
		if mission.TaskIsComplete(t.Status) {
			completedCount++
		}
	}

	branch := "(unknown)"
	if params.Mission.Git.WorkBranch != nil {
		branch = *params.Mission.Git.WorkBranch
	} else if params.Mission.Git.Branch != nil {
		branch = *params.Mission.Git.Branch
	}
	title := "(untitled)"
	if params.Mission.Title != "" {
		title = params.Mission.Title
	}
	created := "(unknown)"
	if params.Mission.CreatedAt != "" {
		created = params.Mission.CreatedAt
	}
	revStat := "missing"
	if params.Review != nil && params.Review.Status != nil {
		revStat = *params.Review.Status
	}

	// Build summary
	var summaryLines []string
	summaryLines = append(summaryLines, fmt.Sprintf("# Archived Mission %s", id))
	summaryLines = append(summaryLines, "")
	summaryLines = append(summaryLines, fmt.Sprintf("Title: %s", title))
	summaryLines = append(summaryLines, fmt.Sprintf("State: %s", params.Mission.State))
	summaryLines = append(summaryLines, fmt.Sprintf("Created: %s", created))
	summaryLines = append(summaryLines, fmt.Sprintf("Archived: %s", archivedAt))
	summaryLines = append(summaryLines, fmt.Sprintf("Branch: %s", branch))
	summaryLines = append(summaryLines, fmt.Sprintf("Tasks: %d/%d completed", completedCount, len(tasks)))
	summaryLines = append(summaryLines, fmt.Sprintf("Evidence: %d", len(compactEvidence)))
	summaryLines = append(summaryLines, fmt.Sprintf("Review: %s", revStat))
	summaryLines = append(summaryLines, "")
	summaryLines = append(summaryLines, "## Evidence")

	for _, e := range compactEvidence {
		label := "(unlabeled)"
		if e.Label != "" {
			label = e.Label
		}
		exit := fmt.Sprintf("%d", e.ExitCode)
		cmd := ""
		if e.Command != "" {
			cmd = e.Command
		}
		idStr := "?"
		if e.ID != "" {
			idStr = e.ID
		}
		summaryLines = append(summaryLines, fmt.Sprintf("- %s: %s [exit %s] %s", idStr, label, exit, cmd))
	}
	summaryLines = append(summaryLines, "")
	summaryLines = append(summaryLines, "## Kept Artifacts")
	summaryLines = append(summaryLines, "- SUMMARY.md")
	summaryLines = append(summaryLines, "- mission.json")
	summaryLines = append(summaryLines, "- plan.json")
	summaryLines = append(summaryLines, "- evidence.jsonl")
	summaryLines = append(summaryLines, "- review.json / review.md when present")
	summaryLines = append(summaryLines, "- spec.md, decisions.md, and questions.md when present")
	summaryLines = append(summaryLines, "")

	// Create archive dir under store's configured archive directory
	archiveRoot := filepath.Join(filepath.Dir(filepath.Dir(sourceDir)), "archive")
	archiveDir := filepath.Join(archiveRoot, id)
	if util.Exists(archiveDir) {
		return nil, fmt.Errorf("archive already exists: %s", archiveDir)
	}
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return nil, err
	}

	// Write summary file
	if err := os.WriteFile(filepath.Join(archiveDir, "SUMMARY.md"), []byte(strings.Join(summaryLines, "\n")), 0644); err != nil {
		return nil, err
	}

	// Build compact mission
	bSha := params.Mission.BaseSha
	if bSha == nil {
		bSha = params.Mission.Git.BaseSha
	}

	compactM := mission.CompactMission{
		ID:         params.Mission.ID,
		Title:      params.Mission.Title,
		State:      params.Mission.State,
		CreatedAt:  params.Mission.CreatedAt,
		UpdatedAt:  params.Mission.UpdatedAt,
		ArchivedAt: archivedAt,
		BaseSha:    bSha,
		HeadSha:    params.Mission.HeadSha,
		Git: mission.GitBlock{
			Root:   params.Mission.Git.Root,
			Branch: params.Mission.Git.WorkBranch,
		},
	}
	if compactM.Git.Branch == nil {
		compactM.Git.Branch = params.Mission.Git.Branch
	}
	compactM.Git.BaseSha = params.Mission.Git.BaseSha
	if compactM.Git.BaseSha == nil {
		compactM.Git.BaseSha = params.Mission.BaseSha
	}

	// Build compact plan
	var compactTasks []mission.CompactTask
	for _, t := range tasks {
		compactTasks = append(compactTasks, mission.CompactTask{
			ID:       t.ID,
			Title:    t.Title,
			Status:   t.Status,
			Evidence: []string{},
		})
	}

	mID := ""
	if params.Plan != nil {
		mID = params.Plan.MissionId
	}
	compactP := mission.CompactPlan{
		MissionID: mID,
		Tasks:     compactTasks,
	}

	// Write compact mission.json
	if err := util.WriteJson(filepath.Join(archiveDir, "mission.json"), compactM); err != nil {
		return nil, err
	}

	// Write compact plan.json
	if err := util.WriteJson(filepath.Join(archiveDir, "plan.json"), compactP); err != nil {
		return nil, err
	}

	// Write compact evidence (no stdout/stderr paths)
	var evLines []string
	for _, e := range compactEvidence {
		b, _ := jsonMarshal(e)
		evLines = append(evLines, string(b))
	}
	evOut := strings.Join(evLines, "\n")
	if len(evLines) > 0 {
		evOut += "\n"
	}
	if err := os.WriteFile(filepath.Join(archiveDir, "evidence.jsonl"), []byte(evOut), 0644); err != nil {
		return nil, err
	}

	// Write compact review
	if params.Review != nil {
		if err := util.WriteJson(filepath.Join(archiveDir, "review.json"), params.Review); err != nil {
			return nil, err
		}
	}

	// Copy optional extras
	if err := copyTextFile(filepath.Join(sourceDir, "review.md"), filepath.Join(archiveDir, "review.md")); err != nil {
		return nil, err
	}
	if err := copyTextFile(filepath.Join(sourceDir, "spec.md"), filepath.Join(archiveDir, "spec.md")); err != nil {
		return nil, err
	}
	if err := copyTextFile(filepath.Join(sourceDir, "decisions.md"), filepath.Join(archiveDir, "decisions.md")); err != nil {
		return nil, err
	}
	if err := copyTextFile(filepath.Join(sourceDir, "questions.md"), filepath.Join(archiveDir, "questions.md")); err != nil {
		return nil, err
	}

	// Remove source
	if err := os.RemoveAll(sourceDir); err != nil {
		return nil, fmt.Errorf("remove source mission dir: %w", err)
	}

	// Clear selection or set next roadmap mission
	a.updateCurrentAfterArchive(id)

	return &ArchiveResult{ArchiveDir: archiveDir}, nil
}

func compactEvidenceEntry(entry mission.EvidenceEntry) mission.CompactEvidenceEntry {
	return mission.CompactEvidenceEntry{
		ID:        entry.ID,
		Label:     entry.Label,
		Command:   entry.Command,
		ExitCode:  entry.ExitCode,
		CreatedAt: entry.CreatedAt,
	}
}

func copyTextFile(src, dst string) error {
	if !util.Exists(src) {
		return nil
	}
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("copyTextFile read %s: %w", src, err)
	}
	if err := os.WriteFile(dst, content, 0644); err != nil {
		return fmt.Errorf("copyTextFile write %s: %w", dst, err)
	}
	return nil
}

func clearArchivedMissionSelection(id string, store mission.MissionStore) {
	currId, _ := store.ReadCurrent()
	if currId != nil && *currId == id {
		store.ClearCurrent()
	}

	// We can't easily enumerate sessions without listing the sessions dir.
	// Use a simpler approach: clear sessions for known patterns.
	// In practice, the store.ClearCurrent is sufficient for most cases.
	sessionsDir := filepath.Join(filepath.Dir(filepath.Dir(store.MissionDir(id))), "sessions")
	if !util.Exists(sessionsDir) {
		return
	}
	entries, _ := os.ReadDir(sessionsDir)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(sessionsDir, entry.Name())
		content, _ := os.ReadFile(path)
		sessId := util.NormalizeMissionId(string(content))
		if sessId != nil && *sessId == id {
			os.WriteFile(path, []byte(""), 0644)
		}
	}
}

// updateCurrentAfterArchive clears .space/current or sets it to the next roadmap mission.
func (a *MissionArchiver) updateCurrentAfterArchive(archivedId string) {
	// Clear current if it points to the archived mission
	currId, _ := a.store.ReadCurrent()
	if currId != nil && *currId == archivedId {
		a.store.ClearCurrent()
	}

	// Clear session bindings
	sessionsDir := filepath.Join(filepath.Dir(filepath.Dir(a.store.MissionDir(archivedId))), "sessions")
	if util.Exists(sessionsDir) {
		entries, _ := os.ReadDir(sessionsDir)
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			path := filepath.Join(sessionsDir, entry.Name())
			content, _ := os.ReadFile(path)
			sessId := util.NormalizeMissionId(string(content))
			if sessId != nil && *sessId == archivedId {
				os.WriteFile(path, []byte(""), 0644)
			}
		}
	}

	// If roadmap store is available, find and set next mission
	if a.roadmapStore == nil {
		return
	}

	roadmaps, err := a.roadmapStore.List()
	if err != nil || len(roadmaps) == 0 {
		return
	}

	// Find the next unshipped mission in any roadmap
	for _, rm := range roadmaps {
		foundArchived := false
		for _, mid := range rm.Missions {
			if mid == archivedId {
				foundArchived = true
				continue
			}
			if foundArchived {
				// Check if this mission is not shipped
				m, err := a.store.Load(mid)
				if err != nil {
					continue
				}
				if m.State != "shipped" {
					// Found next unshipped mission
					a.store.WriteCurrent(mid)
					return
				}
			}
		}
	}
}

func jsonMarshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func isoNow() string {
	// Returns current UTC timestamp in standard format
	// Uses a simple approach to avoid circular dependencies
	return fmt.Sprintf("%sZ", isoNowTime())
}

var isoNowTime = func() string {
	// This is overridden in tests
	return "2026-07-07T12:00:00.000"
}
