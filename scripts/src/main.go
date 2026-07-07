package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"spacecraft/internal/archive"
	"spacecraft/internal/closeout"
	"spacecraft/internal/config"
	"spacecraft/internal/gitutil"
	"spacecraft/internal/id"
	"spacecraft/internal/mission"
	"spacecraft/internal/resolver"
	"spacecraft/internal/state"
	"spacecraft/internal/workflow"
)

// CLI dependencies set during init.
var (
	cfg   *config.Config
	store *mission.FSStore
	r     *resolver.Resolver
	ss    *state.StateSetter
	ws    *workflow.Snapshot
	cc    *closeout.Checker
	arc   *archive.ReadinessChecker
	ar    *archive.MissionArchiver
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft:", err)
		os.Exit(1)
	}
	cfg, err = config.NewConfig(cwd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "spacecraft:", err)
		os.Exit(1)
	}
	store = mission.NewFSStore(cfg)
	r = resolver.New(store, gitutil.OSCommandRunner{}, nil)
	ss = state.NewSetter(store)
	ws = workflow.NewSnapshot(store)
	cc = closeout.NewChecker(store, gitutil.OSCommandRunner{})
	arc = archive.NewReadinessChecker(store)
	ar = archive.NewArchiver(store)

	if len(os.Args) < 2 {
		fmt.Print(usage())
		os.Exit(0)
	}

	command := os.Args[1]
	args := os.Args[2:]
	exitCode := 0

	switch command {
	case "init":
		exitCode = initCmd(args)
	case "new":
		exitCode = newCmd(args)
	case "current":
		exitCode = currentCmd()
	case "resolve":
		exitCode = resolveCmd(args)
	case "missions":
		exitCode = missionsCmd()
	case "use":
		exitCode = useCmd(args)
	case "bind-branch":
		exitCode = bindBranchCmd(args)
	case "status":
		exitCode = statusCmd(args)
	case "flow", "workflow":
		exitCode = workflowCmd(args)
	case "git-info":
		exitCode = gitInfoCmd()
	case "git-suggest":
		exitCode = gitSuggestCmd(args)
	case "set-state":
		exitCode = setStateCmd(args)
	case "clarify-status":
		exitCode = clarifyStatusCmd(args)
	case "evidence":
		exitCode = evidenceCmd(args)
	case "validate":
		exitCode = validateCmd()
	case "closeout-check":
		exitCode = closeoutCmd()
	case "archive":
		exitCode = archiveCmd(args)
	case "-h", "--help", "help":
		fmt.Print(usage())
	default:
		fmt.Fprintf(os.Stderr, "spacecraft: unknown command %q.\n\n%s", command, usage())
		exitCode = 1
	}
	os.Exit(exitCode)
}

func usage() string {
	return `Spacecraft local mission helper

Usage:
  spacecraft init
  spacecraft new <title>
  spacecraft current
  spacecraft resolve [selector] [--json]
  spacecraft missions
  spacecraft use <number|id|title>
  spacecraft bind-branch [selector]
  spacecraft status
  spacecraft flow [--json]
  spacecraft git-info
  spacecraft git-suggest [type] [slug]
  spacecraft set-state <state>
  spacecraft clarify-status <open|clear|deferred>
  spacecraft evidence <label> -- <command...>
  spacecraft validate
  spacecraft closeout-check
  spacecraft archive [selector]
`
}

// ---------- command implementations ----------

func initCmd(args []string) int {
	if err := os.MkdirAll(cfg.MissionsDir(), 0755); err != nil {
		return printErr(err)
	}
	if _, err := os.Stat(cfg.CurrentFile()); os.IsNotExist(err) {
		os.WriteFile(cfg.CurrentFile(), []byte(""), 0644)
	}
	fmt.Println("Spacecraft initialized at .space/")
	return 0
}

