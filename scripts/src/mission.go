package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func initSpacecraft(silent bool) {
	err := os.MkdirAll(MISSIONS_DIR, 0755)
	if err != nil {
		fail(err.Error())
	}
	err = ensureCurrentFile()
	if err != nil {
		fail(err.Error())
	}
	if !silent {
		fmt.Println("Spacecraft initialized at .space/")
	}
}

func readCurrentMissionId(required bool) *string {
	if !exists(CURRENT_FILE) {
		if required {
			fail("No current Spacecraft mission. Start one with /sc-start <title>.")
		}
		return nil
	}
	content, err := os.ReadFile(CURRENT_FILE)
	if err != nil {
		if required {
			fail("No current Spacecraft mission. Start one with /sc-start <title>.")
		}
		return nil
	}
	val := strings.TrimSpace(string(content))
	if val == "" {
		if required {
			fail("No current Spacecraft mission. Start one with /sc-start <title>.")
		}
		return nil
	}
	if norm := normalizeMissionId(val); norm != nil {
		return norm
	}
	return &val
}

func missionDir(id string) string {
	return filepath.Join(MISSIONS_DIR, id)
}

func currentSessionKey() *string {
	for _, env := range []string{"SPACECRAFT_SESSION", "OPENCODE_SESSION_ID", "CODEX_SESSION_ID"} {
		if val := os.Getenv(env); val != "" {
			return &val
		}
	}
	return nil
}

func sessionFilePath() *string {
	key := currentSessionKey()
	if key == nil {
		return nil
	}
	safeKey := slugify(*key)
	if len(safeKey) > 80 {
		safeKey = safeKey[:80]
	}
	if safeKey == "" {
		return nil
	}
	res := filepath.Join(SPACE_DIR, "sessions", safeKey+".current")
	return &res
}

func writeSessionMissionId(id string) *string {
	file := sessionFilePath()
	if file == nil {
		return nil
	}
	os.MkdirAll(filepath.Dir(*file), 0755)
	os.WriteFile(*file, []byte(id+"\n"), 0644)
	return file
}

func readSessionMissionId() *string {
	file := sessionFilePath()
	if file == nil || !exists(*file) {
		return nil
	}
	content, _ := os.ReadFile(*file)
	return normalizeMissionId(strings.TrimSpace(string(content)))
}

func createMission(title string) {
	if strings.TrimSpace(title) == "" {
		fail("Missing mission title.\n\n" + usage())
	}
	initSpacecraft(true)
	id := missionId()
	dir := missionDir(id)
	os.MkdirAll(dir, 0755)
	os.MkdirAll(filepath.Join(dir, "outputs"), 0755)
	os.MkdirAll(filepath.Join(dir, "design"), 0755)

	now := isoNow()
	git := gitInfo()

	var branch *string
	if git.Branch != "" {
		branch = &git.Branch
	}
	var root *string
	if git.Root != "" {
		root = &git.Root
	}
	var sha *string
	if git.Sha != "" {
		sha = &git.Sha
	}

	mission := Mission{
		ID:        id,
		Title:     strings.TrimSpace(title),
		State:     "draft",
		CreatedAt: now,
		UpdatedAt: now,
		BaseSha:   sha,
		HeadSha:   nil,
		Git: GitBlock{
			IsRepo:            git.IsRepo,
			Root:              root,
			Branch:            branch,
			BaseSha:           sha,
			DirtyAtStart:      &git.Dirty,
			DirtyFilesAtStart: git.DirtyFiles,
		},
		Artifacts: ArtifactsBlock{
			Spec:       "spec.md",
			Plan:       "plan.json",
			Evidence:   "evidence.jsonl",
			Review:     "review.md",
			ReviewJson: "review.json",
			Questions:  "questions.md",
			Decisions:  "decisions.md",
			Design:     "design/",
		},
		Clarification: ClarificationBlock{
			Status:            "open",
			BlockingQuestions: 0,
			LastQuestion:      nil,
		},
	}

	writeJson(filepath.Join(dir, "mission.json"), mission)
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Mission Spec\n\n## Goal\n\n## User-visible behavior\n\n## Non-goals\n\n## Constraints\n\n## Acceptance checks\n"), 0644)
	writeJson(filepath.Join(dir, "plan.json"), map[string]interface{}{"missionId": id, "tasks": []interface{}{}})
	os.WriteFile(filepath.Join(dir, "evidence.jsonl"), []byte(""), 0644)
	os.WriteFile(filepath.Join(dir, "review.md"), []byte("# Mission Review\n"), 0644)
	os.WriteFile(filepath.Join(dir, "questions.md"), []byte("# Clarification Questions\n\n## Open\n\n## Answered\n"), 0644)
	os.WriteFile(filepath.Join(dir, "decisions.md"), []byte("# Mission Decisions\n\n## Confirmed\n\n## Assumptions\n"), 0644)
	writeJson(filepath.Join(dir, "review.json"), map[string]interface{}{"status": "not-reviewed", "findings": []interface{}{}})
	os.WriteFile(CURRENT_FILE, []byte(id+"\n"), 0644)
	sess := writeSessionMissionId(id)

	fmt.Printf("Created Spacecraft mission %s\n", id)
	if sess != nil {
		fmt.Printf("Session: %s\n", displayPath(*sess))
	}
	if git.IsRepo {
		b := "(detached)"
		if git.Branch != "" {
			b = git.Branch
		}
		s := "(no commit)"
		if git.Sha != "" && len(git.Sha) >= 12 {
			s = git.Sha[:12]
		}
		d := ""
		if git.Dirty {
			d = fmt.Sprintf(" dirty:%d", git.DirtyFiles)
		}
		fmt.Printf("Git: %s %s%s\n", b, s, d)
	} else {
		fmt.Println("Git: not a git worktree. Use only for discovery/design/read-only work, or explicitly accept no-git implementation risk.")
	}
	fmt.Println("Next: /sc-plan")
}

