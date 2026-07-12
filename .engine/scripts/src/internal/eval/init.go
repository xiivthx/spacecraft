package eval

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"spacecraft/internal/util"
)

// Init scaffolds an eval directory for a mission.
// It is idempotent: existing files are not overwritten.
func Init(evalsDir, missionID string) error {
	if missionID == "" {
		return errors.New("eval init: missing mission ID")
	}

	dir := filepath.Join(evalsDir, missionID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("eval init: create directory: %w", err)
	}

	if err := scaffoldIfAbsent(filepath.Join(dir, "rubric.json"), defaultRubricJSON()); err != nil {
		return err
	}
	if err := scaffoldIfAbsent(filepath.Join(dir, "dataset.json"), defaultDatasetJSON()); err != nil {
		return err
	}
	if err := scaffoldIfAbsent(filepath.Join(dir, "config.json"), defaultConfigJSON()); err != nil {
		return err
	}

	fmt.Printf("Scaffolded eval directory for %s: %s\n", missionID, dir)
	return nil
}

// LoadConfig reads the eval config for a mission, falling back to defaults.
func LoadConfig(evalsDir, missionID string) EvalConfig {
	path := filepath.Join(evalsDir, missionID, "config.json")
	var cfg EvalConfig
	if err := util.ReadJson(path, &cfg); err != nil {
		return DefaultEvalConfig()
	}
	if cfg.CoverageThreshold <= 0 || cfg.CoverageThreshold > 1.0 {
		cfg.CoverageThreshold = DefaultCoverageThreshold
	}
	return cfg
}

// LoadRubric reads the rubric for a mission.
func LoadRubric(evalsDir, missionID string) (*EvalRubric, error) {
	path := filepath.Join(evalsDir, missionID, "rubric.json")
	var rubric EvalRubric
	if err := util.ReadJson(path, &rubric); err != nil {
		return nil, fmt.Errorf("eval: load rubric: %w", err)
	}
	return &rubric, nil
}

// LoadDataset reads the eval dataset for a mission.
func LoadDataset(evalsDir, missionID string) (*EvalDataset, error) {
	path := filepath.Join(evalsDir, missionID, "dataset.json")
	var dataset EvalDataset
	if err := util.ReadJson(path, &dataset); err != nil {
		return nil, fmt.Errorf("eval: load dataset: %w", err)
	}
	return &dataset, nil
}

// Coverage calculates eval coverage ratio: covered_checks / total_checks.
// datasetSize is the number of labelled examples in the dataset.
// planTaskCount is the number of tasks in the plan.
func Coverage(datasetSize, planTaskCount int) float64 {
	if planTaskCount == 0 {
		return 1.0
	}
	return float64(datasetSize) / float64(planTaskCount)
}

func scaffoldIfAbsent(path string, content []byte) error {
	if util.Exists(path) {
		return nil
	}
	return os.WriteFile(path, content, 0644)
}

func defaultRubricJSON() []byte {
	rubric := EvalRubric{
		Dimensions: StdDimensions(),
	}
	b, _ := marshalJSON(rubric)
	return b
}

func defaultDatasetJSON() []byte {
	dataset := EvalDataset{
		Examples: []EvalDatasetEntry{},
	}
	b, _ := marshalJSON(dataset)
	return b
}

func defaultConfigJSON() []byte {
	cfg := DefaultEvalConfig()
	b, _ := marshalJSON(cfg)
	return b
}

func marshalJSON(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

// MarshalJSON marshals v to indented JSON.
func MarshalJSON(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}
