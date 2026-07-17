package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
	out := res.stdout + res.stderr
	if !strings.Contains(out, id) {
		t.Fatalf("resolve output missing mission id %s\n%s", id, out)
	}
}

func TestResolveExplicitSelector(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07RES02"
	writeMission(t, dir, id, "active")
	initGitRepo(t, dir, "main")

	res := runCLI(t, dir, "resolve", id)
	if stringsContainsUnknown(res.stderr) {
		t.Fatalf("resolve unknown: %s", res.stderr)
	}
	if res.code != 0 {
		t.Fatalf("resolve %s exit=%d stderr=%s", id, res.code, res.stderr)
	}
	if !strings.Contains(res.stdout+res.stderr, id) {
		t.Fatalf("resolve output missing %s\n%s", id, res.stdout)
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
