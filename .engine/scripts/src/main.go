package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
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
	"spacecraft/internal/eval"
	"spacecraft/internal/gitutil"
	"spacecraft/internal/hooks"
	"spacecraft/internal/id"
	"spacecraft/internal/mission"
	"spacecraft/internal/research"
	"spacecraft/internal/roadmap"
	"spacecraft/internal/resolver"
	"spacecraft/internal/state"
	"spacecraft/internal/trace"
	"spacecraft/internal/workflow"
)

// CLI dependencies set during init.
var (
	cfg          *config.Config
	store        *mission.FSStore
	traceStore   *trace.FSTraceStore
	roadmapStore *roadmap.FSStore
	r            *resolver.Resolver
	ss          *state.StateSetter
	ws          *workflow.Snapshot
	cc          *closeout.Checker
	arc         *archive.ReadinessChecker
	ar          *archive.MissionArchiver
	hooksCfg    *hooks.Config

	// Overridable endpoints for testing.
	braveBaseURL   = "https://api.search.brave.com"
	npmRegistryURL = "https://registry.npmjs.org"
	goProxyURL     = "https://proxy.golang.org"
	pypiURL        = "https://pypi.org"
	cratesURL      = "https://crates.io"
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
	traceStore = trace.NewFSTraceStore(cfg)
	roadmapStore = roadmap.NewFSStore(cfg)
	r = resolver.New(store, gitutil.OSCommandRunner{}, nil)
	ss = state.NewSetter(store)
	ws = workflow.NewSnapshot(store)
	ws.SetCommandsDir(filepath.Join(cfg.Root(), ".opencode", "commands"))
	cc = closeout.NewChecker(store, gitutil.OSCommandRunner{})
	arc = archive.NewReadinessChecker(store)
	ar = archive.NewArchiverWithRoadmap(store, roadmapStore)
	hooksCfg, _ = hooks.LoadConfig(filepath.Join(cfg.SpaceDir(), "hooks.json"))

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
	case "research":
		exitCode = researchCmd(args)
	case "check-deps":
		exitCode = checkDepsCmd(args)
	case "eval":
		exitCode = evalCmd(args)
	case "traces":
		exitCode = tracesCmd(args)
	case "cost":
		exitCode = costCmd(args)
	case "roadmap":
		exitCode = roadmapCmd(args)
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
  spacecraft research <query> [flags]
  spacecraft check-deps [flags]
  spacecraft eval <mission-id>
  spacecraft eval init <mission-id>
  spacecraft traces <mission-id> [--json] [--verbose] [--flat]
  spacecraft cost [--mission <id>] [--all]
  spacecraft roadmap new <title> [--desc <text>]
  spacecraft roadmap add <roadmap-id> <mission-id> [--after <target-mission-id>]
  spacecraft roadmap remove <roadmap-id> <mission-id>
  spacecraft roadmap show <roadmap-id>
  spacecraft roadmap list
  spacecraft roadmap continue <roadmap-id>
  spacecraft roadmap archive <roadmap-id>
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

	ctx := context.Background()
	if hooks.Fire(ctx, hooksCfg, "mission.created") != nil {
		return 1
	}
	if hooks.Fire(ctx, hooksCfg, "mission.state.changed") != nil {
		return 1
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
		if out.Selected != nil {
			return 0
		}
		return 1
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
	fmt.Println("Loop: /sc-build Txx -> checkpoint commit -> next task, until a gate blocks.")
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
		return printErr("Missing state.\nAllowed states: draft, planned, built, ready, shipped, blocked")
	}
	res := requireResolved("set-state")
	state := args[0]
	if err := ss.SetState(res.Selected.ID, state); err != nil {
		return printErr(err)
	}
	if hooks.Fire(context.Background(), hooksCfg, "mission.state.changed") != nil {
		return 1
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

	// I13: reject placeholder evidence — exitCode 0 with no real output
	if result.code == 0 && len(strings.TrimSpace(result.stdout)) == 0 {
		return printErr("Evidence command produced no output (exit 0, empty stdout).\n" +
			"Use a command that demonstrates actual behavior, not echo/placeholder.\n\n" + usage())
	}

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

	fmt.Printf("Evidence: %s\n", eid)
	fmt.Printf("Exit code: %d\n", result.code)

	if err := store.AppendEvidence(id, &entry); err != nil {
		return printErr("Failed to append evidence:", err)
	}
	_ = hooks.Fire(context.Background(), hooksCfg, "mission.evidence.appended")
	return 0
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
	_ = hooks.Fire(context.Background(), hooksCfg, "mission.validated")
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
	_ = hooks.Fire(context.Background(), hooksCfg, "mission.archived")
	fmt.Printf("Archived mission %s\n", id)
	fmt.Printf("Archive: %s\n", rel(result.ArchiveDir))
	return 0
}

// ---------- research command ----------

// normalizeDeepArgs converts bare "--deep" occurrences to "--deep=true" so
// that flag.FlagSet can treat --deep as an optional string flag.
func normalizeDeepArgs(args []string) []string {
	out := make([]string, 0, len(args))
	for i, a := range args {
		if a == "--deep" && (i == len(args)-1 || strings.HasPrefix(args[i+1], "-")) {
			out = append(out, "--deep=true")
			continue
		}
		out = append(out, a)
	}
	return out
}

func researchCmd(args []string) int {
	args = normalizeDeepArgs(args)

	fs := flag.NewFlagSet("research", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() { fmt.Print(researchUsage()) }

	var (
		scope    string
		deepMode string
		results  int
		timeout  time.Duration
		jsonOut  bool
		noSave   bool
	)
	fs.StringVar(&scope, "scope", "", "Force a scope")
	fs.StringVar(&deepMode, "deep", "", "Deep research (true=browser-use, nlm=notebooklm)")
	fs.IntVar(&results, "results", 5, "Number of results")
	fs.DurationVar(&timeout, "timeout", 10*time.Second, "Request timeout")
	fs.BoolVar(&jsonOut, "json", false, "JSON output")
	fs.BoolVar(&noSave, "no-save", false, "Skip persistence")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}
	args = fs.Args()

	if len(args) == 0 {
		return printErr("Missing query.\n\n" + researchUsage())
	}
	query := strings.TrimSpace(strings.Join(args, " "))
	if query == "" {
		return printErr("Missing query.\n\n" + researchUsage())
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if deepMode != "" && deepMode != "true" && deepMode != "nlm" {
		return fatalErr("Invalid --deep value:", deepMode, "(expected 'true' or 'nlm')")
	}

	var (
		outputResults interface{}
		outputSource  string
		outputScope   string
		outputMethod  string
		outputDeep    *research.DeepResult
	)

	// Determine if the query looks like a package identifier.
	// Simple heuristic: single token, no spaces, common package patterns.
	pkgFound := false
	if isPackageQuery(query) {
		// Try registry lookup across all registries.
		pkg, src := lookupPackage(ctx, query, timeout)
		if pkg != nil {
			outputResults = pkg
			outputSource = src
			pkgFound = true
		}
		// Registry lookup failed — fall through to Brave Search.
	}

	if !pkgFound {
		// Scope detection.
		{
			detectedScope := scope
			if detectedScope == "" {
				// Auto-detect from query + project manifests.
				manifests := detectManifests()
				detectedScope = research.SmartScope(query, manifests)
			}
			outputScope = detectedScope

			// Load scope config once (built-in defaults + .space/scopes.json override).
			scopeCfg, err := research.LoadScopes(filepath.Join(cfg.Root(), ".space", "scopes.json"))
			if err != nil {
				return fatalErr("Failed to load search scopes:", err)
			}

			// Validate custom scope against known scopes.
			if scope != "" && detectedScope != "" {
				if _, ok := scopeCfg.Scopes[detectedScope]; !ok {
					return fatalErr(fmt.Sprintf("Unknown scope: %q", scope))
				}
			}

			// Get Brave Search API key.
			apiKey := os.Getenv("SPACECRAFT_BRAVE_API_KEY")
			if apiKey == "" {
				return fatalErr("SPACECRAFT_BRAVE_API_KEY environment variable is not set. Get a key at https://brave.com/search/api/")
			}

			// Build domain list from scope config.
			var domains []string
			if detectedScope != "" {
				if s, ok := scopeCfg.Scopes[detectedScope]; ok {
					domains = s.Domains
				}
			}

			brave := research.NewBraveClient(apiKey, braveBaseURL)
			searchResults, err := brave.Search(ctx, query, domains, results)
			if err != nil {
				return fatalErr("Brave Search error:", err)
			}

			// Truncate to --results limit.
			if len(searchResults) > results {
				searchResults = searchResults[:results]
			}

			outputResults = searchResults
			outputSource = "brave-search"

			// Deep analysis: if --deep is set, analyze the top result URL.
			if deepMode != "" && len(searchResults) > 0 {
				topURL := searchResults[0].URL
				deepResult, deepErr := runDeepAnalysis(ctx, topURL, query, deepMode, timeout)
				if deepErr != nil {
					fmt.Fprintf(os.Stderr, "Deep analysis warning: %v\n", deepErr)
				} else {
					outputDeep = deepResult
				}
				switch deepMode {
				case "true":
					outputMethod = "browser-use"
				case "nlm":
					outputMethod = "nlm"
				}
			}

			if len(searchResults) == 0 {
				fmt.Println("No results found.")
				return 1
			}
		}
	}

	// Determine persistence directory.
	persistDir := ""
	if !noSave {
		res := r.Resolve("")
		if res.Selected != nil {
			persistDir = filepath.Join(store.MissionDir(res.Selected.ID), "research")
		} else {
			persistDir = filepath.Join(cfg.Root(), ".space", "research")
		}
	}

	opts := research.FormatOptions{
		JSON:         jsonOut,
		Query:        query,
		Scope:        outputScope,
		Source:       outputSource,
		Method:       outputMethod,
		DeepAnalysis: outputDeep,
		Timestamp:    time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		PersistDir:   persistDir,
	}

	if err := research.FormatResults(os.Stdout, outputResults, opts); err != nil {
		return printErr("Formatting error:", err)
	}

	// Persist to disk.
	if persistDir != "" {
		prefix := slugify(query)
		if len(prefix) > 40 {
			prefix = prefix[:40]
		}
		path, err := research.SaveResults(persistDir, prefix, outputResults, opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to save results: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "\nSaved: %s\n", rel(path))
		}
	}

	return 0
}

// runDeepAnalysis invokes the deep-research backend for a single URL.
// deep: "true" = browser-use, "nlm" = notebooklm-mcp-cli.
func runDeepAnalysis(ctx context.Context, url, query, deepMode string, timeout time.Duration) (*research.DeepResult, error) {
	switch deepMode {
	case "true":
		runner := research.NewBrowserUseRunner(&research.OSExecutor{})
		if !runner.IsAvailable() {
			return nil, fmt.Errorf("browser-use not available. Install with: %s", runner.InstallInstructions())
		}
		// Extend timeout for deep analysis (use 2x the search timeout or default 30s).
		deepCtx, cancel := context.WithTimeout(ctx, timeout*2)
		defer cancel()
		return runner.Analyze(deepCtx, url)
	case "nlm":
		runner := research.NewNotebookLMRunner(&research.OSExecutor{})
		if !runner.IsAvailable() {
			return nil, fmt.Errorf("notebooklm-mcp-cli not available. Install with: %s", runner.InstallInstructions())
		}
		deepCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
		return runner.Analyze(deepCtx, query)
	default:
		return nil, fmt.Errorf("unknown deep mode: %s", deepMode)
	}
}

func researchUsage() string {
	return `spacecraft research <query> [flags]

Flags:
  --scope <name>     Force a scope (react, tailwindcss, nextjs, storybook, postgresql, go, rust)
  --deep [true|nlm]  Deep research (true=browser-use, nlm=notebooklm)
  --results <n>      Number of results (default 5)
  --timeout <d>      Request timeout (default 10s)
  --json             JSON output
  --no-save          Skip persistence
`
}

// isPackageQuery returns true if the query looks like a package identifier.
// Requires at least one of: dot, slash, or @-scope prefix to avoid triggering
// 4 registry lookups for every single-token search query.
func isPackageQuery(query string) bool {
	if strings.Contains(query, " ") {
		return false
	}
	// Must contain a dot, slash, or @ to look like a package identifier.
	if !strings.ContainsAny(query, "./@") {
		return false
	}
	for _, r := range query {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '.' || r == '-' ||
			r == '/' || r == '_' || r == '@' {
			continue
		}
		return false
	}
	return len(query) > 0
}

// lookupPackage tries each registry in order and returns the first match.
func lookupPackage(ctx context.Context, query string, timeout time.Duration) (*research.PackageInfo, string) {
	clients := []struct {
		name   string
		client interface{ Lookup(context.Context, string) (*research.PackageInfo, error) }
	}{
		{"npm", research.NewNpmClientWithTimeout(npmRegistryURL, timeout)},
		{"go", research.NewGoProxyClientWithTimeout(goProxyURL, timeout)},
		{"pypi", research.NewPypiClientWithTimeout(pypiURL, timeout)},
		{"crates", research.NewCargoClientWithTimeout(cratesURL, timeout)},
	}
	for _, c := range clients {
		pkg, err := c.client.Lookup(ctx, query)
		if err == nil && pkg != nil && pkg.Name != "" {
			return pkg, c.name
		}
	}
	fmt.Fprintf(os.Stderr, "No registry match for package query %q (tried: npm, go, pypi, crates)\n", query)
	return nil, ""
}

// detectManifests returns a list of manifest filenames found in the project root.
func detectManifests() []string {
	root := cfg.Root()
	var manifests []string
	for _, name := range []string{"go.mod", "package.json", "requirements.txt", "pyproject.toml", "Cargo.toml"} {
		if _, err := os.Stat(filepath.Join(root, name)); err == nil {
			manifests = append(manifests, name)
		}
	}
	return manifests
}

// ---------- eval command ----------

func evalCmd(args []string) int {
	if len(args) == 0 {
		fmt.Print(evalUsage())
		return 1
	}

	sub := args[0]
	switch sub {
	case "init":
		return evalInitCmd(args[1:])
	case "-h", "--help", "help":
		fmt.Print(evalUsage())
		return 0
	default:
		return evalRunCmd(args)
	}
}

func evalRunCmd(args []string) int {
	res := resolveOrUseCurrent(args)
	if res.Safety != "safe" || res.Selected == nil {
		return printErr(resolver.FormatResolutionBlock(res, "eval"))
	}
	id := res.Selected.ID

	runner := eval.NewRunner(store, cfg.EvalsDir())
	runResult, err := runner.Run(id)
	if err != nil {
		return printErr("eval:", err)
	}

	fmt.Printf("Eval complete for %s\n", id)
	fmt.Printf("Deterministic: allPassed=%v\n", runResult.EvalResult.Deterministic.AllPassed)
	fmt.Printf("Rubric average: %.1f/4.0\n", runResult.EvalResult.Scorecard.Average)
	if runResult.EvalResult.LMJudge.Fallback {
		fmt.Printf("LM Judge: fallback (%s)\n", runResult.EvalResult.LMJudge.FallbackReason)
	} else {
		fmt.Printf("LM Judge: %d/4 (model: %s)\n", runResult.EvalResult.LMJudge.Score, runResult.EvalResult.LMJudge.Model)
	}
	fmt.Printf("Coverage: %.2f (%d/%d checks)\n", runResult.EvalResult.Coverage, runResult.EvalResult.CoveredChecks, runResult.EvalResult.TotalChecks)
	if runResult.EvalResult.CoverageSatisfied {
		fmt.Println("Coverage threshold: satisfied")
	} else {
		fmt.Println("Coverage threshold: NOT satisfied")
	}
	fmt.Printf("Evidence: %s\n", runResult.Entry.ID)

	return 0
}

func evalInitCmd(args []string) int {
	res := resolveOrUseCurrent(args)
	if res.Safety != "safe" || res.Selected == nil {
		return printErr(resolver.FormatResolutionBlock(res, "eval init"))
	}
	id := res.Selected.ID
	if err := eval.Init(cfg.EvalsDir(), id); err != nil {
		return printErr(err)
	}
	return 0
}

func resolveOrUseCurrent(args []string) mission.ResolveOutput {
	if len(args) > 0 {
		return r.Resolve(strings.Join(args, " "))
	}
	result, err := r.RequireResolved("eval")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Resolution failed or blocked. Provide a mission id or selector.")
		return r.Resolve("")
	}
	return result
}

