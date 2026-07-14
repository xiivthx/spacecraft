package eval

import "testing"

func TestDefaultEvalConfig(t *testing.T) {
	cfg := DefaultEvalConfig()
	if cfg.CoverageThreshold != DefaultCoverageThreshold {
		t.Errorf("expected coverage threshold %v, got %v", DefaultCoverageThreshold, cfg.CoverageThreshold)
	}
}

func TestStdDimensionsCount(t *testing.T) {
	dims := StdDimensions()
	if len(dims) != 5 {
		t.Errorf("expected 5 standard dimensions, got %d", len(dims))
	}
}
