package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateRejectsMalformedEvidenceJSONL(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07VAL01"
	writeMission(t, dir, id, "active")

	evidencePath := filepath.Join(dir, ".space", "missions", id, "evidence.jsonl")
	malformed := "not-json\n{\"label\":\"ok\",\"command\":\"echo\",\"output\":\"x\",\"ts\":\"2026-01-01T00:00:00Z\"}\n"
	if err := os.WriteFile(evidencePath, []byte(malformed), 0644); err != nil {
		t.Fatal(err)
	}

	for _, cmd := range []string{"val", "validate"} {
		t.Run(cmd, func(t *testing.T) {
			res := runCLI(t, dir, cmd, id)
			if cmd == "validate" && stringsContainsUnknown(res.stderr) {
				t.Fatalf("validate unknown: %s", res.stderr)
			}
			if res.code == 0 {
				t.Fatalf("%s accepted malformed evidence.jsonl; want nonzero\nstdout=%s", cmd, res.stdout)
			}
		})
	}
}

func TestValidateRejectsEvidenceMissingRequiredFields(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07VAL02"
	writeMission(t, dir, id, "active")

	evidencePath := filepath.Join(dir, ".space", "missions", id, "evidence.jsonl")
	// Valid JSON object but missing required evidence fields.
	bad := `{"foo":"bar"}` + "\n"
	if err := os.WriteFile(evidencePath, []byte(bad), 0644); err != nil {
		t.Fatal(err)
	}

	res := runCLI(t, dir, "val", id)
	if res.code == 0 {
		t.Fatalf("val accepted evidence entry missing label/command/output/ts; want nonzero\nstdout=%s", res.stdout)
	}
}

func TestValidateAcceptsWellFormedMission(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07VAL03"
	writeMission(t, dir, id, "active")

	evidencePath := filepath.Join(dir, ".space", "missions", id, "evidence.jsonl")
	good := `{"label":"unit","command":"echo hi","output":"hi\n","ts":"2026-01-01T00:00:00Z","exitCode":0}` + "\n"
	if err := os.WriteFile(evidencePath, []byte(good), 0644); err != nil {
		t.Fatal(err)
	}

	res := runCLI(t, dir, "val", id)
	if res.code != 0 {
		t.Fatalf("val well-formed exit=%d\nstdout=%s\nstderr=%s", res.code, res.stdout, res.stderr)
	}
}

func TestValidateRejectsMissingArtifacts(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07VAL04"
	writeMission(t, dir, id, "active")
	os.Remove(filepath.Join(dir, ".space", "missions", id, "spec.md"))

	res := runCLI(t, dir, "val", id)
	if res.code == 0 {
		t.Fatal("val must fail when spec.md missing")
	}
}
