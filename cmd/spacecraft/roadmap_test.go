package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRoadmapLifecycle(t *testing.T) {
	dir := spaceRoot(t)
	writeMission(t, dir, "M07MAP01", "active")

	// Canonical command name after restore.
	create := runCLI(t, dir, "roadmap", "new", "Cursor Native", "--desc", "restore cli")
	if stringsContainsUnknown(create.stderr) {
		t.Fatalf("roadmap unknown: %s", create.stderr)
	}
	if create.code != 0 {
		t.Fatalf("roadmap new exit=%d stderr=%s", create.code, create.stderr)
	}

	rid := "cursor-native"
	path := filepath.Join(dir, ".space", "roadmaps", rid+".json")
	if _, err := os.Stat(path); err != nil {
		// Fall back: discover created id from ls/show output or files.
		rid = findRoadmapID(t, dir)
		path = filepath.Join(dir, ".space", "roadmaps", rid+".json")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("roadmap file missing after new: %v\nstdout=%s", err, create.stdout)
	}

	add := runCLI(t, dir, "roadmap", "add", rid, "M07MAP01", "--desc", "first mission")
	if add.code != 0 {
		t.Fatalf("roadmap add exit=%d stderr=%s", add.code, add.stderr)
	}

	show := runCLI(t, dir, "roadmap", "show", rid)
	if show.code != 0 {
		t.Fatalf("roadmap show exit=%d stderr=%s", show.code, show.stderr)
	}
	if !strings.Contains(show.stdout, "M07MAP01") {
		t.Fatalf("show missing mission\n%s", show.stdout)
	}

	ls := runCLI(t, dir, "roadmap", "ls")
	if ls.code != 0 {
		// Accept "list" subcommand as alternate spelling during restore.
		ls = runCLI(t, dir, "roadmap", "list")
	}
	if ls.code != 0 {
		t.Fatalf("roadmap ls/list exit=%d stderr=%s", ls.code, ls.stderr)
	}
	if !strings.Contains(ls.stdout, rid) && !strings.Contains(ls.stdout, "Cursor Native") {
		t.Fatalf("ls missing roadmap\n%s", ls.stdout)
	}

	rm := runCLI(t, dir, "roadmap", "rm", rid, "M07MAP01")
	if rm.code != 0 {
		rm = runCLI(t, dir, "roadmap", "remove", rid, "M07MAP01")
	}
	if rm.code != 0 {
		t.Fatalf("roadmap rm/remove exit=%d stderr=%s", rm.code, rm.stderr)
	}

	show2 := runCLI(t, dir, "roadmap", "show", rid)
	if show2.code != 0 {
		t.Fatalf("roadmap show after rm exit=%d", show2.code)
	}
	var raw map[string]any
	if err := json.Unmarshal([]byte(show2.stdout), &raw); err == nil {
		missions, _ := raw["missions"].([]any)
		if len(missions) != 0 {
			t.Fatalf("missions still present after rm: %v", missions)
		}
	}

	arch := runCLI(t, dir, "roadmap", "archive", rid)
	if arch.code != 0 {
		t.Fatalf("roadmap archive exit=%d stderr=%s", arch.code, arch.stderr)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("roadmap still in roadmaps/ after archive")
	}
}

func TestMapAliasLifecycle(t *testing.T) {
	dir := spaceRoot(t)
	writeMission(t, dir, "M07MAP02", "active")

	res := runCLI(t, dir, "map", "new", "Alias Roadmap")
	if res.code != 0 {
		t.Fatalf("map new exit=%d stderr=%s", res.code, res.stderr)
	}
	rid := "alias-roadmap"
	if _, err := os.Stat(filepath.Join(dir, ".space", "roadmaps", rid+".json")); err != nil {
		rid = findRoadmapID(t, dir)
	}

	add := runCLI(t, dir, "map", "add", rid, "M07MAP02")
	if add.code != 0 {
		t.Fatalf("map add exit=%d stderr=%s", add.code, add.stderr)
	}
	ls := runCLI(t, dir, "map", "ls")
	if ls.code != 0 {
		t.Fatalf("map ls exit=%d stderr=%s", ls.code, ls.stderr)
	}
}

func findRoadmapID(t *testing.T, dir string) string {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(dir, ".space", "roadmaps"))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".json") {
			return strings.TrimSuffix(e.Name(), ".json")
		}
	}
	t.Fatal("no roadmap json found")
	return ""
}
