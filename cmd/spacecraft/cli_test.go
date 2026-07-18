package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var spacecraftBin string

func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "spacecraft-bin-")
	if err != nil {
		panic(err)
	}
	spacecraftBin = filepath.Join(tmp, "spacecraft")
	cmd := exec.Command("go", "build", "-o", spacecraftBin, ".")
	cmd.Dir = mustAbs(".")
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmp)
		panic("go build failed: " + err.Error() + "\n" + string(out))
	}
	code := m.Run()
	os.RemoveAll(tmp)
	os.Exit(code)
}

func mustAbs(p string) string {
	a, err := filepath.Abs(p)
	if err != nil {
		panic(err)
	}
	return a
}

type cliResult struct {
	stdout string
	stderr string
	code   int
}

func runCLI(t *testing.T, dir string, args ...string) cliResult {
	t.Helper()
	cmd := exec.Command(spacecraftBin, args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	code := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else {
			t.Fatalf("runCLI %v: %v", args, err)
		}
	}
	return cliResult{stdout: stdout.String(), stderr: stderr.String(), code: code}
}

func spaceRoot(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".space", "missions"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, ".space", "roadmaps"), 0755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func writeMission(t *testing.T, root, id, state string) {
	t.Helper()
	dir := filepath.Join(root, ".space", "missions", id)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	m := map[string]any{
		"id":        id,
		"title":     "Test Mission",
		"state":     state,
		"createdAt": "2026-01-01T00:00:00Z",
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "mission.json"), append(data, '\n'), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "spec.md"), []byte("# Spec\n"), 0644); err != nil {
		t.Fatal(err)
	}
	plan := map[string]any{
		"planName":  "test",
		"missionId": id,
		"tasks":     []any{},
	}
	pdata, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "plan.json"), append(pdata, '\n'), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "evidence.jsonl"), []byte(""), 0644); err != nil {
		t.Fatal(err)
	}
}

func initGitRepo(t *testing.T, dir, branch string) {
	t.Helper()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_CONFIG_NOSYSTEM=1",
			"GIT_TERMINAL_PROMPT=0",
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	// Empty template avoids writing default hooks (works under sandbox).
	run("init", "--template=")
	run("config", "user.email", "test@example.com")
	run("config", "user.name", "Test")
	keep := filepath.Join(dir, ".gitkeep")
	if err := os.WriteFile(keep, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}
	run("add", ".gitkeep")
	run("commit", "-m", "chore: init")
	if branch != "" && branch != "main" && branch != "master" {
		run("checkout", "-b", branch)
	}
}

// Public command surface from plan / main (lean CLI — no research/check-deps stubs).
var restoredCommands = []string{
	"init",
	"new",
	"missions",
	"use",
	"current",
	"resolve",
	"status",
	"flow",
	"bind-branch",
	"git-info",
	"git-suggest",
	"set-state",
	"clarify-status",
	"evidence",
	"validate",
	"closeout-check",
	"archive",
	"roadmap",
	"help",
}

// Removed stubs that must not appear in help or dispatch.
var removedCommands = []string{"research", "check-deps"}

func TestHelpListsRestoredCommands(t *testing.T) {
	dir := spaceRoot(t)
	res := runCLI(t, dir, "help")
	if res.code != 0 {
		t.Fatalf("help exit = %d, want 0\nstderr=%s", res.code, res.stderr)
	}
	out := res.stdout + res.stderr
	var missing []string
	for _, cmd := range restoredCommands {
		if cmd == "help" {
			continue
		}
		// Require the command as a CLI entry (avoid substring hits like "evidence" in prose).
		needle := "spacecraft " + cmd
		if !strings.Contains(out, needle) {
			missing = append(missing, cmd)
		}
	}
	if len(missing) > 0 {
		t.Fatalf("help missing restored commands %v\noutput:\n%s", missing, out)
	}
	var present []string
	for _, cmd := range removedCommands {
		needle := "spacecraft " + cmd
		if strings.Contains(out, needle) {
			present = append(present, cmd)
		}
	}
	if len(present) > 0 {
		t.Fatalf("help must not list removed commands %v\noutput:\n%s", present, out)
	}
}

