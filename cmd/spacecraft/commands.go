package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func hasHelpFlag(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			return true
		}
	}
	return false
}

// resolveActive returns the mission ID from .space/current, falling back to the branch.
func resolveActive(spaceDir, mid string) string {
	if cur := readCurrent(spaceDir); cur != "" && missionExists(spaceDir, cur) {
		return cur
	}
	if mid != "" && missionExists(spaceDir, mid) {
		return mid
	}
	return ""
}

func initCmd(spaceDir string) int {
	if err := os.MkdirAll(filepath.Join(spaceDir, "missions"), 0755); err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft init:", err)
		return 1
	}
	if err := os.MkdirAll(filepath.Join(spaceDir, "roadmaps"), 0755); err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft init:", err)
		return 1
	}
	fmt.Println("Spacecraft initialized at .space/")
	return 0
}

func newCmd(args []string, spaceDir string) int {
	if hasHelpFlag(args) {
		fmt.Println("Usage: spacecraft new <title>")
		return 0
	}
	title := strings.TrimSpace(strings.Join(args, " "))
	if title == "" {
		fmt.Fprintln(os.Stderr, "spacecraft new: missing mission title")
		return 1
	}

	id := newMissionID()
	dir := missionDir(spaceDir, id)
	if err := os.MkdirAll(filepath.Join(dir, "outputs"), 0755); err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft new:", err)
		return 1
	}

	now := time.Now().UTC().Format(time.RFC3339)
	m := map[string]any{
		"id":        id,
		"title":     title,
		"state":     "active",
		"branches":  []any{},
		"createdAt": now,
	}
	data, _ := json.MarshalIndent(m, "", "  ")
	if err := os.WriteFile(filepath.Join(dir, "mission.json"), append(data, '\n'), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft new:", err)
		return 1
	}
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# "+title+"\n\n## What\n\n## Why\n"), 0644)
	plan := map[string]any{"planName": "", "missionId": id, "tasks": []any{}}
	pdata, _ := json.MarshalIndent(plan, "", "  ")
	os.WriteFile(filepath.Join(dir, "plan.json"), append(pdata, '\n'), 0644)
	os.WriteFile(filepath.Join(dir, "evidence.jsonl"), []byte(""), 0644)
	writeCurrent(spaceDir, id)

	fmt.Printf("Created mission %s\n", id)
	fmt.Println("Next: /sc-run")
	return 0
}

func missionsCmd(spaceDir string) int {
	ids := listMissionIDs(spaceDir)
	if len(ids) == 0 {
		fmt.Println("No missions.")
		return 0
	}
	current := readCurrent(spaceDir)
	for i, id := range ids {
		title, state := "(untitled)", "unknown"
		if m, err := readMission(spaceDir, id); err == nil {
			if t, ok := m["title"].(string); ok && t != "" {
				title = t
			}
			if s, ok := m["state"].(string); ok && s != "" {
				state = s
			}
		}
		marker := ""
		if id == current {
			marker = " *"
		}
		fmt.Printf("%d. %s (%s) state:%s%s\n", i+1, title, id, state, marker)
	}
	return 0
}

func useCmd(args []string, spaceDir string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: spacecraft use <number|id|title>")
		return 1
	}
	id := normalizeID(args[0])
	if !missionExists(spaceDir, id) {
		fmt.Fprintf(os.Stderr, "spacecraft use: no mission matches %q\n", args[0])
		return 1
	}
	if err := writeCurrent(spaceDir, id); err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft use:", err)
		return 1
	}
	fmt.Printf("Selected mission %s\n", id)
	return 0
}

func currentCmd(spaceDir string) int {
	cur := readCurrent(spaceDir)
	if cur == "" {
		fmt.Println("No current mission. Use spacecraft new then /sc-run.")
		return 0
	}
	fmt.Println(cur)
	return 0
}

func resolveCmd(args []string, spaceDir, mid string) int {
	sel := ""
	for _, a := range args {
		if !strings.HasPrefix(a, "--") {
			sel = a
			break
		}
	}
	if sel != "" {
		id := normalizeID(sel)
		if missionExists(spaceDir, id) {
			fmt.Printf("Mission: %s\n", id)
			return 0
		}
		fmt.Fprintf(os.Stderr, "spacecraft resolve: no mission matches %q\n", sel)
		return 1
	}
	if id := resolveActive(spaceDir, mid); id != "" {
		fmt.Printf("Mission: %s\n", id)
		return 0
	}
	fmt.Fprintln(os.Stderr, "spacecraft resolve: no mission resolved")
	return 1
}

