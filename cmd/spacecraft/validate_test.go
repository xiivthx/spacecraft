package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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

func TestValidateRejectsMismatchedOutputHash(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07VAL10"
	writeMission(t, dir, id, "active")

	evidencePath := filepath.Join(dir, ".space", "missions", id, "evidence.jsonl")
	// Well-formed required fields, but outputHash is not SHA-256 of output.
	bad := `{"label":"unit","command":"echo hi","output":"hi\n","ts":"2026-01-01T00:00:00Z","exitCode":0,"outputHash":"0000000000000000000000000000000000000000000000000000000000000000"}` + "\n"
	if err := os.WriteFile(evidencePath, []byte(bad), 0644); err != nil {
		t.Fatal(err)
	}

	for _, cmd := range []string{"val", "validate"} {
		t.Run(cmd, func(t *testing.T) {
			res := runCLI(t, dir, cmd, id)
			if cmd == "validate" && stringsContainsUnknown(res.stderr) {
				t.Fatalf("validate unknown: %s", res.stderr)
			}
			if res.code == 0 {
				t.Fatalf("%s accepted evidence with mismatched outputHash; want nonzero\nstdout=%s", cmd, res.stdout)
			}
			out := res.stdout + res.stderr
			if !strings.Contains(out, "line 1") {
				t.Fatalf("%s mismatch message must identify evidence line; want %q\nstdout=%s\nstderr=%s",
					cmd, "line 1", res.stdout, res.stderr)
			}
			lower := strings.ToLower(out)
			if !strings.Contains(lower, "outputhash") && !strings.Contains(lower, "hash") {
				t.Fatalf("%s mismatch message must mention outputHash or hash\nstdout=%s\nstderr=%s",
					cmd, res.stdout, res.stderr)
			}
		})
	}
}

func TestValidateAcceptsMatchingOutputHash(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07VAL11"
	writeMission(t, dir, id, "active")

	evidencePath := filepath.Join(dir, ".space", "missions", id, "evidence.jsonl")
	// Independent expected: SHA-256 of "hi\n" (lowercase hex).
	const matchingHash = "98ea6e4f216f2fb4b69fff9b3a44842c38686ca685f3f55dc48c5d3fb1107be4"
	good := `{"label":"unit","command":"echo hi","output":"hi\n","ts":"2026-01-01T00:00:00Z","exitCode":0,"outputHash":"` + matchingHash + `"}` + "\n"
	if err := os.WriteFile(evidencePath, []byte(good), 0644); err != nil {
		t.Fatal(err)
	}

	res := runCLI(t, dir, "val", id)
	if res.code != 0 {
		t.Fatalf("val matching outputHash exit=%d\nstdout=%s\nstderr=%s", res.code, res.stdout, res.stderr)
	}
}

