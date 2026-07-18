package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeReadyCloseoutMission creates a mission that passes closeout when
// SPACECRAFT_CLOSEOUT_SKIP_CHANGELOG=1 is set (good review + evidence with exitCode).
func writeReadyCloseoutMission(t *testing.T, root, id string) {
	t.Helper()
	writeMission(t, root, id, "ready")
	dir := filepath.Join(root, ".space", "missions", id)
	if err := writeCurrent(filepath.Join(root, ".space"), id); err != nil {
		t.Fatal(err)
	}

	evidence := `{"label":"unit","command":"echo hi","output":"hi\n","ts":"2026-01-01T00:00:00Z","exitCode":0}` + "\n"
	if err := os.WriteFile(filepath.Join(dir, "evidence.jsonl"), []byte(evidence), 0644); err != nil {
		t.Fatal(err)
	}

	review := map[string]any{
		"status": "ready",
		"findings": []any{
			map[string]any{"severity": "low", "blocksShip": false, "summary": "nits"},
		},
		"releaseReadiness": map[string]any{
			"changelog": map[string]any{"status": "ready"},
			"specNote":  map[string]any{"status": "ready"},
		},
	}
	rdata, err := json.MarshalIndent(review, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "review.json"), append(rdata, '\n'), 0644); err != nil {
		t.Fatal(err)
	}
}

func writeReviewJSON(t *testing.T, root, id string, review map[string]any) {
	t.Helper()
	dir := filepath.Join(root, ".space", "missions", id)
	data, err := json.MarshalIndent(review, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "review.json"), append(data, '\n'), 0644); err != nil {
		t.Fatal(err)
	}
}

func writeEvidenceLine(t *testing.T, root, id, line string) {
	t.Helper()
	path := filepath.Join(root, ".space", "missions", id, "evidence.jsonl")
	if err := os.WriteFile(path, []byte(line), 0644); err != nil {
		t.Fatal(err)
	}
}

func goodReview() map[string]any {
	return map[string]any{
		"status":   "ready",
		"findings": []any{},
		"releaseReadiness": map[string]any{
			"changelog": map[string]any{"status": "ready"},
			"specNote":  map[string]any{"status": "ready"},
		},
	}
}

func TestCloseoutFailsMissingReviewJSON(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07CLO01"
	writeMission(t, dir, id, "ready")
	if err := writeCurrent(filepath.Join(dir, ".space"), id); err != nil {
		t.Fatal(err)
	}
	writeEvidenceLine(t, dir, id,
		`{"label":"unit","command":"echo","output":"x","ts":"2026-01-01T00:00:00Z","exitCode":0}`+"\n")

	res := runCLIWithEnv(t, dir, []string{"SPACECRAFT_CLOSEOUT_SKIP_CHANGELOG=1"}, "closeout-check")
	if res.code == 0 {
		t.Fatalf("expected fail without review.json\nstdout=%s", res.stdout)
	}
	if !strings.Contains(res.stdout, "Closeout blocked for "+id) {
		t.Fatalf("want blocked header\nstdout=%s", res.stdout)
	}
	if !strings.Contains(res.stdout, "missing review.json") {
		t.Fatalf("want missing review.json\nstdout=%s", res.stdout)
	}
}

func TestCloseoutFailsWrongState(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07CLO02"
	writeReadyCloseoutMission(t, dir, id)
	// Downgrade state to in_progress.
	m, err := readMission(filepath.Join(dir, ".space"), id)
	if err != nil {
		t.Fatal(err)
	}
	m["state"] = "in_progress"
	data, _ := json.MarshalIndent(m, "", "  ")
	if err := os.WriteFile(filepath.Join(dir, ".space", "missions", id, "mission.json"), append(data, '\n'), 0644); err != nil {
		t.Fatal(err)
	}

	res := runCLIWithEnv(t, dir, []string{"SPACECRAFT_CLOSEOUT_SKIP_CHANGELOG=1"}, "closeout-check")
	if res.code == 0 {
		t.Fatalf("expected fail for wrong state\nstdout=%s", res.stdout)
	}
	if !strings.Contains(res.stdout, "state is in_progress") {
		t.Fatalf("want state problem\nstdout=%s", res.stdout)
	}
}

func TestCloseoutFailsClarifyOpen(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07CLO03"
	writeReadyCloseoutMission(t, dir, id)
	if err := os.WriteFile(filepath.Join(dir, ".space", "missions", id, "clarify-status"), []byte("open\n"), 0644); err != nil {
		t.Fatal(err)
	}

	res := runCLIWithEnv(t, dir, []string{"SPACECRAFT_CLOSEOUT_SKIP_CHANGELOG=1"}, "closeout-check")
	if res.code == 0 {
		t.Fatalf("expected fail for clarify open\nstdout=%s", res.stdout)
	}
	if !strings.Contains(res.stdout, "clarify") {
		t.Fatalf("want clarify problem\nstdout=%s", res.stdout)
	}
}

func TestCloseoutFailsEmptyEvidence(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07CLO04"
	writeReadyCloseoutMission(t, dir, id)
	writeEvidenceLine(t, dir, id, "")

	res := runCLIWithEnv(t, dir, []string{"SPACECRAFT_CLOSEOUT_SKIP_CHANGELOG=1"}, "closeout-check")
	if res.code == 0 {
		t.Fatalf("expected fail for empty evidence\nstdout=%s", res.stdout)
	}
	out := res.stdout + res.stderr
	if !strings.Contains(out, "evidence") {
		t.Fatalf("want evidence problem\nstdout=%s", res.stdout)
	}
}

