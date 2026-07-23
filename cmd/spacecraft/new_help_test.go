package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewHelpFlagShowsUsageNoMission(t *testing.T) {
	for _, flag := range []string{"--help", "-h"} {
		dir := spaceRoot(t)
		res := runCLI(t, dir, "new", flag)
		if res.code != 0 {
			t.Fatalf("new %s exit=%d stderr=%s", flag, res.code, res.stderr)
		}
		if !strings.Contains(res.stdout, "Usage: spacecraft new <title>") {
			t.Fatalf("new %s stdout=%q want usage line", flag, res.stdout)
		}
		entries, err := os.ReadDir(filepath.Join(dir, ".space", "missions"))
		if err != nil {
			t.Fatal(err)
		}
		if len(entries) != 0 {
			t.Fatalf("new %s created mission(s): %v", flag, entries)
		}
	}
}

func TestNewRealTitleStillCreatesMission(t *testing.T) {
	dir := spaceRoot(t)
	res := runCLI(t, dir, "new", "Real Title")
	if res.code != 0 {
		t.Fatalf("new exit=%d stderr=%s", res.code, res.stderr)
	}
	entries, err := os.ReadDir(filepath.Join(dir, ".space", "missions"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("new \"Real Title\" mission count=%d want 1", len(entries))
	}
}

func TestMapNewHelpFlagShowsUsageNoRoadmap(t *testing.T) {
	dir := spaceRoot(t)
	res := runCLI(t, dir, "map", "new", "--help")
	if res.code != 0 {
		t.Fatalf("map new --help exit=%d stderr=%s", res.code, res.stderr)
	}
	if !strings.Contains(res.stdout, "Usage: spacecraft map new <title> [--desc <text>]") {
		t.Fatalf("map new --help stdout=%q want usage line", res.stdout)
	}
	entries, err := os.ReadDir(filepath.Join(dir, ".space", "roadmaps"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("map new --help created roadmap(s): %v", entries)
	}
}