func newCmd(args []string) int {
	title := strings.TrimSpace(strings.Join(args, " "))
	if title == "" {
		return printErr("Missing mission title.\n\n" + usage())
	}

	os.MkdirAll(cfg.MissionsDir(), 0755)
	mid, err := id.MissionId()
	if err != nil {
		return printErr(err)
	}
	dir := cfg.MissionDir(mid)
	os.MkdirAll(dir, 0755)
	os.MkdirAll(filepath.Join(dir, "outputs"), 0755)
	os.MkdirAll(filepath.Join(dir, "design"), 0755)

	now := isoNow()
	git := gitutil.GitInfo(gitutil.OSCommandRunner{})

	if err := store.Create(&mission.Mission{
		ID:        mid,
		Title:     title,
		State:     "draft",
		CreatedAt: now,
		UpdatedAt: now,
		BaseSha:   strPtr(ifStr(git.Sha != "", git.Sha, "")),
		Git: mission.GitBlock{
			IsRepo:            git.IsRepo,
			Root:              strPtr(ifStr(git.Root != "", git.Root, "")),
			Branch:            strPtr(ifStr(git.Branch != "", git.Branch, "")),
			BaseSha:           strPtr(ifStr(git.Sha != "", git.Sha, "")),
			DirtyAtStart:      &git.Dirty,
			DirtyFilesAtStart: git.DirtyFiles,
		},
		Artifacts: mission.ArtifactsBlock{
			Spec: "spec.md", Plan: "plan.json", Evidence: "evidence.jsonl",
			Review: "review.md", ReviewJson: "review.json",
			Questions: "questions.md", Decisions: "decisions.md", Design: "design/",
		},
		Clarification: mission.ClarificationBlock{Status: "open"},
	}); err != nil {
		return printErr(err)
	}

	// Write stubs
	os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Mission Spec\n\n## Goal\n\n## User-visible behavior\n\n## Non-goals\n\n## Constraints\n\n## Acceptance checks\n"), 0644)
	os.WriteFile(filepath.Join(dir, "plan.json"), []byte(`{"missionId":"`+mid+`","tasks":[]}`+"\n"), 0644)
	os.WriteFile(filepath.Join(dir, "evidence.jsonl"), []byte(""), 0644)
	os.WriteFile(filepath.Join(dir, "review.md"), []byte("# Mission Review\n"), 0644)
	os.WriteFile(filepath.Join(dir, "questions.md"), []byte("# Clarification Questions\n\n## Open\n\n## Answered\n"), 0644)
	os.WriteFile(filepath.Join(dir, "decisions.md"), []byte("# Mission Decisions\n\n## Confirmed\n\n## Assumptions\n"), 0644)
	os.WriteFile(filepath.Join(dir, "review.json"), []byte(`{"status":"not-reviewed","findings":[]}`+"\n"), 0644)
	store.WriteCurrent(mid)
	sessFile := writeSessionBinding(mid)

	fmt.Printf("Created Spacecraft mission %s\n", mid)
	if sessFile != "" {
		fmt.Printf("Session: %s\n", rel(sessFile))
	}
	if git.IsRepo {
		b := git.Branch
		if b == "" {
			b = "(detached)"
		}
		s := git.Sha
		if s == "" {
			s = "(no commit)"
		} else if len(s) >= 12 {
			s = s[:12]
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
	return 0
}

func currentCmd() int {
	id, _ := store.ReadCurrent()
	if id == nil {
		fmt.Println("No current Spacecraft mission. Start one with /sc-start <title>.")
		return 0
	}
	fmt.Println(*id)
	return 0
}

func resolveCmd(args []string) int {
	sel := ""
	if len(args) > 0 && args[0] != "--json" {
		sel = args[0]
	}
	out := r.Resolve(sel)
	jsonOut := false
	for _, a := range args {
		if a == "--json" {
			jsonOut = true
		}
	}
	if jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		enc.Encode(out)
		return 0
	}
	if out.Selected != nil {
		fmt.Printf("Mission: %s\n", out.Selected.ID)
		return 0
	}
	fmt.Println("No mission resolved.")
	return 1
}

func missionsCmd() int {
	records, err := store.List()
	if err != nil {
		return printErr(err)
	}
	ordered := missionDisplayRecords(records)
	if len(ordered) == 0 {
		fmt.Println("No missions.")
		return 0
	}
	res := r.Resolve("")
	sigMap := make(map[string]string)
	for _, s := range res.Signals {
		if s.MissionId != nil {
			sigMap[*s.MissionId] = s.Source
		}
	}
	for i, rec := range ordered {
		num := i + 1
		title := displayTitle(rec.Mission.Title)
		state := rec.Mission.State
		active := ""
		if !rec.Active {
			active = " shipped"
		}
		bh := ""
		if len(rec.Branches) > 0 {
			bh = " [" + strings.Join(rec.Branches, ",") + "]"
		}
		sig := ""
		if s, ok := sigMap[rec.ID]; ok {
			sig = " <" + s + ">"
		}
		fmt.Printf("%d. %s (%s) state:%s%s%s%s\n", num, title, rec.ID, state, active, sig, bh)
	}
	return 0
}

func useCmd(args []string) int {
	if len(args) == 0 {
		return printErr("Missing selector.\nUsage: spacecraft use <number|id|title>")
	}
	sel := strings.TrimSpace(strings.Join(args, " "))
	records, _ := store.List()
	ordered := missionDisplayRecords(records)
	rec := findSelector(records, sel, ordered)
	if rec == nil {
		return printErr(fmt.Sprintf("No mission matches %q.\nRun 'spacecraft missions' to see available missions.", sel))
	}
	if err := store.WriteCurrent(rec.ID); err != nil {
		return printErr("Failed to write .space/current:", err)
	}
	sessFile := writeSessionBinding(rec.ID)
	fmt.Printf("Selected mission %s (%s)\n", rec.ID, rec.Mission.Title)
	if sessFile != "" {
		fmt.Printf("Session: %s\n", rel(sessFile))
	}
	return 0
}

func bindBranchCmd(args []string) int {
	sel := ""
	if len(args) > 0 {
		sel = strings.TrimSpace(strings.Join(args, " "))
	}
	var record *mission.MissionRecord
	if sel != "" {
		records, _ := store.List()
		ordered := missionDisplayRecords(records)
		record = findSelector(records, sel, ordered)
		if record == nil {
			return printErr(fmt.Sprintf("No mission matches %q.", sel))
		}
	} else {
		res := requireResolved("bind-branch")
		if res.Selected == nil {
			return printErr("No resolved mission to bind branch.")
		}
		records, _ := store.List()
		rec := findSelector(records, res.Selected.ID, missionDisplayRecords(records))
		if rec == nil {
			return printErr("Selected mission not found on disk.")
		}
		record = rec
	}
	git := gitutil.GitInfo(gitutil.OSCommandRunner{})
	if !git.IsRepo {
		return printErr("Not a git worktree; cannot bind branch.")
	}
	if git.Branch == "" {
		return printErr("No current branch to bind.")
	}
	m, err := store.Load(record.ID)
	if err != nil {
		return printErr("Failed to read mission.json:", err)
	}
	m.WorkBranch = &git.Branch
	m.Git.WorkBranch = &git.Branch
	m.Branch = &git.Branch
	m.UpdatedAt = isoNow()
	if err := store.Save(m); err != nil {
		return printErr("Failed to write mission.json:", err)
	}
	fmt.Printf("Bound branch %s to mission %s\n", git.Branch, record.ID)
	return 0
}

func statusCmd(args []string) int {
	res := r.Resolve("")
	if res.Safety != "safe" || res.Selected == nil {
		fmt.Println("No selected Spacecraft mission. Start one with /sc-start <title>.")
		for _, c := range res.Candidates {
			fmt.Printf("- %s (%s)\n", c.Title, c.ID)
		}
		return 0
	}
	id := res.Selected.ID
	m, err := store.Load(id)
	if err != nil {
		return printErr("Failed to load mission:", err)
	}
	taskCount := 0
	plan, _ := store.LoadPlan(id)
	if plan != nil && plan.Tasks != nil {
		taskCount = len(plan.Tasks)
	}
	evCount, _ := store.CountEvidence(id)
	src := "unknown"
	if res.Source != nil {
		src = *res.Source
	}
	fmt.Printf("Mission: %s\n", m.ID)
	fmt.Printf("Selected by: %s\n", src)
	if res.Safety != "safe" {
		fmt.Printf("Mission safety: %s\n", res.Safety)
	}
	if res.CurrentMissionId != nil {
		fmt.Printf("Current: %s\n", *res.CurrentMissionId)
	}
	if len(res.Conflicts) > 0 {
		fmt.Println("Conflicts:")
		for _, c := range res.Conflicts {
			fmt.Printf("- %s\n", c.Type)
		}
	}
	fmt.Printf("Title: %s\n", m.Title)
	fmt.Printf("State: %s\n", m.State)
	if m.Clarification.Status != "" {
		fmt.Printf("Clarification: %s\n", m.Clarification.Status)
		fmt.Printf("Blocking questions: %d\n", m.Clarification.BlockingQuestions)
	}
	git := gitutil.GitInfo(gitutil.OSCommandRunner{})
	if git.IsRepo {
		b := git.Branch
		if b == "" {
			b = "(detached)"
		}
		sha := git.Sha
		if sha == "" {
			sha = "(no commit)"
		} else if len(sha) >= 12 {
			sha = sha[:12]
		}
		st := " clean"
		if git.Dirty {
			st = fmt.Sprintf(" dirty:%d", git.DirtyFiles)
		}
		fmt.Printf("Git: %s %s%s\n", b, sha, st)
		if m.BaseSha != nil && git.Sha != "" && *m.BaseSha != git.Sha {
			fmt.Printf("Mission base: %s\n", (*m.BaseSha)[:12])
		}
	} else {
		fmt.Println("Git: not a git worktree")
	}
	fmt.Printf("Tasks: %d\n", taskCount)
	fmt.Printf("Evidence: %d\n", evCount)
	reviewStatus := "not-reviewed"
	if rev, err := store.LoadReview(id); err == nil && rev.Status != nil && *rev.Status != "" {
		reviewStatus = *rev.Status
	}
	fmt.Printf("Review: %s\n", reviewStatus)
	return 0
}

func workflowCmd(args []string) int {
	jsonOut := len(args) > 0 && args[0] == "--json"
	res := r.Resolve("")
	if res.Safety != "safe" || res.Selected == nil {
		return printErr(resolver.FormatResolutionBlock(res, "flow") + "\n")
	}
	snap, err := ws.Build(res, res.Selected.ID)
	if err != nil {
		return printErr(err)
	}
	if jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		enc.Encode(snap)
		return 0
	}
	readyStr := "ready"
	if len(snap.Blockers) > 0 {
		readyStr = "blocked"
	}
	fmt.Printf("Workflow: %s\n", readyStr)
	fmt.Printf("Mission: %s (%s)\n", snap.Title, snap.MissionID)
	fmt.Printf("State: %s\n", snap.State)
	fmt.Printf("Tasks: %d/%d completed\n", snap.Tasks.Completed, snap.Tasks.Total)
	fmt.Printf("Evidence: %d\n", snap.EvidenceCount)
	if snap.NextTask != nil {
		d := "(unnamed task)"
		if snap.NextTask.ID != nil && *snap.NextTask.ID != "" && snap.NextTask.Title != nil {
			d = fmt.Sprintf("%s %s", *snap.NextTask.ID, *snap.NextTask.Title)
		} else if snap.NextTask.ID != nil && *snap.NextTask.ID != "" {
			d = *snap.NextTask.ID
		} else if snap.NextTask.Title != nil {
			d = *snap.NextTask.Title
		}
		fmt.Printf("Next task: %s\n", d)
	}
	fmt.Printf("Next: %s\n", snap.Next)
	fmt.Println("Loop: /sc-work Txx -> /sc-verify Txx -> checkpoint commit -> next task, until a gate blocks.")
	fmt.Printf("Checkpoint: %s\n", snap.CheckpointPolicy)
	if len(snap.Blockers) > 0 {
		fmt.Println("Blockers:")
		for _, b := range snap.Blockers {
			fmt.Printf("- %s\n", b)
		}
	}
	return 0
}

