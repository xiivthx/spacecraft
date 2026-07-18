package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMapNextSkipsReadyReturnsPlanned(t *testing.T) {
	dir := spaceRoot(t)
	writeMission(t, dir, "M07NEXT1", "planned")
	writeMission(t, dir, "M07NEXT2", "ready")

	create := runCLI(t, dir, "map", "new", "Next Order")
	if create.code != 0 {
		t.Fatalf("map new exit=%d stderr=%s", create.code, create.stderr)
	}
	rid := "next-order"
	if _, err := os.Stat(filepath.Join(dir, ".space", "roadmaps", rid+".json")); err != nil {
		rid = findRoadmapID(t, dir)
	}

	// ready first in roadmap order; next must skip it and return planned
	if res := runCLI(t, dir, "map", "add", rid, "M07NEXT2", "--desc", "ready mission"); res.code != 0 {
		t.Fatalf("map add ready exit=%d stderr=%s", res.code, res.stderr)
	}
	if res := runCLI(t, dir, "map", "add", rid, "M07NEXT1", "--desc", "planned mission"); res.code != 0 {
		t.Fatalf("map add planned exit=%d stderr=%s", res.code, res.stderr)
	}

	next := runCLI(t, dir, "map", "next", rid)
	if next.code != 0 {
		t.Fatalf("map next exit=%d stderr=%s", next.code, next.stderr)
	}
	want := "M07NEXT1: planned mission (state=planned)"
	if !strings.Contains(next.stdout, want) {
		t.Fatalf("map next want %q\nstdout=%s", want, next.stdout)
	}
	if strings.Contains(next.stdout, "M07NEXT2") {
		t.Fatalf("map next should skip ready mission\nstdout=%s", next.stdout)
	}
}

func TestMapNextAllComplete(t *testing.T) {
	dir := spaceRoot(t)
	writeMission(t, dir, "M07DONE1", "ready")
	writeMission(t, dir, "M07DONE2", "shipped")

	create := runCLI(t, dir, "map", "new", "All Done")
	if create.code != 0 {
		t.Fatalf("map new exit=%d stderr=%s", create.code, create.stderr)
	}
	rid := "all-done"
	if _, err := os.Stat(filepath.Join(dir, ".space", "roadmaps", rid+".json")); err != nil {
		rid = findRoadmapID(t, dir)
	}

	if res := runCLI(t, dir, "map", "add", rid, "M07DONE1", "--desc", "ready one"); res.code != 0 {
		t.Fatalf("map add exit=%d stderr=%s", res.code, res.stderr)
	}
	if res := runCLI(t, dir, "map", "add", rid, "M07DONE2", "--desc", "shipped one"); res.code != 0 {
		t.Fatalf("map add exit=%d stderr=%s", res.code, res.stderr)
	}

	next := runCLI(t, dir, "map", "next", rid)
	if next.code != 0 {
		t.Fatalf("map next exit=%d stderr=%s", next.code, next.stderr)
	}
	if !strings.Contains(next.stdout, "All missions complete.") {
		t.Fatalf("want All missions complete.\nstdout=%s", next.stdout)
	}
}

func TestMapUseCurrentRoundtrip(t *testing.T) {
	dir := spaceRoot(t)

	create := runCLI(t, dir, "map", "new", "Current Map")
	if create.code != 0 {
		t.Fatalf("map new exit=%d stderr=%s", create.code, create.stderr)
	}
	rid := "current-map"
	if _, err := os.Stat(filepath.Join(dir, ".space", "roadmaps", rid+".json")); err != nil {
		rid = findRoadmapID(t, dir)
	}

	use := runCLI(t, dir, "map", "use", rid)
	if use.code != 0 {
		t.Fatalf("map use exit=%d stderr=%s", use.code, use.stderr)
	}

	data, err := os.ReadFile(filepath.Join(dir, ".space", "current-roadmap"))
	if err != nil {
		t.Fatalf("current-roadmap missing: %v", err)
	}
	if got := strings.TrimSpace(string(data)); got != rid {
		t.Fatalf("current-roadmap=%q want %q", got, rid)
	}

	cur := runCLI(t, dir, "map", "current")
	if cur.code != 0 {
		t.Fatalf("map current exit=%d stderr=%s", cur.code, cur.stderr)
	}
	if got := strings.TrimSpace(cur.stdout); got != rid {
		t.Fatalf("map current=%q want %q", got, rid)
	}
}
