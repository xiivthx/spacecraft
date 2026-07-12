package trace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"spacecraft/internal/config"
)

func setupStore(t *testing.T) (TraceStore, string) {
	t.Helper()
	dir := t.TempDir()
	cfg, err := config.NewConfig(dir, config.WithTraceStoreDir(filepath.Join(dir, ".space", "traces")))
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}
	return NewFSTraceStore(cfg), dir
}

func writeFile(t *testing.T, dir, relPath, content string) {
	t.Helper()
	fullPath := filepath.Join(dir, relPath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
}

func TestLoadTracesEmptyDir(t *testing.T) {
	store, _ := setupStore(t)
	entries, err := store.LoadTraces("M07N6P7I4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestLoadTracesValidJSONL(t *testing.T) {
	store, dir := setupStore(t)
	traceDir := filepath.Join(dir, ".space", "traces")
	writeFile(t, traceDir, "M07N6P7I4.jsonl", `{"id":"E01","seq":1,"ts":"2026-07-12T10:30:00.000Z","type":"tool_call","tool":"bash","args":{"command":"go build"},"latencyMs":2340,"inputTokens":0,"outputTokens":0}
`)

	entries, err := store.LoadTraces("M07N6P7I4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.ID != "E01" {
		t.Errorf("ID = %q, want %q", e.ID, "E01")
	}
	if e.Seq != 1 {
		t.Errorf("Seq = %d, want %d", e.Seq, 1)
	}
	if e.LatencyMs != 2340 {
		t.Errorf("LatencyMs = %d, want %d", e.LatencyMs, 2340)
	}
	if e.Type != EventToolCall {
		t.Errorf("Type = %q, want %q", e.Type, EventToolCall)
	}
	if e.Tool == nil || *e.Tool != "bash" {
		t.Errorf("Tool = %v, want bash", e.Tool)
	}
}

func TestLoadTracesSkipsEmptyLines(t *testing.T) {
	store, dir := setupStore(t)
	traceDir := filepath.Join(dir, ".space", "traces")
	writeFile(t, traceDir, "M07N6P7I4.jsonl", `
{"id":"E01","seq":1,"ts":"2026-07-12T10:30:00.000Z","type":"checkpoint","latencyMs":0,"inputTokens":0,"outputTokens":0}

{"id":"E02","seq":2,"ts":"2026-07-12T10:31:00.000Z","type":"checkpoint","latencyMs":0,"inputTokens":0,"outputTokens":0}

`)

	entries, err := store.LoadTraces("M07N6P7I4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].ID != "E01" {
		t.Errorf("first ID = %q, want E01", entries[0].ID)
	}
	if entries[1].ID != "E02" {
		t.Errorf("second ID = %q, want E02", entries[1].ID)
	}
}

func TestHasTracesMissing(t *testing.T) {
	store, _ := setupStore(t)
	if store.HasTraces("M07N6P7I4") {
		t.Error("HasTraces should return false for missing file")
	}
}

func TestHasTracesEmptyFile(t *testing.T) {
	store, dir := setupStore(t)
	traceDir := filepath.Join(dir, ".space", "traces")
	writeFile(t, traceDir, "M07N6P7I4.jsonl", "\n\n\n")

	if store.HasTraces("M07N6P7I4") {
		t.Error("HasTraces should return false for file with only empty lines")
	}
}

func TestHasTracesWithEntries(t *testing.T) {
	store, dir := setupStore(t)
	traceDir := filepath.Join(dir, ".space", "traces")
	writeFile(t, traceDir, "M07N6P7I4.jsonl", `{"id":"E01","seq":1,"ts":"2026-07-12T10:30:00.000Z","type":"checkpoint","latencyMs":0,"inputTokens":0,"outputTokens":0}
`)

	if !store.HasTraces("M07N6P7I4") {
		t.Error("HasTraces should return true for file with entries")
	}
}

func TestListMissionsWithTracesEmpty(t *testing.T) {
	store, _ := setupStore(t)
	ids, err := store.ListMissionsWithTraces()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected 0 ids, got %d", len(ids))
	}
}

func TestListMissionsWithTraces(t *testing.T) {
	store, dir := setupStore(t)
	traceDir := filepath.Join(dir, ".space", "traces")
	writeFile(t, traceDir, "M07AAAAAA.jsonl", `{"id":"E01","seq":1,"ts":"2026-07-12T10:30:00.000Z","type":"checkpoint","latencyMs":0,"inputTokens":0,"outputTokens":0}`)
	writeFile(t, traceDir, "M07BBBBBB.jsonl", `{"id":"E01","seq":1,"ts":"2026-07-12T10:30:00.000Z","type":"checkpoint","latencyMs":0,"inputTokens":0,"outputTokens":0}`)
	writeFile(t, traceDir, "not-jsonl.txt", "some text")

	ids, err := store.ListMissionsWithTraces()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sort.Strings(ids)
	if len(ids) != 2 {
		t.Fatalf("expected 2 ids, got %d: %v", len(ids), ids)
	}
	if ids[0] != "M07AAAAAA" {
		t.Errorf("ids[0] = %q, want M07AAAAAA", ids[0])
	}
	if ids[1] != "M07BBBBBB" {
		t.Errorf("ids[1] = %q, want M07BBBBBB", ids[1])
	}
}

func TestLoadTracesArgsPreserved(t *testing.T) {
	store, dir := setupStore(t)
	traceDir := filepath.Join(dir, ".space", "traces")
	writeFile(t, traceDir, "M07N6P7I4.jsonl", `{"id":"E01","seq":1,"ts":"2026-07-12T10:30:00.000Z","type":"tool_call","tool":"bash","args":{"command":"echo hello","workdir":"/tmp"},"latencyMs":500,"inputTokens":0,"outputTokens":0}
`)

	entries, err := store.LoadTraces("M07N6P7I4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Args == nil {
		t.Fatal("Args should not be nil")
	}
	var args map[string]string
	if err := json.Unmarshal(entries[0].Args, &args); err != nil {
		t.Fatalf("failed to unmarshal args: %v", err)
	}
	if args["command"] != "echo hello" {
		t.Errorf("args[command] = %q, want echo hello", args["command"])
	}
	if args["workdir"] != "/tmp" {
		t.Errorf("args[workdir] = %q, want /tmp", args["workdir"])
	}
}