func gitInfoCmd() int {
	git := gitutil.GitInfo(gitutil.OSCommandRunner{})
	if !git.Available {
		fmt.Println("Git: command unavailable")
		return 0
	}
	if !git.IsRepo {
		fmt.Println("Git: not a git worktree")
		return 0
	}
	fmt.Println("Git: worktree detected")
	r := git.Root
	if r == "" {
		r = "(unknown)"
	}
	fmt.Printf("Root: %s\n", r)
	b := git.Branch
	if b == "" {
		b = "(detached)"
	}
	fmt.Printf("Branch: %s\n", b)
	s := git.Sha
	if s == "" {
		s = "(no commit)"
	}
	fmt.Printf("HEAD: %s\n", s)
	st := "clean"
	if git.Dirty {
		st = fmt.Sprintf("dirty (%d files)", git.DirtyFiles)
	}
	fmt.Printf("Status: %s\n", st)
	return 0
}

func gitSuggestCmd(args []string) int {
	res := requireResolved("git-suggest")
	id := res.Selected.ID
	m, _ := store.Load(id)

	branchTypes := map[string]bool{
		"feat": true, "fix": true, "docs": true, "refactor": true,
		"test": true, "build": true, "ci": true, "chore": true,
		"perf": true, "style": true, "issue": true, "release": true,
	}
	reqType := ""
	if len(args) > 0 {
		reqType = strings.ToLower(args[0])
	}
	typ := "feat"
	var slugParts []string
	if branchTypes[reqType] {
		typ = reqType
		slugParts = args[1:]
	} else {
		slugParts = args
	}
	slugBase := strings.Join(slugParts, " ")
	if slugBase == "" && m != nil && m.Title != "" {
		slugBase = m.Title
	} else if slugBase == "" {
		slugBase = "mission"
	}
	slug := slugify(slugBase)
	mp := strings.ToLower(id)
	if strings.HasPrefix(slug, mp+"-") {
		slug = slug[len(mp)+1:]
	} else if strings.HasPrefix(slug, mp) && slug != mp {
		slug = slug[len(mp):]
	}
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "mission"
	}
	branch := fmt.Sprintf("%s/%s/%s", typ, mp, slug)
	ct := typ
	if typ == "issue" || typ == "release" {
		ct = "chore"
	}
	fmt.Println("Spacecraft git strategy: release branching")
	fmt.Printf("Mission: %s (%s)\n", res.Selected.Title, res.Selected.ID)
	src := "unknown"
	if res.Source != nil {
		src = *res.Source
	}
	fmt.Printf("Selected by: %s\n", src)
	fmt.Println("Base: latest main")
	fmt.Println("Main: protected; no direct writes")
	fmt.Printf("Branch: %s\n", branch)
	if typ == "release" {
		fmt.Println("Rule: release branches are for version/changelog/spec preparation only.")
	} else {
		fmt.Println("Rule: one branch per feature/issue/tightly scoped change.")
	}
	fmt.Println("Final branch history: target 1-3 commits; maximum 5 unless justified.")
	fmt.Println("Before merge: rebase on latest main, verify, bump version, update changelog/spec.")
	fmt.Println("Merge: git merge --no-ff <branch>")
	fmt.Println("After merge: create annotated version tag.")
	fmt.Println("Use a separate worktree for large, risky, or multi-session branches.")
	fmt.Println("")
	fmt.Println("Conventional Commits format:")
	fmt.Println("<type>[optional scope]: <description>")
	fmt.Println("")
	fmt.Println("Common types: feat, fix, docs, refactor, test, build, ci, chore, perf, style, revert")
	fmt.Println("Examples:")
	if typ == "release" {
		fmt.Println("chore: prepare release v0.1.0")
	} else {
		fmt.Printf("%s: add focused mission change\n", ct)
	}
	fmt.Println("docs: update changelog for v0.2.0")
	fmt.Println("chore: bump version to v0.2.0")
	return 0
}