func printCurrent() {
	id := readCurrentMissionId(false)
	if id == nil {
		fmt.Println("No current Spacecraft mission. Start one with /sc-start <title>.")
		return
	}
	fmt.Println(*id)
}

func missionBranchNames(m *Mission) []string {
	res := []string{}
	if m == nil {
		return res
	}
	if m.Branch != nil && *m.Branch != "" {
		res = append(res, *m.Branch)
	}
	if m.WorkBranch != nil && *m.WorkBranch != "" {
		res = append(res, *m.WorkBranch)
	}
	if m.Git.WorkBranch != nil && *m.Git.WorkBranch != "" {
		res = append(res, *m.Git.WorkBranch)
	}
	return res
}

func missionActive(m *Mission) bool {
	if m == nil {
		return false
	}
	return m.State != "shipped"
}

func listMissionRecords() []MissionRecord {
	if !exists(MISSIONS_DIR) {
		return nil
	}
	entries, _ := os.ReadDir(MISSIONS_DIR)
	var records []MissionRecord
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		id := entry.Name()
		mPath := filepath.Join(MISSIONS_DIR, id, "mission.json")
		if !exists(mPath) {
			continue
		}
		var m Mission
		if err := readJson(mPath, &m); err == nil {
			records = append(records, MissionRecord{
				ID:       id,
				Mission:  &m,
				Dir:      missionDir(id),
				Active:   missionActive(&m),
				Branches: missionBranchNames(&m),
			})
		}
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].ID < records[j].ID
	})
	return records
}

func missionDisplayRecords(records []MissionRecord) []MissionRecord {
	res := make([]MissionRecord, len(records))
	copy(res, records)
	sort.Slice(res, func(i, j int) bool {
		if res[i].Active != res[j].Active {
			return res[i].Active
		}
		return res[i].ID > res[j].ID
	})
	return res
}

func displayMissionTitle(title string) string {
	text := strings.TrimSpace(title)
	if text == "" {
		return "(untitled)"
	}
	if len(text) <= 88 {
		return text
	}
	return text[:85] + "..."
}

func findMissionRecord(records []MissionRecord, id string) *MissionRecord {
	for i := range records {
		if records[i].ID == id {
			return &records[i]
		}
	}
	return nil
}

