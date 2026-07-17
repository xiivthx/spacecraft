package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEvidenceFailingCommandReturnsNonzero(t *testing.T) {
	dir := spaceRoot(t)
	writeMission(t, dir, "M07EVI01", "active")

	res := runCLI(t, dir, "evi", "--mission", "M07EVI01", "fail-case", "--", "sh", "-c", "echo boom; exit 42")
	if res.code == 0 {
		t.Fatalf("evi with failing command returned 0; want nonzero (captured exit)\nstdout=%s\nstderr=%s", res.stdout, res.stderr)
	}
	if res.code != 42 {
		// Prefer exact propagation; nonzero is the minimum contract.
		t.Logf("note: exit=%d (prefer 42 from captured command)", res.code)
	}
}

func TestEvidenceCommandRecordsExitStatusInJSONL(t *testing.T) {
	dir := spaceRoot(t)
	writeMission(t, dir, "M07EVI02", "active")

	// Use evi alias so the capture path runs; assert exit propagation + JSONL schema.
	res := runCLI(t, dir, "evi", "--mission", "M07EVI02", "exit-status", "--", "sh", "-c", "echo out; exit 7")

	path := filepath.Join(dir, ".space", "missions", "M07EVI02", "evidence.jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read evidence.jsonl: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) == 0 || lines[0] == "" {
		t.Fatal("evidence.jsonl empty")
	}
	last := lines[len(lines)-1]

	var entry map[string]any
	if err := json.Unmarshal([]byte(last), &entry); err != nil {
		t.Fatalf("evidence JSONL not JSON: %v\nline=%s", err, last)
	}

	exitVal, ok := entry["exitCode"]
	if !ok {
		t.Errorf("evidence JSONL missing exitCode field\nentry=%s", last)
	} else {
		code, ok := exitVal.(float64)
		if !ok {
			t.Errorf("exitCode has unexpected type %T (%v)", exitVal, exitVal)
		} else if int(code) != 7 {
			t.Errorf("recorded exitCode=%d, want 7", int(code))
		}
	}

	if res.code == 0 {
		t.Errorf("evi with failing command returned 0; want nonzero (captured exit 7)")
	} else if res.code != 7 {
		t.Logf("note: process exit=%d (prefer exact captured exit 7)", res.code)
	}
}

func TestEvidenceSuccessStillZero(t *testing.T) {
	dir := spaceRoot(t)
	writeMission(t, dir, "M07EVI03", "active")

	res := runCLI(t, dir, "evi", "--mission", "M07EVI03", "ok", "--", "echo", "hello")
	if res.code != 0 {
		t.Fatalf("evi success exit=%d, want 0\nstderr=%s", res.code, res.stderr)
	}

	path := filepath.Join(dir, ".space", "missions", "M07EVI03", "evidence.jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var entry map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(string(data))), &entry); err != nil {
		t.Fatal(err)
	}
	if _, ok := entry["exitCode"]; !ok {
		if _, ok := entry["exit"]; !ok {
			t.Fatalf("success evidence must still record exitCode\nentry=%s", data)
		}
	}
}

func TestEvidenceAliasAndCanonicalBothWork(t *testing.T) {
	dir := spaceRoot(t)
	writeMission(t, dir, "M07EVI04", "active")

	evi := runCLI(t, dir, "evi", "--mission", "M07EVI04", "via-evi", "--", "echo", "a")
	if evi.code != 0 {
		t.Fatalf("evi exit=%d stderr=%s", evi.code, evi.stderr)
	}
	ev := runCLI(t, dir, "evidence", "--mission", "M07EVI04", "via-evidence", "--", "echo", "b")
	if strings.Contains(ev.stderr, "unknown command") {
		t.Fatalf("evidence unknown: %s", ev.stderr)
	}
	if ev.code != 0 {
		t.Fatalf("evidence exit=%d stderr=%s", ev.code, ev.stderr)
	}
}