func setStateCmd(args []string) int {
	if len(args) == 0 {
		return printErr("Missing state.\nAllowed states: draft, planned, verifying, reviewing, ready, shipped, blocked")
	}
	res := requireResolved("set-state")
	state := args[0]
	if err := ss.SetState(res.Selected.ID, state); err != nil {
		return printErr(err)
	}
	fmt.Printf("Spacecraft mission %s state: %s\n", res.Selected.ID, stateDisplay(state))
	return 0
}

func stateDisplay(s string) string {
	switch s {
	case "specified":
		return "draft"
	case "implementing":
		return "planned"
	}
	return s
}

func clarifyStatusCmd(args []string) int {
	if len(args) == 0 {
		return printErr("Missing clarification status.\nAllowed statuses: open, clear, deferred")
	}
	res := requireResolved("clarify-status")
	if err := ss.SetClarificationStatus(res.Selected.ID, args[0]); err != nil {
		return printErr(err)
	}
	fmt.Printf("Spacecraft mission %s clarification: %s\n", res.Selected.ID, args[0])
	return 0
}

func evidenceCmd(args []string) int {
	sep := -1
	for i, a := range args {
		if a == "--" {
			sep = i
			break
		}
	}
	if sep == -1 {
		return printErr("Missing -- before evidence command.\n\n" + usage())
	}
	label := strings.TrimSpace(strings.Join(args[:sep], " "))
	cmdParts := args[sep+1:]
	if label == "" {
		return printErr("Missing evidence label.\n\n" + usage())
	}
	if len(cmdParts) == 0 {
		return printErr("Missing command after --.\n\n" + usage())
	}
	res := r.Resolve("")
	if res.Safety != "safe" || res.Selected == nil {
		return printErr(resolver.FormatResolutionBlock(res, "evidence"))
	}
	id := res.Selected.ID
	os.MkdirAll(filepath.Join(store.MissionDir(id), "outputs"), 0755)
	eid, stdoutP, stderrP, err := store.ReserveEvidencePath(id)
	if err != nil {
		return printErr("Failed to reserve evidence path:", err)
	}
	result := execCmd(cmdParts)
	os.WriteFile(stdoutP, []byte(result.stdout), 0644)
	os.WriteFile(stderrP, []byte(result.stderr), 0644)
	entry := mission.EvidenceEntry{
		ID:        eid,
		Label:     label,
		Command:   cmdToStr(cmdParts),
		ExitCode:  result.code,
		Stdout:    rel(stdoutP),
		Stderr:    rel(stderrP),
		CreatedAt: isoNow(),
	}
	if err := store.AppendEvidence(id, &entry); err != nil {
		return printErr("Failed to append evidence:", err)
	}
	fmt.Printf("Evidence: %s\n", eid)
	fmt.Printf("Exit code: %d\n", result.code)
	return result.code
}

