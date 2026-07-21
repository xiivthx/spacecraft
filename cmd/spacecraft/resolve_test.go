package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestResolveFromFeatBranch covers T2: with .space/current unset and branch
// feat/<B>/… where B exists, resolve prints Mission: B.
func TestResolveFromFeatBranch(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07RES01"
	writeMission(t, dir, id, "active")
	initGitRepo(t, dir, "feat/"+id+"/restore-cli")

	res := runCLI(t, dir, "resolve")
	if stringsContainsUnknown(res.stderr) {
		t.Fatalf("resolve unknown command: %s", res.stderr)
	}
	if res.code != 0 {
		t.Fatalf("resolve exit=%d stderr=%s stdout=%s", res.code, res.stderr, res.stdout)
	}
	want := "Mission: " + id
	out := res.stdout + res.stderr
	if !strings.Contains(out, want) {
		t.Fatalf("resolve: want %q (branch when current unset), got:\n%s", want, out)
	}
}

// TestResolveExplicitSelector covers T5: resolve <C> wins over .space/current=A
// and branch feat/<B>/….
func TestResolveExplicitSelector(t *testing.T) {
	dir, currentID, branchID := currentOverridesBranchDir(t)
	explicitID := "M07EXPC1"
	writeMission(t, dir, explicitID, "active")

	res := runCLI(t, dir, "resolve", explicitID)
	if stringsContainsUnknown(res.stderr) {
		t.Fatalf("resolve unknown: %s", res.stderr)
	}
	if res.code != 0 {
		t.Fatalf("resolve %s exit=%d stderr=%s", explicitID, res.code, res.stderr)
	}
	out := res.stdout + res.stderr
	want := "Mission: " + explicitID
	if !strings.Contains(out, want) {
		t.Fatalf("resolve: want %q (explicit over current/branch), got:\n%s", want, out)
	}
	if strings.Contains(out, "Mission: "+currentID) {
		t.Fatalf("resolve: current %s must not win over explicit %s\n%s", currentID, explicitID, out)
	}
	if strings.Contains(out, "Mission: "+branchID) {
		t.Fatalf("resolve: branch %s must not win over explicit %s\n%s", branchID, explicitID, out)
	}
}

func TestResolveUnknownSelectorFails(t *testing.T) {
	dir := spaceRoot(t)
	writeMission(t, dir, "M07RES03", "active")
	initGitRepo(t, dir, "main")

	res := runCLI(t, dir, "resolve", "M07DOESNOTEXIST")
	if stringsContainsUnknown(res.stderr) {
		t.Fatalf("resolve unknown command: %s", res.stderr)
	}
	if res.code == 0 {
		t.Fatal("resolve bad selector must fail")
	}
}

func TestEvidenceUsesBranchMissionWhenNoFlag(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07RES04"
	writeMission(t, dir, id, "active")
	initGitRepo(t, dir, "feat/"+id+"/branch-mission")

	res := runCLI(t, dir, "evi", "from-branch", "--", "echo", "branched")
	if res.code != 0 {
		t.Fatalf("evi without --mission on feat branch exit=%d stderr=%s", res.code, res.stderr)
	}
	path := filepath.Join(dir, ".space", "missions", id, "evidence.jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "from-branch") {
		t.Fatalf("evidence not appended to branch mission\n%s", data)
	}
}

func TestEvidenceUsesCurrentOverBranch(t *testing.T) {
	dir, currentID, branchID := currentOverridesBranchDir(t)

	res := runCLI(t, dir, "evi", "from-current", "--", "echo", "currented")
	if res.code != 0 {
		t.Fatalf("evi without --mission exit=%d stderr=%s", res.code, res.stderr)
	}

	currentPath := filepath.Join(dir, ".space", "missions", currentID, "evidence.jsonl")
	data, err := os.ReadFile(currentPath)
	if err != nil {
		t.Fatalf("evidence should append to current mission %s: %v", currentID, err)
	}
	if !strings.Contains(string(data), "from-current") {
		t.Fatalf("evidence not appended to current mission %s\n%s", currentID, data)
	}

	branchPath := filepath.Join(dir, ".space", "missions", branchID, "evidence.jsonl")
	if branchData, err := os.ReadFile(branchPath); err == nil && strings.Contains(string(branchData), "from-current") {
		t.Fatalf("evidence must not append to branch mission %s when current=%s\n%s", branchID, currentID, branchData)
	}
}