func findMissionBySelector(records []MissionRecord, selector string, orderedRecords []MissionRecord) *MissionRecord {
	text := strings.TrimSpace(selector)
	if text == "" {
		return nil
	}
	if num, err := strconv.Atoi(text); err == nil && num > 0 && num <= len(orderedRecords) {
		return &orderedRecords[num-1]
	}
	if id := normalizeMissionId(text); id != nil {
		return findMissionRecord(records, *id)
	}
	var exact []MissionRecord
	for _, r := range records {
		if r.Mission != nil && r.Mission.Title == text {
			exact = append(exact, r)
		}
	}
	if len(exact) == 1 {
		return &exact[0]
	}
	var matches []MissionRecord
	norm := strings.ToLower(text)
	for _, r := range records {
		if r.Mission != nil && strings.Contains(strings.ToLower(r.Mission.Title), norm) {
			matches = append(matches, r)
		}
	}
	if len(matches) == 1 {
		return &matches[0]
	}
	return nil
}

func printMissions() {
	records := listMissionRecords()
	ordered := missionDisplayRecords(records)
	if len(ordered) == 0 {
		fmt.Println("No missions.")
		return
	}
	// Resolve for safety/conflict info to print signal info
	res := resolveMission("")
	signalByID := make(map[string]string)
	for _, sig := range res.Signals {
		if sig.MissionId != nil {
			signalByID[*sig.MissionId] = sig.Source
		}
	}

	for i, r := range ordered {
		num := i + 1
		title := displayMissionTitle(r.Mission.Title)
		state := r.Mission.State
		active := ""
		if !r.Active {
			active = " shipped"
		}
		branchHint := ""
		if len(r.Branches) > 0 {
			branchHint = " [" + strings.Join(r.Branches, ",") + "]"
		}
		sig := ""
		if s, ok := signalByID[r.ID]; ok {
			sig = " <" + s + ">"
		}
		fmt.Printf("%d. %s (%s) state:%s%s%s%s\n", num, title, r.ID, state, active, sig, branchHint)
	}
}

func useMission(args []string) {
	if len(args) == 0 {
		fail("Missing selector.\nUsage: spacecraft use <number|id|title>")
	}
	selector := strings.TrimSpace(strings.Join(args, " "))
	records := listMissionRecords()
	ordered := missionDisplayRecords(records)
	record := findMissionBySelector(records, selector, ordered)
	if record == nil {
		fail(fmt.Sprintf("No mission matches %q.\nRun 'spacecraft missions' to see available missions.", selector))
	}
	// Write to .space/current
	err := os.WriteFile(CURRENT_FILE, []byte(record.ID+"\n"), 0644)
	if err != nil {
		fail("Failed to write .space/current: " + err.Error())
	}
	// Write to session binding if available
	sessFile := writeSessionMissionId(record.ID)
	fmt.Printf("Selected mission %s (%s)\n", record.ID, record.Mission.Title)
	if sessFile != nil {
		fmt.Printf("Session: %s\n", displayPath(*sessFile))
	}
}

func bindBranch(args []string) {
	selector := ""
	if len(args) > 0 {
		selector = strings.TrimSpace(strings.Join(args, " "))
	}
	var record *MissionRecord
	if selector != "" {
		records := listMissionRecords()
		ordered := missionDisplayRecords(records)
		record = findMissionBySelector(records, selector, ordered)
		if record == nil {
			fail(fmt.Sprintf("No mission matches %q.", selector))
		}
	} else {
		res := requireResolvedMission("bind-branch")
		if res.Selected == nil {
			fail("No resolved mission to bind branch.")
		}
		rec := findMissionRecord(listMissionRecords(), res.Selected.ID)
		if rec == nil {
			fail("Selected mission not found on disk.")
		}
		record = rec
	}

	git := gitInfo()
	if !git.IsRepo {
		fail("Not a git worktree; cannot bind branch.")
	}
	if git.Branch == "" {
		fail("No current branch to bind.")
	}

	// Read mission.json
	var m Mission
	mPath := filepath.Join(record.Dir, "mission.json")
	if err := readJson(mPath, &m); err != nil {
		fail("Failed to read mission.json: " + err.Error())
	}

	m.WorkBranch = &git.Branch
	m.Git.WorkBranch = &git.Branch
	m.Branch = &git.Branch
	m.UpdatedAt = isoNow()

	if err := writeJson(mPath, m); err != nil {
		fail("Failed to write mission.json: " + err.Error())
	}

	fmt.Printf("Bound branch %s to mission %s\n", git.Branch, record.ID)
}