type execResult struct {
	code   int
	stdout string
	stderr string
}

func execCmd(parts []string) execResult {
	if len(parts) == 0 {
		return execResult{code: 1}
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	var outB, errB bytes.Buffer
	cmd.Stdout = &outB
	cmd.Stderr = &errB
	cmd.Dir = cfg.Root()
	cmd.Env = os.Environ()

	err := cmd.Run()
	code := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else {
			code = 127
			errB.WriteString(err.Error() + "\n")
		}
	}
	return execResult{code: code, stdout: outB.String(), stderr: errB.String()}
}

func validateCmd() int {
	res := requireResolved("validate")
	errs := ss.ValidateMission(res.Selected.ID)
	if errs != nil {
		fmt.Fprintf(os.Stderr, "Spacecraft mission %s is invalid:\n", res.Selected.ID)
		for _, e := range errs.Errors {
			fmt.Fprintf(os.Stderr, "- %s\n", e)
		}
		return 1
	}
	fmt.Printf("Spacecraft mission %s is valid.\n", res.Selected.ID)
	return 0
}

func closeoutCmd() int {
	res := requireResolved("closeout-check")
	id := res.Selected.ID
	m, errM := store.Load(id)
	if errM != nil {
		return printErr("Missing mission.json:", errM)
	}
	var plan *mission.Plan
	if p, err := store.LoadPlan(id); err == nil {
		plan = p
	}
	var review *mission.Review
	if r, err := store.LoadReview(id); err == nil {
		review = r
	}
	evCount, _ := store.CountEvidence(id)
	result := cc.Check(id, m, plan, review, evCount)
	if len(result.Errors) > 0 {
		fmt.Printf("Spacecraft closeout blocked for %s:\n", id)
		for _, e := range result.Errors {
			fmt.Printf("- %s\n", e)
		}
		if len(result.Warnings) > 0 {
			fmt.Println("Warnings:")
			for _, w := range result.Warnings {
				fmt.Printf("- %s\n", w)
			}
		}
		return 1
	}
	fmt.Printf("Spacecraft closeout ready for %s.\n", id)
	fmt.Println("Next: rebase already satisfied, run final verification if needed, merge with git merge --no-ff <branch>, tag release, then delete merged branch unless kept.")
	if len(result.Warnings) > 0 {
		fmt.Println("Warnings:")
		for _, w := range result.Warnings {
			fmt.Printf("- %s\n", w)
		}
	}
	return 0
}