func evalUsage() string {
	return `spacecraft eval <mission-id>         Run eval suite against mission evidence
spacecraft eval init <mission-id>   Scaffold eval directory
`
}

// ---------- check-deps command ----------

func checkDepsCmd(args []string) int {
	fs := flag.NewFlagSet("check-deps", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() { fmt.Print(checkDepsUsage()) }

	var (
		timeout     time.Duration
		concurrency int
		jsonOut     bool
		registry    string
	)
	fs.DurationVar(&timeout, "timeout", 10*time.Second, "Request timeout")
	fs.IntVar(&concurrency, "concurrency", 4, "Concurrent registry lookups")
	fs.StringVar(&registry, "registry", "", "Limit to one ecosystem")
	fs.BoolVar(&jsonOut, "json", false, "JSON output")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}

	root := cfg.Root()
	var manifestFiles []string

	if registry != "" {
		// Limit to specific ecosystem manifest.
		switch registry {
		case "go":
			manifestFiles = []string{filepath.Join(root, "go.mod")}
		case "npm":
			manifestFiles = []string{filepath.Join(root, "package.json")}
		case "pypi":
			for _, n := range []string{"requirements.txt", "pyproject.toml"} {
				p := filepath.Join(root, n)
				if _, err := os.Stat(p); err == nil {
					manifestFiles = append(manifestFiles, p)
				}
			}
		case "crates":
			manifestFiles = []string{filepath.Join(root, "Cargo.toml")}
		default:
			return fatalErr(fmt.Sprintf("Unknown registry: %q. Supported: go, npm, pypi, crates", registry))
		}
	} else {
		for _, name := range []string{"go.mod", "package.json", "requirements.txt", "pyproject.toml", "Cargo.toml"} {
			p := filepath.Join(root, name)
			if _, err := os.Stat(p); err == nil {
				manifestFiles = append(manifestFiles, p)
			}
		}
	}

	if len(manifestFiles) == 0 {
		fmt.Println("No supported dependency manifest files found in", root)
		fmt.Println("Supported: go.mod, package.json, requirements.txt, pyproject.toml, Cargo.toml")
		return 0
	}

	// Parse all manifests into a flat dependency list.
	type depWithSource struct {
		research.Dependency
		LatestVersion string
		UpgradeType   string
	}
	var flatDeps []research.Dependency
	for _, mf := range manifestFiles {
		deps, err := research.ParseManifest(mf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %s: %v\n", filepath.Base(mf), err)
			continue
		}
		flatDeps = append(flatDeps, deps...)
	}

	// Concurrent registry lookups bounded by --concurrency.
	type depResult struct {
		index        int
		depWithSource depWithSource
		err          bool
	}
	sem := make(chan struct{}, concurrency)
	resultCh := make(chan depResult, len(flatDeps))

	for i, d := range flatDeps {
		go func(idx int, dep research.Dependency) {
			sem <- struct{}{}
			defer func() { <-sem }()

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			latest, err := lookupLatestVersion(ctx, dep, timeout)

			dws := depWithSource{Dependency: dep}
			if err != nil {
				dws.LatestVersion = "?"
				dws.UpgradeType = "ERROR"
				resultCh <- depResult{index: idx, depWithSource: dws, err: true}
				return
			}
			dws.LatestVersion = latest
			dws.UpgradeType = compareNumericVersions(dep.CurrentVersion, latest)
			resultCh <- depResult{index: idx, depWithSource: dws}
		}(i, d)
	}

	// Collect results in original order.
	ordered := make([]depWithSource, len(flatDeps))
	hasErrors := false
	for range flatDeps {
		r := <-resultCh
		ordered[r.index] = r.depWithSource
		if r.err {
			hasErrors = true
		}
	}
	allDeps := ordered

	if jsonOut {
		type depJSON struct {
			Name           string `json:"name"`
			Ecosystem      string `json:"ecosystem"`
			CurrentVersion string `json:"current_version"`
			LatestVersion  string `json:"latest_version"`
			UpgradeType    string `json:"upgrade_type"`
		}
		var out []depJSON
		for _, d := range allDeps {
			out = append(out, depJSON{
				Name:           d.Name,
				Ecosystem:      d.Ecosystem,
				CurrentVersion: d.CurrentVersion,
				LatestVersion:  d.LatestVersion,
				UpgradeType:    d.UpgradeType,
			})
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(out)
	} else {
		for _, d := range allDeps {
			flag := ""
			switch d.UpgradeType {
			case "MAJOR UPGRADE":
				flag = " [MAJOR UPGRADE]"
			case "MINOR UPGRADE":
				flag = " [MINOR UPGRADE]"
			case "PATCH":
				flag = " [PATCH]"
			case "current":
				flag = " [current]"
			case "ERROR":
				flag = " [ERROR]"
			}
			fmt.Printf("%-50s %-10s → %-10s%s\n", d.Name, d.CurrentVersion, d.LatestVersion, flag)
		}
	}

	if hasErrors {
		return 1
	}
	return 0
}

func checkDepsUsage() string {
	return `spacecraft check-deps [flags]

Flags:
  --registry <ecosystem>  Limit to one ecosystem (go, npm, pypi, crates)
  --timeout <d>           Request timeout (default 10s)
  --concurrency <n>       Concurrent registry lookups (default 4)
  --json                  JSON output
`
}

// lookupLatestVersion looks up the latest version of a dependency using the matching registry.
// timeout is used for the HTTP client timeout.
func lookupLatestVersion(ctx context.Context, dep research.Dependency, timeout time.Duration) (string, error) {
	switch dep.Ecosystem {
	case "go":
		c := research.NewGoProxyClientWithTimeout(goProxyURL, timeout)
		pkg, err := c.Lookup(ctx, dep.Name)
		if err != nil {
			return "", err
		}
		return pkg.LatestVersion, nil
	case "npm":
		c := research.NewNpmClientWithTimeout(npmRegistryURL, timeout)
		pkg, err := c.Lookup(ctx, dep.Name)
		if err != nil {
			return "", err
		}
		return pkg.LatestVersion, nil
	case "pypi":
		c := research.NewPypiClientWithTimeout(pypiURL, timeout)
		pkg, err := c.Lookup(ctx, dep.Name)
		if err != nil {
			return "", err
		}
		return pkg.LatestVersion, nil
	case "crates":
		c := research.NewCargoClientWithTimeout(cratesURL, timeout)
		pkg, err := c.Lookup(ctx, dep.Name)
		if err != nil {
			return "", err
		}
		return pkg.LatestVersion, nil
	default:
		return "", fmt.Errorf("unknown ecosystem: %s", dep.Ecosystem)
	}
}

// compareNumericVersions compares two version strings and returns the upgrade type.
// Strips leading non-numeric characters (v, ^, ~, =, >, <) and pre-release suffixes,
// then compares major.minor.patch components.
func compareNumericVersions(current, latest string) string {
	cur := extractNumericVersion(current)
	lat := extractNumericVersion(latest)

	curParts := splitVersion(cur)
	latParts := splitVersion(lat)

	if len(curParts) == 0 || len(latParts) == 0 {
		return "current"
	}

	// Normalize to at least 3 parts. Extra parts beyond 3 are used
	// in tiebreakers (e.g., 1.2.3.4 vs 1.2.3.5 → PATCH upgrade).
	for len(curParts) < len(latParts) {
		curParts = append(curParts, 0)
	}
	for len(latParts) < len(curParts) {
		latParts = append(latParts, 0)
	}

	// Check pre-release: if the base versions match but one side has a
	// pre-release suffix, resolve accordingly.
	curBase := extractNumericVersion(current)
	latBase := extractNumericVersion(latest)
	curHasPre := strings.Contains(current, "-") && !strings.Contains(latest, "-")
	latHasPre := strings.Contains(latest, "-") && !strings.Contains(current, "-")

	if curBase == latBase {
		if latHasPre {
			// Latest is a pre-release while current is a full release → no upgrade.
			return "current"
		}
		if curHasPre {
			// Current is a pre-release and latest is a full release → upgrade.
			return "PATCH"
		}
	}

	for i := 0; i < len(latParts); i++ {
		if latParts[i] > curParts[i] {
			switch {
			case i == 0:
				return "MAJOR UPGRADE"
			case i == 1:
				return "MINOR UPGRADE"
			default:
				return "PATCH"
			}
		}
		if latParts[i] < curParts[i] {
			return "current"
		}
	}
	return "current"
}

func extractNumericVersion(v string) string {
	// Strip leading non-digit characters.
	v = strings.TrimSpace(v)
	// Remove leading v, ^, ~, =, >, <
	v = strings.TrimLeft(v, "v^~=>< ")
	// Take only the numeric version part (before any space or -).
	if idx := strings.IndexAny(v, " -"); idx >= 0 {
		v = v[:idx]
	}
	return v
}

func splitVersion(v string) []int {
	parts := strings.Split(v, ".")
	var nums []int
	for _, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			break
		}
		nums = append(nums, n)
	}
	return nums
}

