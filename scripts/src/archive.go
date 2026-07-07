package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func readEvidenceEntries(filePath string) []EvidenceEntry {
	if !exists(filePath) {
		return nil
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}
	var entries []EvidenceEntry
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var entry EvidenceEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			entry.ID = fmt.Sprintf("invalid-line-%d", i+1)
			entry.Label = "Invalid evidence entry"
			entry.Command = ""
			entry.ExitCode = 1
			entry.Stdout = ""
			entry.Stderr = ""
			entry.CreatedAt = ""
		}
		entries = append(entries, entry)
	}
	return entries
}

func compactEvidenceEntry(entry EvidenceEntry) CompactEvidenceEntry {
	return CompactEvidenceEntry{
		ID:        entry.ID,
		Label:     entry.Label,
		Command:   entry.Command,
		ExitCode:  entry.ExitCode,
		CreatedAt: entry.CreatedAt,
	}
}

func archiveReadinessErrors(plan *Plan, review *Review, entries []EvidenceEntry) []string {
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

	var tasks []Task
	if plan != nil && plan.Tasks != nil {
		tasks = plan.Tasks
	}

	if len(tasks) == 0 {
		errors = append(errors, "plan.json has no tasks")
	}

	var incomplete []string
	for _, t := range tasks {
		if t.Status == nil || *t.Status != "completed" {
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

	blocking := blockingReviewFindings(review)
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

	if review != nil {
		errors = append(errors, releaseReadinessErrors(review.ReleaseReadiness)...)
	}

	return errors
}

func copyArchiveText(src, dst string) bool {
	if !exists(src) {
		return false
	}
	content, err := os.ReadFile(src)
	if err != nil {
		return false
	}
	os.WriteFile(dst, content, 0644)
	return true
}

func clearArchivedMissionSelection(id string) {
	var currId *string
	cBytes, err := os.ReadFile(CURRENT_FILE)
	if err == nil {
		currId = normalizeMissionId(string(cBytes))
	}
	if currId != nil && *currId == id {
		os.WriteFile(CURRENT_FILE, []byte(""), 0644)
	}
	sessionsDir := filepath.Join(SPACE_DIR, "sessions")
	if !exists(sessionsDir) {
		return
	}
	entries, _ := os.ReadDir(sessionsDir)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(sessionsDir, entry.Name())
		content, _ := os.ReadFile(path)
		sessId := normalizeMissionId(strings.TrimSpace(string(content)))
		if sessId != nil && *sessId == id {
			os.WriteFile(path, []byte(""), 0644)
		}
	}
}

func archiveMission(args []string) {
	selector := ""
	if len(args) > 0 {
		selector = strings.Join(args, " ")
	}
	var res ResolveOutput
	if selector != "" {
		res = resolveMission(selector)
	} else {
		res = requireResolvedMission("archive")
	}
	if res.Safety != "safe" || res.Selected == nil {
		fail(formatResolutionBlock(res, "archive"))
	}

	id := res.Selected.ID
	sourceDir := missionDir(id)

	var mission Mission
	errM := readJson(filepath.Join(sourceDir, "mission.json"), &mission)
	if errM != nil {
		fail(errM.Error())
	}
	if mission.State != "shipped" {
		fail(fmt.Sprintf("Archive blocked: mission %s state is %s. Archive only shipped missions.", id, mission.State))
	}

	var plan Plan
	errP := readJson(filepath.Join(sourceDir, "plan.json"), &plan)

	var review Review
	errR := readJson(filepath.Join(sourceDir, "review.json"), &review)

	entries := readEvidenceEntries(filepath.Join(sourceDir, "evidence.jsonl"))

	var pPtr *Plan
	if errP == nil {
		pPtr = &plan
	}
	var rPtr *Review
	if errR == nil {
		rPtr = &review
	}

	errs := archiveReadinessErrors(pPtr, rPtr, entries)
	if len(errs) > 0 {
		fail(fmt.Sprintf("Archive blocked for %s:\n- %s", id, strings.Join(errs, "\n- ")))
	}

	os.MkdirAll(ARCHIVE_DIR, 0755)
	archiveDir := filepath.Join(ARCHIVE_DIR, id)
	if exists(archiveDir) {
		fail(fmt.Sprintf("Archive already exists: %s", displayPath(archiveDir)))
	}
	os.MkdirAll(archiveDir, 0755)

	archivedAt := isoNow()
	var compactEvidence []CompactEvidenceEntry
	for _, e := range entries {
		compactEvidence = append(compactEvidence, compactEvidenceEntry(e))
	}

	tasks := []Task{}
	if pPtr != nil && pPtr.Tasks != nil {
		tasks = pPtr.Tasks
	}
	completedCount := 0
	for _, t := range tasks {
		if t.Status != nil && (*t.Status == "completed" || *t.Status == "done") {
			completedCount++
		}
	}

	branch := "(unknown)"
	if mission.Git.WorkBranch != nil {
		branch = *mission.Git.WorkBranch
	} else if mission.Git.Branch != nil {
		branch = *mission.Git.Branch
	}
	title := "(untitled)"
	if mission.Title != "" {
		title = mission.Title
	}
	created := "(unknown)"
	if mission.CreatedAt != "" {
		created = mission.CreatedAt
	}
	revStat := "missing"
	if rPtr != nil && rPtr.Status != nil {
		revStat = *rPtr.Status
	}

	var summaryLines []string
	summaryLines = append(summaryLines, fmt.Sprintf("# Archived Mission %s", id))
	summaryLines = append(summaryLines, "")
	summaryLines = append(summaryLines, fmt.Sprintf("Title: %s", title))
	summaryLines = append(summaryLines, fmt.Sprintf("State: %s", mission.State))
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
		exit := "?"
		exit = fmt.Sprintf("%d", e.ExitCode)
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

	os.WriteFile(filepath.Join(archiveDir, "SUMMARY.md"), []byte(strings.Join(summaryLines, "\n")), 0644)

	bSha := mission.BaseSha
	if bSha == nil {
		bSha = mission.Git.BaseSha
	}

	compactM := CompactMission{
		ID:         mission.ID,
		Title:      mission.Title,
		State:      mission.State,
		CreatedAt:  mission.CreatedAt,
		UpdatedAt:  mission.UpdatedAt,
		ArchivedAt: archivedAt,
		BaseSha:    bSha,
		HeadSha:    mission.HeadSha,
		Git: GitBlock{
			Root:   mission.Git.Root,
			Branch: mission.Git.WorkBranch,
		},
	}
	if compactM.Git.Branch == nil {
		compactM.Git.Branch = mission.Git.Branch
	}
	compactM.Git.BaseSha = mission.Git.BaseSha
	if compactM.Git.BaseSha == nil {
		compactM.Git.BaseSha = mission.BaseSha
	}

	writeJson(filepath.Join(archiveDir, "mission.json"), compactM)

	var compactTasks []CompactTask
	for _, t := range tasks {
		compactTasks = append(compactTasks, CompactTask{
			ID:       t.ID,
			Title:    t.Title,
			Status:   t.Status,
			Evidence: []string{},
		})
	}

	mID := ""
	if pPtr != nil {
		mID = pPtr.MissionId
	}
	compactP := CompactPlan{
		MissionID: mID,
		Tasks:     compactTasks,
	}
	writeJson(filepath.Join(archiveDir, "plan.json"), compactP)

	var evStr []string
	for _, e := range compactEvidence {
		b, _ := json.Marshal(e)
		evStr = append(evStr, string(b))
	}
	evOut := strings.Join(evStr, "\n")
	if len(evStr) > 0 {
		evOut += "\n"
	}
	os.WriteFile(filepath.Join(archiveDir, "evidence.jsonl"), []byte(evOut), 0644)

	if rPtr != nil {
		writeJson(filepath.Join(archiveDir, "review.json"), *rPtr)
	}

	copyArchiveText(filepath.Join(sourceDir, "review.md"), filepath.Join(archiveDir, "review.md"))
	copyArchiveText(filepath.Join(sourceDir, "spec.md"), filepath.Join(archiveDir, "spec.md"))
	copyArchiveText(filepath.Join(sourceDir, "decisions.md"), filepath.Join(archiveDir, "decisions.md"))
	copyArchiveText(filepath.Join(sourceDir, "questions.md"), filepath.Join(archiveDir, "questions.md"))

	os.RemoveAll(sourceDir)
	clearArchivedMissionSelection(id)

	fmt.Printf("Archived mission %s\n", id)
	fmt.Printf("Archive: %s\n", displayPath(archiveDir))
}