func archiveCmd(args []string) int {
	var res mission.ResolveOutput
	if len(args) > 0 {
		res = r.Resolve(strings.Join(args, " "))
	} else {
		var err error
		res, err = r.RequireResolved("archive")
		if err != nil {
			return printErr(resolver.FormatResolutionBlock(res, "archive"))
		}
	}
	if res.Safety != "safe" || res.Selected == nil {
		return printErr(resolver.FormatResolutionBlock(res, "archive"))
	}
	id := res.Selected.ID
	m, err := store.Load(id)
	if err != nil {
		return printErr(err)
	}
	if m.State != "shipped" {
		return printErr(fmt.Sprintf("Archive blocked: mission %s state is %s. Archive only shipped missions.", id, m.State))
	}
	var plan *mission.Plan
	if p, err := store.LoadPlan(id); err == nil {
		plan = p
	}
	var review *mission.Review
	if r, err := store.LoadReview(id); err == nil {
		review = r
	}
	entries, _ := store.ReadEvidenceEntries(id)
	errs := arc.CheckReadiness(id, plan, review, entries)
	if errs != nil {
		return printErr(fmt.Sprintf("Archive blocked for %s:\n- %s", id, strings.Join(errs.Errors, "\n- ")))
	}
	result, err := ar.Archive(archive.ArchiveParams{
		ID: id, Mission: m, Plan: plan, Review: review, EvidenceEntries: entries,
	})
	if err != nil {
		return printErr(err)
	}
	fmt.Printf("Archived mission %s\n", id)
	fmt.Printf("Archive: %s\n", rel(result.ArchiveDir))
	return 0
}