// ---------- traces command ----------

func tracesCmd(args []string) int {
	fs := flag.NewFlagSet("traces", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() { fmt.Print(tracesUsage()) }

	var (
		jsonOut bool
		verbose bool
		flat    bool
	)
	fs.BoolVar(&jsonOut, "json", false, "JSONL output")
	fs.BoolVar(&verbose, "verbose", false, "Show full tool args")
	fs.BoolVar(&flat, "flat", false, "Disable sub-agent nesting indent")

	flagArgs, positional := splitFlags(args)
	if err := fs.Parse(flagArgs); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}

	if len(positional) == 0 {
		return printErr("Missing mission ID.\n\n" + tracesUsage())
	}
	missionID := strings.TrimSpace(positional[0])

	// Resolve the mission id to validate it exists.
	res := r.Resolve(missionID)
	if res.Safety != "safe" || res.Selected == nil {
		return printErr(resolver.FormatResolutionBlock(res, "traces"))
	}

	// Load traces.
	entries, err := traceStore.LoadTraces(missionID)
	if err != nil {
		return printErr("Failed to load traces:", err)
	}

	if jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		for _, e := range entries {
			enc.Encode(e)
		}
		return 0
	}

	if len(entries) == 0 {
		fmt.Println("No trace events recorded.")
		return 0
	}

	title := res.Selected.Title
	opts := trace.TraceDisplayOptions{Verbose: verbose, Flat: flat}
	output := trace.FormatTraceTable(entries, missionID, title, opts)
	fmt.Print(output)
	return 0
}