// Legacy evidence.jsonl lines predate outputHash; validate must still accept them.
func TestValidateAcceptsLegacyEvidenceWithoutOutputHash(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07VAL12"
	writeMission(t, dir, id, "active")

	evidencePath := filepath.Join(dir, ".space", "missions", id, "evidence.jsonl")
	legacy := `{"label":"unit","command":"echo hi","output":"hi\n","ts":"2026-01-01T00:00:00Z","exitCode":0}` + "\n"
	if err := os.WriteFile(evidencePath, []byte(legacy), 0644); err != nil {
		t.Fatal(err)
	}

	for _, cmd := range []string{"val", "validate"} {
		t.Run(cmd, func(t *testing.T) {
			res := runCLI(t, dir, cmd, id)
			if cmd == "validate" && stringsContainsUnknown(res.stderr) {
				t.Fatalf("validate unknown: %s", res.stderr)
			}
			if res.code != 0 {
				t.Fatalf("%s must accept well-formed evidence omitting outputHash; exit=%d\nstdout=%s\nstderr=%s",
					cmd, res.code, res.stdout, res.stderr)
			}
		})
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

func TestValidateAcceptsEvidenceWithoutExitCodeNonStrict(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07VAL05"
	writeMission(t, dir, id, "active")

	evidencePath := filepath.Join(dir, ".space", "missions", id, "evidence.jsonl")
	noExit := `{"label":"unit","command":"echo hi","output":"hi\n","ts":"2026-01-01T00:00:00Z"}` + "\n"
	if err := os.WriteFile(evidencePath, []byte(noExit), 0644); err != nil {
		t.Fatal(err)
	}

	res := runCLI(t, dir, "val", id)
	if res.code != 0 {
		t.Fatalf("non-strict must accept evidence without exitCode; exit=%d\nstdout=%s\nstderr=%s",
			res.code, res.stdout, res.stderr)
	}
}

func writePlanWithTasks(t *testing.T, root, id string, tasks []any) {
	t.Helper()
	plan := map[string]any{
		"planName":  "test",
		"missionId": id,
		"tasks":     tasks,
	}
	data, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, ".space", "missions", id, "plan.json")
	if err := os.WriteFile(path, append(data, '\n'), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestValidateStrictFailsEmptyEvidence(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07VAL06"
	writeMission(t, dir, id, "active")

	res := runCLI(t, dir, "validate", "--strict", id)
	if res.code == 0 {
		t.Fatalf("strict must fail empty evidence\nstdout=%s", res.stdout)
	}
	if !strings.Contains(strings.ToLower(res.stdout+res.stderr), "evidence") {
		t.Fatalf("want evidence mention\nstdout=%s\nstderr=%s", res.stdout, res.stderr)
	}
}

func TestValidateStrictFailsMissingExitCode(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07VAL07"
	writeMission(t, dir, id, "active")
	path := filepath.Join(dir, ".space", "missions", id, "evidence.jsonl")
	line := `{"label":"unit","command":"echo","output":"x","ts":"2026-01-01T00:00:00Z"}` + "\n"
	if err := os.WriteFile(path, []byte(line), 0644); err != nil {
		t.Fatal(err)
	}

	res := runCLI(t, dir, "validate", "--strict", id)
	if res.code == 0 {
		t.Fatalf("strict must fail missing exitCode\nstdout=%s", res.stdout)
	}
	if !strings.Contains(res.stdout+res.stderr, "exitCode") {
		t.Fatalf("want exitCode mention\nstdout=%s\nstderr=%s", res.stdout, res.stderr)
	}
}

func TestValidateStrictFailsDoneTaskWithoutMatchingEvidence(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07VAL08"
	writeMission(t, dir, id, "active")
	writePlanWithTasks(t, dir, id, []any{
		map[string]any{
			"id":       "T1",
			"title":    "Do thing",
			"status":   "done",
			"evidence": []string{"t1-pass"},
		},
	})
	path := filepath.Join(dir, ".space", "missions", id, "evidence.jsonl")
	line := `{"label":"other","command":"echo","output":"x","ts":"2026-01-01T00:00:00Z","exitCode":0}` + "\n"
	if err := os.WriteFile(path, []byte(line), 0644); err != nil {
		t.Fatal(err)
	}

	res := runCLI(t, dir, "validate", "--strict", id)
	if res.code == 0 {
		t.Fatalf("strict must fail done task without matching evidence\nstdout=%s", res.stdout)
	}
}

func TestValidateStrictPassesDoneTaskWithMatchingEvidence(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07VAL09"
	writeMission(t, dir, id, "active")
	writePlanWithTasks(t, dir, id, []any{
		map[string]any{
			"id":       "T1",
			"title":    "Do thing",
			"status":   "done",
			"evidence": []string{"t1-pass"},
		},
	})
	path := filepath.Join(dir, ".space", "missions", id, "evidence.jsonl")
	line := `{"label":"t1-pass","command":"echo","output":"x","ts":"2026-01-01T00:00:00Z","exitCode":0}` + "\n"
	if err := os.WriteFile(path, []byte(line), 0644); err != nil {
		t.Fatal(err)
	}

	res := runCLI(t, dir, "validate", "--strict", id)
	if res.code != 0 {
		t.Fatalf("strict pass exit=%d\nstdout=%s\nstderr=%s", res.code, res.stdout, res.stderr)
	}
}