// ---------- helpers ----------

func requireResolved(cmd string) mission.ResolveOutput {
	r2, err := r.RequireResolved(cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Resolution failed or blocked.")
		os.Exit(1)
	}
	return r2
}

func missionDisplayRecords(records []mission.MissionRecord) []mission.MissionRecord {
	res := make([]mission.MissionRecord, len(records))
	copy(res, records)
	sort.Slice(res, func(i, j int) bool {
		if res[i].Active != res[j].Active {
			return res[i].Active
		}
		return res[i].ID > res[j].ID
	})
	return res
}

func displayTitle(title string) string {
	t := strings.TrimSpace(title)
	if t == "" {
		return "(untitled)"
	}
	if len(t) <= 88 {
		return t
	}
	return t[:85] + "..."
}

func findSelector(records []mission.MissionRecord, sel string, ordered []mission.MissionRecord) *mission.MissionRecord {
	text := strings.TrimSpace(sel)
	if text == "" {
		return nil
	}
	if num, err := strconv.Atoi(text); err == nil && num > 0 && num <= len(ordered) {
		return &ordered[num-1]
	}
	if id := normMissionID(text); id != nil {
		for i := range records {
			if records[i].ID == *id {
				return &records[i]
			}
		}
	}
	var exact []mission.MissionRecord
	for _, r := range records {
		if r.Mission != nil && r.Mission.Title == text {
			exact = append(exact, r)
		}
	}
	if len(exact) == 1 {
		return &exact[0]
	}
	var matches []mission.MissionRecord
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

func writeSessionBinding(id string) string {
	for _, env := range []string{"SPACECRAFT_SESSION", "OPENCODE_SESSION_ID", "CODEX_SESSION_ID"} {
		if key := os.Getenv(env); key != "" {
			if err := store.WriteSession(key, id); err == nil {
				return store.SessionFilePath(key)
			}
		}
	}
	return ""
}

func rel(p string) string {
	r, err := filepath.Rel(cfg.Root(), p)
	if err != nil || r == "" {
		return p
	}
	return r
}

func isoNow() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ifStr(cond bool, t, f string) string {
	if cond {
		return t
	}
	return f
}

func printErr(v ...interface{}) int {
	fmt.Fprintln(os.Stderr, v...)
	return 1
}

// ID normalization
var legacyRe = regexp.MustCompile(`\b[Mm]-(\d{8}-\d{6})\b`)
var compactRe = regexp.MustCompile(`(?:^|[^A-Za-z0-9])([Mm][0-9A-Za-z]{8})(?:$|[^A-Za-z0-9])`)

func normMissionID(value string) *string {
	text := strings.TrimSpace(value)
	if m := legacyRe.FindStringSubmatch(text); m != nil {
		s := "M-" + m[1]
		return &s
	}
	if m := compactRe.FindStringSubmatch(text); m != nil {
		s := strings.ToUpper(m[1])
		return &s
	}
	return nil
}

func slugify(value string) string {
	text := strings.ToLower(strings.TrimSpace(value))
	var sb strings.Builder
	for _, r := range text {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
		} else {
			sb.WriteRune('-')
		}
	}
	slug := sb.String()
	parts := strings.FieldsFunc(slug, func(r rune) bool { return r == '-' })
	if len(parts) == 0 {
		return "mission"
	}
	slug = strings.Join(parts, "-")
	slug = strings.Trim(slug, "-")
	if len(slug) > 60 {
		slug = slug[:60]
		slug = strings.TrimRight(slug, "-")
	}
	if slug == "" {
		return "mission"
	}
	return slug
}

// safeCmdQuote quotes a command part for display
var safeRe = regexp.MustCompile(`^[A-Za-z0-9_./:=@%+-]+$`)

func cmdToStr(parts []string) string {
	var res []string
	for _, p := range parts {
		if safeRe.MatchString(p) {
			res = append(res, p)
		} else {
			esc := strings.ReplaceAll(p, "'", "'\\''")
			res = append(res, "'"+esc+"'")
		}
	}
	return strings.Join(res, " ")
}