func TestCloseoutFailsMissingExitCode(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07CLO05"
	writeReadyCloseoutMission(t, dir, id)
	writeEvidenceLine(t, dir, id,
		`{"label":"unit","command":"echo","output":"x","ts":"2026-01-01T00:00:00Z"}`+"\n")

	res := runCLIWithEnv(t, dir, []string{"SPACECRAFT_CLOSEOUT_SKIP_CHANGELOG=1"}, "closeout-check")
	if res.code == 0 {
		t.Fatalf("expected fail for missing exitCode\nstdout=%s", res.stdout)
	}
	if !strings.Contains(res.stdout, "exitCode") {
		t.Fatalf("want exitCode problem\nstdout=%s", res.stdout)
	}
}

func TestCloseoutFailsCriticalFinding(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07CLO06"
	writeReadyCloseoutMission(t, dir, id)
	r := goodReview()
	r["findings"] = []any{
		map[string]any{"severity": "critical", "blocksShip": false, "summary": "bad"},
	}
	writeReviewJSON(t, dir, id, r)

	res := runCLIWithEnv(t, dir, []string{"SPACECRAFT_CLOSEOUT_SKIP_CHANGELOG=1"}, "closeout-check")
	if res.code == 0 {
		t.Fatalf("expected fail for critical finding\nstdout=%s", res.stdout)
	}
	if !strings.Contains(strings.ToLower(res.stdout), "critical") {
		t.Fatalf("want critical problem\nstdout=%s", res.stdout)
	}
}

func TestCloseoutFailsBlocksShip(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07CLO07"
	writeReadyCloseoutMission(t, dir, id)
	r := goodReview()
	r["findings"] = []any{
		map[string]any{"severity": "medium", "blocksShip": true, "summary": "blocker"},
	}
	writeReviewJSON(t, dir, id, r)

	res := runCLIWithEnv(t, dir, []string{"SPACECRAFT_CLOSEOUT_SKIP_CHANGELOG=1"}, "closeout-check")
	if res.code == 0 {
		t.Fatalf("expected fail for blocksShip\nstdout=%s", res.stdout)
	}
	if !strings.Contains(res.stdout, "blocksShip") && !strings.Contains(strings.ToLower(res.stdout), "block") {
		t.Fatalf("want blocksShip problem\nstdout=%s", res.stdout)
	}
}

func TestCloseoutFailsDeferredChangelogReadiness(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07CLO08"
	writeReadyCloseoutMission(t, dir, id)
	r := goodReview()
	r["releaseReadiness"] = map[string]any{
		"changelog": map[string]any{"status": "deferred"},
		"specNote":  map[string]any{"status": "ready"},
	}
	writeReviewJSON(t, dir, id, r)

	res := runCLIWithEnv(t, dir, []string{"SPACECRAFT_CLOSEOUT_SKIP_CHANGELOG=1"}, "closeout-check")
	if res.code == 0 {
		t.Fatalf("expected fail for deferred changelog readiness\nstdout=%s", res.stdout)
	}
	if !strings.Contains(strings.ToLower(res.stdout), "changelog") {
		t.Fatalf("want changelog readiness problem\nstdout=%s", res.stdout)
	}
}

func TestCloseoutFailsNoChangelogCommits(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07CLO09"
	writeReadyCloseoutMission(t, dir, id)

	// SKIP unset: non-git temp dir - changelog check must fail (proves SKIP required for pass path).
	res := runCLI(t, dir, "closeout-check")
	if res.code == 0 {
		t.Fatalf("expected fail without CHANGELOG skip\nstdout=%s", res.stdout)
	}
	out := strings.ToUpper(res.stdout + res.stderr)
	if !strings.Contains(out, "CHANGELOG") && !strings.Contains(out, "GIT") {
		t.Fatalf("want CHANGELOG or git problem\nstdout=%s\nstderr=%s", res.stdout, res.stderr)
	}
}

func TestCloseoutPassesFullFixture(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07CLO10"
	writeReadyCloseoutMission(t, dir, id)

	res := runCLIWithEnv(t, dir, []string{"SPACECRAFT_CLOSEOUT_SKIP_CHANGELOG=1"}, "closeout-check")
	if res.code != 0 {
		t.Fatalf("expected pass\nstdout=%s\nstderr=%s", res.stdout, res.stderr)
	}
	if !strings.Contains(res.stdout, "Closeout ready for "+id) {
		t.Fatalf("want ready message\nstdout=%s", res.stdout)
	}
}

func TestShipCheckAliasDispatches(t *testing.T) {
	dir := spaceRoot(t)
	id := "M07CLO11"
	writeReadyCloseoutMission(t, dir, id)

	res := runCLIWithEnv(t, dir, []string{"SPACECRAFT_CLOSEOUT_SKIP_CHANGELOG=1"}, "ship-check")
	if strings.Contains(res.stderr, "unknown command") {
		t.Fatalf("ship-check unknown: %s", res.stderr)
	}
	if res.code != 0 {
		t.Fatalf("ship-check exit=%d\nstdout=%s\nstderr=%s", res.code, res.stdout, res.stderr)
	}
}