// TestEvidenceMissionFlagOverridesCurrentAndBranch covers T5: evi --mission <C>
// writes to C even when current=A and branch=feat/<B>/….
// RED skipped: --mission already bypasses resolveActive after T4 (regression lock).
func TestEvidenceMissionFlagOverridesCurrentAndBranch(t *testing.T) {
	dir, currentID, branchID := currentOverridesBranchDir(t)
	explicitID := "M07EXPC2"
	writeMission(t, dir, explicitID, "active")

	res := runCLI(t, dir, "evi", "--mission", explicitID, "from-flag", "--", "echo", "flagged")
	if res.code != 0 {
		t.Fatalf("evi --mission exit=%d stderr=%s", res.code, res.stderr)
	}

	flagPath := filepath.Join(dir, ".space", "missions", explicitID, "evidence.jsonl")
	data, err := os.ReadFile(flagPath)
	if err != nil {
		t.Fatalf("evidence should append to --mission target %s: %v", explicitID, err)
	}
	if !strings.Contains(string(data), "from-flag") {
		t.Fatalf("evidence not appended to --mission mission %s\n%s", explicitID, data)
	}

	for _, otherID := range []string{currentID, branchID} {
		otherPath := filepath.Join(dir, ".space", "missions", otherID, "evidence.jsonl")
		if otherData, err := os.ReadFile(otherPath); err == nil && strings.Contains(string(otherData), "from-flag") {
			t.Fatalf("evidence must not append to %s when --mission=%s\n%s", otherID, explicitID, otherData)
		}
	}
}

// currentOverridesBranchDir sets .space/current to mission A while the git
// branch is feat/<B>/… (both missions exist). Used by T1 priority tests.
func currentOverridesBranchDir(t *testing.T) (dir, currentID, branchID string) {
	t.Helper()
	dir = spaceRoot(t)
	currentID = "M07CURA1"
	branchID = "M07CURB1"
	writeMission(t, dir, currentID, "active")
	writeMission(t, dir, branchID, "active")
	if err := writeCurrent(filepath.Join(dir, ".space"), currentID); err != nil {
		t.Fatal(err)
	}
	initGitRepo(t, dir, "feat/"+branchID+"/other-mission")
	return dir, currentID, branchID
}

func assertMissionLine(t *testing.T, cmd string, res cliResult, wantID, notID string) {
	t.Helper()
	if res.code != 0 {
		t.Fatalf("%s exit=%d stderr=%s stdout=%s", cmd, res.code, res.stderr, res.stdout)
	}
	out := res.stdout + res.stderr
	want := "Mission: " + wantID
	if !strings.Contains(out, want) {
		t.Fatalf("%s: want %q (current over branch), got:\n%s", cmd, want, out)
	}
	if strings.Contains(out, "Mission: "+notID) {
		t.Fatalf("%s: branch mission %s must not win over current %s\n%s", cmd, notID, wantID, out)
	}
}

func TestResolveCurrentOverridesBranch(t *testing.T) {
	dir, currentID, branchID := currentOverridesBranchDir(t)
	res := runCLI(t, dir, "resolve")
	assertMissionLine(t, "resolve", res, currentID, branchID)
}

func TestStatusCurrentOverridesBranch(t *testing.T) {
	dir, currentID, branchID := currentOverridesBranchDir(t)
	res := runCLI(t, dir, "status")
	assertMissionLine(t, "status", res, currentID, branchID)
}

func TestFlowCurrentOverridesBranch(t *testing.T) {
	dir, currentID, branchID := currentOverridesBranchDir(t)
	res := runCLI(t, dir, "flow")
	assertMissionLine(t, "flow", res, currentID, branchID)
}

// TestResolveMissingCurrentFallsThroughToBranch covers T3: .space/current
// points at a missing mission id, so resolve falls through to feat/<B>/….
func TestResolveMissingCurrentFallsThroughToBranch(t *testing.T) {
	dir := spaceRoot(t)
	branchID := "M07CURB2"
	writeMission(t, dir, branchID, "active")
	if err := writeCurrent(filepath.Join(dir, ".space"), "M07MISSING"); err != nil {
		t.Fatal(err)
	}
	initGitRepo(t, dir, "feat/"+branchID+"/other-mission")

	res := runCLI(t, dir, "resolve")
	if res.code != 0 {
		t.Fatalf("resolve exit=%d stderr=%s stdout=%s", res.code, res.stderr, res.stdout)
	}
	want := "Mission: " + branchID
	out := res.stdout + res.stderr
	if !strings.Contains(out, want) {
		t.Fatalf("resolve: want %q (branch when current missing), got:\n%s", want, out)
	}
	if strings.Contains(out, "Mission: M07MISSING") {
		t.Fatalf("resolve must not print missing current id\n%s", out)
	}
}