func statusCmd(spaceDir, mid string) int {
	id := resolveActive(spaceDir, mid)
	if id == "" {
		fmt.Println("No selected mission. Use spacecraft new then /sc-run.")
		return 0
	}
	m, err := readMission(spaceDir, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft status:", err)
		return 1
	}
	fmt.Printf("Mission: %s\n", id)
	if t, ok := m["title"].(string); ok {
		fmt.Printf("Title: %s\n", t)
	}
	if s, ok := m["state"].(string); ok {
		fmt.Printf("State: %s\n", s)
	}
	fmt.Printf("Evidence: %d\n", countEvidence(spaceDir, id))
	return 0
}

func flowCmd(spaceDir, mid string) int {
	id := resolveActive(spaceDir, mid)
	if id == "" {
		fmt.Println("No selected mission. Use spacecraft new then /sc-run.")
		return 0
	}
	m, err := readMission(spaceDir, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft flow:", err)
		return 1
	}
	state, _ := m["state"].(string)
	fmt.Printf("Mission: %s\n", id)
	fmt.Printf("State: %s\n", state)
	fmt.Printf("Next: %s\n", nextStep(state))
	return 0
}

func nextStep(state string) string {
	switch state {
	case "active", "planned":
		return "/sc-run"
	case "in_progress":
		return "/sc-run (continue)"
	case "ready":
		return "/sc-ship"
	case "blocked":
		return "resolve blockers"
	case "shipped":
		return "archive"
	default:
		return "/sc-run or spacecraft new"
	}
}

func countEvidence(spaceDir, id string) int {
	data, err := os.ReadFile(filepath.Join(missionDir(spaceDir, id), "evidence.jsonl"))
	if err != nil {
		return 0
	}
	n := 0
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) != "" {
			n++
		}
	}
	return n
}

func bindBranchCmd(args []string, spaceDir, cwd, mid string) int {
	if hasHelpFlag(args) {
		fmt.Println("Usage: spacecraft bind-branch [selector]")
		return 0
	}
	id := mid
	if len(args) > 0 {
		id = normalizeID(args[0])
	}
	if id == "" {
		id = resolveActive(spaceDir, mid)
	}
	if id == "" || !missionExists(spaceDir, id) {
		fmt.Fprintln(os.Stderr, "spacecraft bind-branch: no mission to bind")
		return 1
	}
	out, err := runCmd(cwd, "git", "branch", "--show-current")
	branch := strings.TrimSpace(out)
	if err != nil || branch == "" {
		fmt.Fprintln(os.Stderr, "spacecraft bind-branch: not a git worktree or no current branch")
		return 1
	}
	m, err := readMission(spaceDir, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft bind-branch:", err)
		return 1
	}
	branches, _ := m["branches"].([]any)
	found := false
	for _, b := range branches {
		if s, ok := b.(string); ok && s == branch {
			found = true
			break
		}
	}
	if !found {
		branches = append(branches, branch)
	}
	m["branches"] = branches
	data, _ := json.MarshalIndent(m, "", "  ")
	if err := os.WriteFile(filepath.Join(missionDir(spaceDir, id), "mission.json"), append(data, '\n'), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft bind-branch:", err)
		return 1
	}
	fmt.Printf("Bound branch %s to mission %s\n", branch, id)
	return 0
}

func gitInfoCmd(cwd string) int {
	out, err := runCmd(cwd, "git", "rev-parse", "--is-inside-work-tree")
	if err != nil || strings.TrimSpace(out) != "true" {
		fmt.Println("Git: not a git worktree")
		return 0
	}
	fmt.Println("Git: worktree detected")
	if root, err := runCmd(cwd, "git", "rev-parse", "--show-toplevel"); err == nil {
		fmt.Printf("Root: %s\n", strings.TrimSpace(root))
	}
	branch, _ := runCmd(cwd, "git", "branch", "--show-current")
	b := strings.TrimSpace(branch)
	if b == "" {
		b = "(detached)"
	}
	fmt.Printf("Branch: %s\n", b)
	if sha, err := runCmd(cwd, "git", "rev-parse", "HEAD"); err == nil {
		fmt.Printf("HEAD: %s\n", strings.TrimSpace(sha))
	}
	status, _ := runCmd(cwd, "git", "status", "--porcelain")
	if strings.TrimSpace(status) == "" {
		fmt.Println("Status: clean")
	} else {
		fmt.Printf("Status: dirty (%d files)\n", len(strings.Split(strings.TrimSpace(status), "\n")))
	}
	return 0
}