func TestDispatchAcceptsRestoredCommands(t *testing.T) {
	dir := spaceRoot(t)
	writeMission(t, dir, "M07TEST01", "active")

	tests := []struct {
		name string
		args []string
		// wantUnknown is true when current thin CLI must reject; after restore, exit must not be "unknown command".
		wantAccepted bool
	}{
		{"init", []string{"init"}, true},
		{"new", []string{"new", "Hello Mission"}, true},
		{"missions", []string{"missions"}, true},
		{"current", []string{"current"}, true},
		{"resolve", []string{"resolve", "M07TEST01"}, true},
		{"status", []string{"status"}, true},
		{"flow", []string{"flow"}, true},
		{"set-state", []string{"set-state", "M07TEST01", "planned"}, true},
		{"evidence", []string{"evidence", "--mission", "M07TEST01", "label", "--", "echo", "ok"}, true},
		{"validate", []string{"validate", "M07TEST01"}, true},
		{"roadmap", []string{"roadmap", "ls"}, true},
		{"closeout-check", []string{"closeout-check"}, true},
		{"git-info", []string{"git-info"}, true},
		{"clarify-status", []string{"clarify-status", "clear"}, true},
		{"archive", []string{"archive", "--help"}, true},
		{"bind-branch", []string{"bind-branch", "--help"}, true},
		{"git-suggest", []string{"git-suggest", "--help"}, true},
		{"use", []string{"use", "M07TEST01"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := runCLI(t, dir, tt.args...)
			if strings.Contains(res.stderr, "unknown command") {
				t.Fatalf("%v rejected as unknown command (exit %d)\nstderr=%s", tt.args, res.code, res.stderr)
			}
			if tt.wantAccepted && res.code == 1 && strings.Contains(res.stderr, "unknown") {
				t.Fatalf("%v not dispatched", tt.args)
			}
		})
	}
}

func TestDispatchRejectsRemovedCommands(t *testing.T) {
	dir := spaceRoot(t)

	for _, cmd := range removedCommands {
		t.Run(cmd, func(t *testing.T) {
			res := runCLI(t, dir, cmd)
			if res.code == 0 {
				t.Fatalf("%s must exit non-zero", cmd)
			}
			if !strings.Contains(res.stderr, "unknown command") {
				t.Fatalf("%s must be unknown command\nstderr=%s", cmd, res.stderr)
			}
		})
	}
}

func TestAliasesDispatchToRestoredCommands(t *testing.T) {
	dir := spaceRoot(t)
	writeMission(t, dir, "M07ALIAS1", "active")

	tests := []struct {
		alias     string
		full      string
		aliasArgs []string
	}{
		{"evi", "evidence", []string{"evi", "--mission", "M07ALIAS1", "alias-ok", "--", "echo", "hi"}},
		{"val", "validate", []string{"val", "M07ALIAS1"}},
		{"state", "set-state", []string{"state", "M07ALIAS1", "planned"}},
		{"map", "roadmap", []string{"map", "ls"}},
	}

	for _, tt := range tests {
		t.Run(tt.alias+"→"+tt.full, func(t *testing.T) {
			res := runCLI(t, dir, tt.aliasArgs...)
			if strings.Contains(res.stderr, "unknown command") {
				t.Fatalf("alias %q unknown: %s", tt.alias, res.stderr)
			}
			// Alias must remain wired; restored help should document the full name.
			help := runCLI(t, dir, "help")
			needle := "spacecraft " + tt.full
			if !strings.Contains(help.stdout+help.stderr, needle) {
				t.Fatalf("help must document canonical command %q (alias %q)", tt.full, tt.alias)
			}
		})
	}
}