func tracesUsage() string {
	return `spacecraft traces <mission-id> [flags]

Flags:
  --json       Raw JSONL output
  --verbose    Show full tool args
  --flat       Disable sub-agent nesting indent
`
}

// ---------- cost command ----------

func costCmd(args []string) int {
	fs := flag.NewFlagSet("cost", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() { fmt.Print(costUsage()) }

	var (
		missionFlag string
		all         bool
	)
	fs.StringVar(&missionFlag, "mission", "", "Show per-model breakdown for one mission")
	fs.BoolVar(&all, "all", false, "Show all missions with trace data")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}

	if missionFlag != "" {
		return costMissionCmd(missionFlag)
	}

	_ = all // --all is the default behavior when no flags

	ids, err := traceStore.ListMissionsWithTraces()
	if err != nil {
		return printErr("Failed to list trace data:", err)
	}
	if len(ids) == 0 {
		fmt.Println("No trace data found.")
		return 0
	}

	var rows []trace.CostRow
	var totalIn, totalOut int
	var totalCost float64

	for _, id := range ids {
		entries, err := traceStore.LoadTraces(id)
		if err != nil || len(entries) == 0 {
			continue
		}
		s := trace.ComputeSummary(entries)
		// Resolve mission title.
		title := id
		if m, lErr := store.Load(id); lErr == nil && m != nil {
			title = fmt.Sprintf("%s — %s", id, m.Title)
			if len(title) > 48 {
				title = title[:45] + "..."
			}
		}
		rows = append(rows, trace.CostRow{
			Mission:   title,
			TokensIn:  s.TotalInTokens,
			TokensOut: s.TotalOutTokens,
			Cost:      s.EstimatedCost,
		})
		totalIn += s.TotalInTokens
		totalOut += s.TotalOutTokens
		totalCost += s.EstimatedCost
	}

	totalLabel := fmt.Sprintf("Total (%d missions)", len(rows))
	total := trace.CostRow{
		Mission:   totalLabel,
		TokensIn:  totalIn,
		TokensOut: totalOut,
		Cost:      totalCost,
	}

	fmt.Print(trace.FormatCostTable(rows, total))
	return 0
}

