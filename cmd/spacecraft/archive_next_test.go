package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestArchiveSuggestsNextRoadmapMission(t *testing.T) {
	dir := spaceRoot(t)
	space := filepath.Join(dir, ".space")
	writeMission(t, dir, "M07ARCH1", "shipped")
	writeMission(t, dir, "M07ARCH2", "planned")
	if err := writeCurrent(space, "M07ARCH1"); err != nil {
		t.Fatal(err)
	}

	create := runCLI(t, dir, "map", "new", "Archive Next")
	if create.code != 0 {
		t.Fatalf("map new exit=%d stderr=%s", create.code, create.stderr)
	}
	rid := "archive-next"
	if _, err := os.Stat(filepath.Join(space, "roadmaps", rid+".json")); err != nil {
		rid = findRoadmapID(t, dir)
	}
	if res := runCLI(t, dir, "map", "add", rid, "M07ARCH1", "--desc", "shipped one"); res.code != 0 {
		t.Fatalf("map add exit=%d stderr=%s", res.code, res.stderr)
	}
	if res := runCLI(t, dir, "map", "add", rid, "M07ARCH2", "--desc", "planned next"); res.code != 0 {
		t.Fatalf("map add exit=%d stderr=%s", res.code, res.stderr)
	}
	if res := runCLI(t, dir, "map", "use", rid); res.code != 0 {
		t.Fatalf("map use exit=%d stderr=%s", res.code, res.stderr)
	}

	arch := runCLI(t, dir, "archive", "M07ARCH1")
	if arch.code != 0 {
		t.Fatalf("archive exit=%d stderr=%s", arch.code, arch.stderr)
	}
	if !strings.Contains(arch.stdout, "Archived mission M07ARCH1") {
		t.Fatalf("want archived line\nstdout=%s", arch.stdout)
	}
	wantHint := "Next mission on roadmap " + rid + ": M07ARCH2: planned next (state=planned)"
	if !strings.Contains(arch.stdout, wantHint) {
		t.Fatalf("want hint %q\nstdout=%s", wantHint, arch.stdout)
	}
	if !strings.Contains(arch.stdout, "Suggested: new session → /sc-discuss M07ARCH2 (then /sc-run)") {
		t.Fatalf("want discuss suggestion\nstdout=%s", arch.stdout)
	}
	if got := readCurrent(space); got != "M07ARCH2" {
		t.Fatalf("current=%q want M07ARCH2", got)
	}
}

func TestArchiveLastMissionNoNextHint(t *testing.T) {
	dir := spaceRoot(t)
	space := filepath.Join(dir, ".space")
	writeMission(t, dir, "M07LAST1", "shipped")
	if err := writeCurrent(space, "M07LAST1"); err != nil {
		t.Fatal(err)
	}

	create := runCLI(t, dir, "map", "new", "Archive Last")
	if create.code != 0 {
		t.Fatalf("map new exit=%d stderr=%s", create.code, create.stderr)
	}
	rid := "archive-last"
	if _, err := os.Stat(filepath.Join(space, "roadmaps", rid+".json")); err != nil {
		rid = findRoadmapID(t, dir)
	}
	if res := runCLI(t, dir, "map", "add", rid, "M07LAST1", "--desc", "only one"); res.code != 0 {
		t.Fatalf("map add exit=%d stderr=%s", res.code, res.stderr)
	}

	arch := runCLI(t, dir, "archive", "M07LAST1")
	if arch.code != 0 {
		t.Fatalf("archive exit=%d stderr=%s", arch.code, arch.stderr)
	}
	if strings.Contains(arch.stdout, "Next mission on roadmap") {
		t.Fatalf("did not expect next hint\nstdout=%s", arch.stdout)
	}
	if got := readCurrent(space); got != "" {
		t.Fatalf("current=%q want cleared", got)
	}
}
