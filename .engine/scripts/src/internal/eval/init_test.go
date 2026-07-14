package eval

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestInitScaffoldsEvalDir(t *testing.T) {
	dir := t.TempDir()
	mid := "M123"

	if err := Init(dir, mid); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	missionDir := filepath.Join(dir, mid)
	for _, name := range []string{"rubric.json", "dataset.json", "config.json"} {
		path := filepath.Join(missionDir, name)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected %s to exist: %v", path, err)
		}
	}
}

func TestInitMissingMissionID(t *testing.T) {
	if err := Init(t.TempDir(), ""); err == nil {
		t.Error("expected error for empty mission ID")
	}
}

func TestInitIdempotent(t *testing.T) {
	dir := t.TempDir()
	mid := "M123"

	if err := Init(dir, mid); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	path := filepath.Join(dir, mid, "config.json")
	if err := os.WriteFile(path, []byte(`{"coverageThreshold":0.1}`), 0644); err != nil {
		t.Fatalf("write custom config: %v", err)
	}

	if err := Init(dir, mid); err != nil {
		t.Fatalf("Init second time failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if string(data) != `{"coverageThreshold":0.1}` {
		t.Error("Init overwrote an existing file")
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	dir := t.TempDir()
	mid := "M123"

	cfg := LoadConfig(dir, mid)
	if cfg.CoverageThreshold != DefaultCoverageThreshold {
		t.Errorf("expected default threshold %v, got %v", DefaultCoverageThreshold, cfg.CoverageThreshold)
	}
}

func TestLoadConfigValidFile(t *testing.T) {
	dir := t.TempDir()
	mid := "M123"
	missionDir := filepath.Join(dir, mid)
	if err := os.MkdirAll(missionDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(missionDir, "config.json"), []byte(`{"coverageThreshold":0.5}`), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg := LoadConfig(dir, mid)
	if cfg.CoverageThreshold != 0.5 {
		t.Errorf("expected threshold 0.5, got %v", cfg.CoverageThreshold)
	}
}

func TestLoadConfigClampsInvalidThreshold(t *testing.T) {
	dir := t.TempDir()
	mid := "M123"
	missionDir := filepath.Join(dir, mid)
	if err := os.MkdirAll(missionDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(missionDir, "config.json"), []byte(`{"coverageThreshold":1.5}`), 0644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg := LoadConfig(dir, mid)
	if cfg.CoverageThreshold != DefaultCoverageThreshold {
		t.Errorf("expected clamped default %v, got %v", DefaultCoverageThreshold, cfg.CoverageThreshold)
	}
}

func TestLoadRubricMissing(t *testing.T) {
	_, err := LoadRubric(t.TempDir(), "M123")
	if err == nil {
		t.Error("expected error loading missing rubric")
	}
}

func TestLoadDatasetMissing(t *testing.T) {
	_, err := LoadDataset(t.TempDir(), "M123")
	if err == nil {
		t.Error("expected error loading missing dataset")
	}
}

func TestCoverage(t *testing.T) {
	if got := Coverage(0, 0); got != 1.0 {
		t.Errorf("Coverage(0,0) = %v, want 1.0", got)
	}
	if got := Coverage(1, 2); got != 0.5 {
		t.Errorf("Coverage(1,2) = %v, want 0.5", got)
	}
	if got := Coverage(3, 0); got != 1.0 {
		t.Errorf("Coverage(3,0) = %v, want 1.0", got)
	}
}

func TestMarshalJSON(t *testing.T) {
	b, err := MarshalJSON(map[string]int{"a": 1})
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}
	if string(b) != "{\n  \"a\": 1\n}" {
		t.Errorf("unexpected output: %s", string(b))
	}
}

func TestDefaultRubricJSON(t *testing.T) {
	b := defaultRubricJSON()
	var rubric EvalRubric
	if err := json.Unmarshal(b, &rubric); err != nil {
		t.Fatalf("unmarshal rubric: %v", err)
	}
	if len(rubric.Dimensions) != len(StdDimensions()) {
		t.Errorf("expected %d dimensions, got %d", len(StdDimensions()), len(rubric.Dimensions))
	}
}

func TestDefaultDatasetJSON(t *testing.T) {
	b := defaultDatasetJSON()
	var dataset EvalDataset
	if err := json.Unmarshal(b, &dataset); err != nil {
		t.Fatalf("unmarshal dataset: %v", err)
	}
	if dataset.Examples == nil {
		t.Error("expected non-nil examples slice")
	}
}

func TestDefaultConfigJSON(t *testing.T) {
	b := defaultConfigJSON()
	var cfg EvalConfig
	if err := json.Unmarshal(b, &cfg); err != nil {
		t.Fatalf("unmarshal config: %v", err)
	}
	if cfg.CoverageThreshold != DefaultCoverageThreshold {
		t.Errorf("expected default threshold, got %v", cfg.CoverageThreshold)
	}
}