func costMissionCmd(missionID string) int {
	res := r.Resolve(missionID)
	if res.Safety != "safe" || res.Selected == nil {
		return printErr(resolver.FormatResolutionBlock(res, "cost"))
	}

	entries, err := traceStore.LoadTraces(missionID)
	if err != nil {
		return printErr("Failed to load traces:", err)
	}
	if len(entries) == 0 {
		fmt.Println("No trace events recorded.")
		return 0
	}

	s := trace.ComputeSummary(entries)
	title := fmt.Sprintf("%s — %s", missionID, res.Selected.Title)
	if len(title) > 48 {
		title = title[:45] + "..."
	}

	rows := []trace.CostRow{{
		Mission:   title,
		TokensIn:  s.TotalInTokens,
		TokensOut: s.TotalOutTokens,
		Cost:      s.EstimatedCost,
	}}
	total := trace.CostRow{
		Mission:   fmt.Sprintf("Total"),
		TokensIn:  s.TotalInTokens,
		TokensOut: s.TotalOutTokens,
		Cost:      s.EstimatedCost,
	}

	fmt.Print(trace.FormatCostTable(rows, total))
	fmt.Println()
	fmt.Print(trace.FormatCostBreakdown(entries))
	return 0
}

func costUsage() string {
	return `spacecraft cost [--mission <id>] [--all]

Flags:
  --mission <id>  Show per-model breakdown for one mission
  --all           Show all missions with trace data (default)
`
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

// fatalErr prints the error to stderr and returns exit code 1.
// Use for: missing API key, Brave failure, unknown scope, --deep mode unavailable,
// check-deps errors, and other operational errors.
func fatalErr(v ...interface{}) int {
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

// splitFlags separates flag-like arguments (--*) from positional arguments.
func splitFlags(args []string) (flags, positional []string) {
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			flags = append(flags, a)
		} else {
			positional = append(positional, a)
		}
	}
	return flags, positional
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

// ---------- roadmap commands ----------

func roadmapCmd(args []string) int {
	if len(args) == 0 {
		return printErr("Missing roadmap subcommand.\nUsage: spacecraft roadmap <new|list|add|remove|show|continue|archive>")
	}
	sub := args[0]
	rest := args[1:]
	switch sub {
	case "new":
		return roadmapNewCmd(rest)
	case "list":
		return roadmapListCmd()
	case "add":
		return roadmapAddCmd(rest)
	case "remove":
		return roadmapRemoveCmd(rest)
	case "show":
		return roadmapShowCmd(rest)
	case "continue":
		return roadmapContinueCmd(rest)
	case "archive":
		return roadmapArchiveCmd(rest)
	default:
		return printErr(fmt.Sprintf("Unknown roadmap subcommand %q.\nUsage: spacecraft roadmap <new|list|add|remove|show|continue|archive>", sub))
	}
}

func roadmapNewCmd(args []string) int {
	desc := ""
	positional := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		if args[i] == "--desc" || args[i] == "-d" {
			if i+1 < len(args) {
				desc = args[i+1]
				i++
			}
		} else if !strings.HasPrefix(args[i], "-") {
			positional = append(positional, args[i])
		}
	}
	title := strings.TrimSpace(strings.Join(positional, " "))
	if title == "" {
		return printErr("Missing roadmap title.\nUsage: spacecraft roadmap new <title> [--desc <text>]")
	}

	mid, err := id.MissionId()
	if err != nil {
		return printErr(err)
	}
	now := time.Now()
	r := &roadmap.Roadmap{
		ID:          mid,
		Title:       title,
		Description: desc,
		Missions:    []string{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := roadmapStore.Create(r); err != nil {
		return printErr("Failed to create roadmap:", err)
	}
	fmt.Println(mid)
	return 0
}

func roadmapListCmd() int {
	rms, err := roadmapStore.List()
	if err != nil {
		return printErr(err)
	}
	if len(rms) == 0 {
		fmt.Println("No roadmaps.")
		return 0
	}
	for _, rm := range rms {
		shipped := 0
		for _, mid := range rm.Missions {
			m, err := store.Load(mid)
			if err == nil && (m.State == "shipped" || m.State == "archived") {
				shipped++
			}
		}
		fmt.Printf("%s %s [%d/%d]\n", rm.ID, rm.Title, shipped, len(rm.Missions))
	}
	return 0
}

// stub implementations for future tasks — compile-error-safe
func roadmapAddCmd(args []string) int {
	afterIdx := -1
	positional := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		if args[i] == "--after" {
			if i+1 < len(args) {
				afterIdx = -2 // flag: need to resolve target
				positional = append(positional, args[i+1])
				i++
			}
		} else if !strings.HasPrefix(args[i], "-") {
			positional = append(positional, args[i])
		}
	}
	if len(positional) < 2 {
		return printErr("Usage: spacecraft roadmap add <roadmap-id> <mission-id> [--after <target-mission-id>]")
	}
	rid := positional[0]
	mid := positional[1]
	var targetMid string
	if afterIdx == -2 && len(positional) >= 3 {
		targetMid = positional[2]
		afterIdx = 0
	}

	rm, err := roadmapStore.Load(rid)
	if err != nil {
		return printErr("roadmap not found: " + rid)
	}

	if _, err := store.Load(mid); err != nil {
		return printErr("mission not found: " + mid)
	}

	for _, m := range rm.Missions {
		if m == mid {
			fmt.Println("mission already in roadmap:", mid)
			return 0
		}
	}

	all, _ := roadmapStore.List()
	for _, other := range all {
		if other.ID == rid {
			continue
		}
		for _, m := range other.Missions {
			if m == mid {
				return printErr("already in roadmap: " + mid + " -> " + other.ID)
			}
		}
	}

	if targetMid != "" {
		found := false
		for i, m := range rm.Missions {
			if m == targetMid {
				rm.Missions = append(rm.Missions[:i+1], append([]string{mid}, rm.Missions[i+1:]...)...)
				found = true
				break
			}
		}
		if !found {
			return printErr("target mission not found in roadmap: " + targetMid)
		}
	} else {
		rm.Missions = append(rm.Missions, mid)
	}

	rm.UpdatedAt = time.Now()
	if err := roadmapStore.Save(rm); err != nil {
		return printErr("Failed to save roadmap:", err)
	}
	fmt.Printf("Added %s to roadmap %s\n", mid, rid)
	return 0
}

func roadmapRemoveCmd(args []string) int {
	if len(args) < 2 {
		return printErr("Usage: spacecraft roadmap remove <roadmap-id> <mission-id>")
	}
	rid := args[0]
	mid := args[1]

	rm, err := roadmapStore.Load(rid)
	if err != nil {
		return printErr("roadmap not found: " + rid)
	}

	found := false
	for i, m := range rm.Missions {
		if m == mid {
			rm.Missions = append(rm.Missions[:i], rm.Missions[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return printErr("mission not found in roadmap: " + mid)
	}

	rm.UpdatedAt = time.Now()
	if err := roadmapStore.Save(rm); err != nil {
		return printErr("Failed to save roadmap:", err)
	}
	fmt.Printf("Removed %s from roadmap %s\n", mid, rid)
	return 0
}
func roadmapShowCmd(args []string) int {
	if len(args) < 1 {
		return printErr("Usage: spacecraft roadmap show <roadmap-id>")
	}
	rid := args[0]
	rm, err := roadmapStore.Load(rid)
	if err != nil {
		return printErr("roadmap not found: " + rid)
	}

	fmt.Println(rm.Title)
	if rm.Description != "" {
		fmt.Println(rm.Description)
	}
	fmt.Println()

	shipped := 0
	total := len(rm.Missions)
	shippedStates := map[string]bool{"shipped": true, "archived": true}

	for _, mid := range rm.Missions {
		marker := "[ ]"
		m, err := store.Load(mid)
		if err == nil && shippedStates[m.State] {
			marker = "[x]"
			shipped++
		}
		label := mid
		if m != nil {
			label = fmt.Sprintf("%s %s", mid, m.Title)
		}
		fmt.Printf("%s  %s\n", marker, label)
	}

	fmt.Println()
	fmt.Println("Issues:")
	if len(rm.Issues) == 0 {
		fmt.Println("No issues")
	} else {
		// Group by phase
		phaseOrder := []string{}
		phaseMap := map[string][]roadmap.Issue{}
		for _, issue := range rm.Issues {
			phase := issue.Phase
			if phase == "" {
				phase = "unassigned"
			}
			if _, exists := phaseMap[phase]; !exists {
				phaseOrder = append(phaseOrder, phase)
			}
			phaseMap[phase] = append(phaseMap[phase], issue)
		}
		for _, phase := range phaseOrder {
			fmt.Printf("\n  %s:\n", phase)
			for _, issue := range phaseMap[phase] {
				marker := "[ ]"
				if issue.State == "closed" {
					marker = "[x]"
				}
				labels := ""
				if len(issue.Labels) > 0 {
					labels = " (" + strings.Join(issue.Labels, ", ") + ")"
				}
				fmt.Printf("    %s #%d %s%s\n", marker, issue.Number, issue.Title, labels)
			}
		}
	}

	fmt.Println()
	nextIdx := -1
	for i, mid := range rm.Missions {
		m, err := store.Load(mid)
		if err != nil || !shippedStates[m.State] {
			nextIdx = i
			break
		}
	}

	if total == 0 {
		fmt.Println("No missions in roadmap.")
	} else if shipped == total {
		fmt.Printf("%d/%d done [done]\n", shipped, total)
	} else {
		msg := fmt.Sprintf("%d/%d done", shipped, total)
		if nextIdx >= 0 {
			msg += fmt.Sprintf(" — next: %s", rm.Missions[nextIdx])
		}
		fmt.Println(msg)
	}
	return 0
}
func roadmapContinueCmd(args []string) int {
	if len(args) < 1 {
		return printErr("Usage: spacecraft roadmap continue <roadmap-id>")
	}
	rid := args[0]
	rm, err := roadmapStore.Load(rid)
	if err != nil {
		return printErr("roadmap not found: " + rid)
	}

	if len(rm.Missions) == 0 {
		return printErr("roadmap has no missions")
	}

	shippedStates := map[string]bool{"shipped": true, "archived": true}

	for _, mid := range rm.Missions {
		m, err := store.Load(mid)
		if err != nil || !shippedStates[m.State] {
			if err := store.WriteCurrent(mid); err != nil {
				return printErr("Failed to set current mission:", err)
			}
			label := mid
			if m != nil {
				label = fmt.Sprintf("%s %s", mid, m.Title)
			}
			fmt.Println(label)
			return 0
		}
	}

	fmt.Println("all missions complete")
	return 0
}
func roadmapArchiveCmd(args []string) int {
	if len(args) < 1 {
		return printErr("Usage: spacecraft roadmap archive <roadmap-id>")
	}
	rid := args[0]
	rm, err := roadmapStore.Load(rid)
	if err != nil {
		return printErr("roadmap not found: " + rid)
	}

	shippedStates := map[string]bool{"shipped": true, "archived": true}
	for _, mid := range rm.Missions {
		m, err := store.Load(mid)
		if err != nil || !shippedStates[m.State] {
			return printErr("mission " + mid + " is not shipped — cannot archive roadmap")
		}
	}

	src := filepath.Join(cfg.RoadmapsDir(), rid+".json")
	dst := filepath.Join(cfg.ArchiveDir(), rid+".json")
	if err := os.MkdirAll(cfg.ArchiveDir(), 0755); err != nil {
		return printErr("Failed to ensure archive dir:", err)
	}
	if err := os.Rename(src, dst); err != nil {
		return printErr("Failed to archive roadmap:", err)
	}

	fmt.Printf("Roadmap %s archived\n", rid)
	return 0
}