func gitSuggestCmd(args []string, mid string) int {
	if hasHelpFlag(args) {
		fmt.Println("Usage: spacecraft git-suggest [type] [slug]")
		return 0
	}
	branchTypes := map[string]bool{
		"feat": true, "fix": true, "docs": true, "refactor": true,
		"test": true, "build": true, "ci": true, "chore": true,
		"perf": true, "style": true,
	}
	typ := "feat"
	var slugParts []string
	if len(args) > 0 && branchTypes[strings.ToLower(args[0])] {
		typ = strings.ToLower(args[0])
		slugParts = args[1:]
	} else {
		slugParts = args
	}
	slug := slugify(strings.Join(slugParts, " "))
	if slug == "" {
		slug = "mission"
	}
	id := mid
	if id == "" {
		id = "M0000000"
	}
	fmt.Printf("Branch: %s/%s/%s\n", typ, id, slug)
	fmt.Println("Commit: Conventional Commits (feat:, fix:, chore:, docs:, test:, refactor:)")
	fmt.Println("Merge: git merge --no-ff <branch>")
	return 0
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var sb strings.Builder
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
		} else {
			sb.WriteRune('-')
		}
	}
	parts := strings.FieldsFunc(sb.String(), func(r rune) bool { return r == '-' })
	return strings.Join(parts, "-")
}

func clarifyStatusCmd(args []string, spaceDir, mid string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: spacecraft clarify-status <open|clear|deferred>")
		return 1
	}
	status := args[0]
	valid := map[string]bool{"open": true, "clear": true, "deferred": true}
	if !valid[status] {
		fmt.Fprintf(os.Stderr, "spacecraft clarify-status: invalid status %q (open|clear|deferred)\n", status)
		return 1
	}
	id := resolveActive(spaceDir, mid)
	if id == "" {
		fmt.Fprintln(os.Stderr, "spacecraft clarify-status: no active mission - pass a mission via branch or 'use'")
		return 1
	}
	path := filepath.Join(missionDir(spaceDir, id), "clarify-status")
	if err := os.WriteFile(path, []byte(status+"\n"), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft clarify-status:", err)
		return 1
	}
	fmt.Printf("Mission %s clarification: %s\n", id, status)
	return 0
}

func closeoutCmd(spaceDir, mid string) int {
	id := resolveActive(spaceDir, mid)
	if id == "" {
		fmt.Fprintln(os.Stderr, "spacecraft closeout-check: no active mission")
		return 1
	}
	m, err := readMission(spaceDir, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft closeout-check:", err)
		return 1
	}

	dir := missionDir(spaceDir, id)
	var problems []string

	for _, f := range []string{"spec.md", "plan.json", "evidence.jsonl", "review.json"} {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			problems = append(problems, "missing "+f)
		}
	}

	state, _ := m["state"].(string)
	if state != "ready" && state != "shipped" {
		problems = append(problems, "state is "+state+", expected ready or shipped")
	}

	if data, err := os.ReadFile(filepath.Join(dir, "clarify-status")); err == nil {
		if strings.TrimSpace(string(data)) == "open" {
			problems = append(problems, "clarify-status is open")
		}
	}

	if _, err := os.Stat(filepath.Join(dir, "evidence.jsonl")); err == nil {
		problems = append(problems, closeoutEvidenceProblems(filepath.Join(dir, "evidence.jsonl"))...)
	}
	if _, err := os.Stat(filepath.Join(dir, "review.json")); err == nil {
		problems = append(problems, closeoutReviewProblems(filepath.Join(dir, "review.json"))...)
	}

	// SPACECRAFT_CLOSEOUT_SKIP_CHANGELOG=1 is for unit tests in temp dirs without
	// git history only. Production never sets this.
	if os.Getenv("SPACECRAFT_CLOSEOUT_SKIP_CHANGELOG") != "1" {
		problems = append(problems, closeoutChangelogProblems(filepath.Dir(spaceDir))...)
	}

	if len(problems) > 0 {
		fmt.Printf("Closeout blocked for %s:\n", id)
		for _, p := range problems {
			fmt.Printf("- %s\n", p)
		}
		return 1
	}
	fmt.Printf("Closeout ready for %s.\n", id)
	return 0
}

func closeoutEvidenceProblems(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{"missing evidence.jsonl"}
	}
	required := []string{"label", "command", "output", "ts"}
	entries := 0
	var problems []string
	for i, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var entry map[string]any
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			problems = append(problems, fmt.Sprintf("evidence line %d not valid JSON", i+1))
			continue
		}
		for _, field := range required {
			if _, ok := entry[field]; !ok {
				problems = append(problems, fmt.Sprintf("evidence line %d missing %s", i+1, field))
			}
		}
		if !isJSONNumber(entry["exitCode"]) {
			problems = append(problems, fmt.Sprintf("evidence line %d missing exitCode (number)", i+1))
		}
		entries++
	}
	if entries == 0 {
		problems = append(problems, "no evidence captured")
	}
	return problems
}

func isJSONNumber(v any) bool {
	switch v.(type) {
	case float64, json.Number:
		return true
	default:
		return false
	}
}

func closeoutReviewProblems(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil // missing file already reported
	}
	var review map[string]any
	if err := json.Unmarshal(data, &review); err != nil {
		return []string{"review.json invalid JSON"}
	}
	var problems []string
	status, _ := review["status"].(string)
	if status != "ready" {
		problems = append(problems, fmt.Sprintf("review.json status is %q, expected \"ready\"", status))
	}
	if findings, ok := review["findings"].([]any); ok {
		for i, raw := range findings {
			f, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			sev, _ := f["severity"].(string)
			if sev == "critical" {
				problems = append(problems, fmt.Sprintf("review finding %d has severity critical", i+1))
			}
			if blocks, ok := f["blocksShip"].(bool); ok && blocks {
				problems = append(problems, fmt.Sprintf("review finding %d has blocksShip=true", i+1))
			}
		}
	}
	rr, ok := review["releaseReadiness"].(map[string]any)
	if !ok {
		problems = append(problems, "review.json releaseReadiness must be an object")
		return problems
	}
	for _, key := range []string{"changelog", "specNote"} {
		item, ok := rr[key].(map[string]any)
		if !ok {
			problems = append(problems, fmt.Sprintf("releaseReadiness.%s must be an object with status", key))
			continue
		}
		st, _ := item["status"].(string)
		if st != "ready" {
			problems = append(problems, fmt.Sprintf("releaseReadiness.%s status is %q, expected \"ready\"", key, st))
		}
	}
	return problems
}

func closeoutChangelogProblems(cwd string) []string {
	bases := []string{"main", "origin/main"}
	var lastErr error
	for _, base := range bases {
		if _, err := runCmd(cwd, "git", "rev-parse", "--verify", base); err != nil {
			lastErr = err
			continue
		}
		out, err := runCmd(cwd, "git", "log", base+"..HEAD", "--", "CHANGELOG.md")
		if err != nil {
			lastErr = err
			continue
		}
		if strings.TrimSpace(out) == "" {
			return []string{"no commits touch CHANGELOG.md since " + base}
		}
		return nil
	}
	if lastErr != nil {
		return []string{"CHANGELOG check failed: neither main nor origin/main usable (or git unavailable)"}
	}
	return []string{"no commits touch CHANGELOG.md"}
}

func archiveCmd(args []string, spaceDir, mid string) int {
	if hasHelpFlag(args) {
		fmt.Println("Usage: spacecraft archive [selector]")
		return 0
	}
	id := mid
	for _, a := range args {
		if !strings.HasPrefix(a, "--") {
			id = normalizeID(a)
			break
		}
	}
	if id == "" {
		id = resolveActive(spaceDir, mid)
	}
	if id == "" || !missionExists(spaceDir, id) {
		fmt.Fprintln(os.Stderr, "spacecraft archive: no mission to archive")
		return 1
	}
	m, err := readMission(spaceDir, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft archive:", err)
		return 1
	}
	if state, _ := m["state"].(string); state != "shipped" {
		fmt.Fprintf(os.Stderr, "spacecraft archive: mission %s state is %s; archive only shipped missions\n", id, state)
		return 1
	}
	archiveDir := filepath.Join(spaceDir, "archive")
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft archive:", err)
		return 1
	}
	if err := os.Rename(missionDir(spaceDir, id), filepath.Join(archiveDir, id)); err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft archive:", err)
		return 1
	}
	fmt.Printf("Archived mission %s\n", id)
	suggestNextAfterArchive(spaceDir, id)
	return 0
}
