package trace

import (
	"encoding/json"
	"testing"
)

func TestComputeSummaryEmpty(t *testing.T) {
	s := ComputeSummary(nil)
	if s.EventCount != 0 {
		t.Errorf("EventCount = %d, want 0", s.EventCount)
	}
	if s.TotalLatencyMs != 0 {
		t.Errorf("TotalLatencyMs = %d, want 0", s.TotalLatencyMs)
	}
	if s.TotalInTokens != 0 {
		t.Errorf("TotalInTokens = %d, want 0", s.TotalInTokens)
	}
	if s.TotalOutTokens != 0 {
		t.Errorf("TotalOutTokens = %d, want 0", s.TotalOutTokens)
	}
	if s.EstimatedCost != 0 {
		t.Errorf("EstimatedCost = %f, want 0", s.EstimatedCost)
	}
	if s.FirstTS != "" {
		t.Errorf("FirstTS = %q, want empty", s.FirstTS)
	}
	if s.LastTS != "" {
		t.Errorf("LastTS = %q, want empty", s.LastTS)
	}
}

func TestComputeSummaryBasic(t *testing.T) {
	entries := []TraceEntry{
		{
			ID:           "E01",
			Seq:          1,
			TS:           "2026-07-12T10:30:00.000Z",
			Type:         EventToolCall,
			LatencyMs:    2340,
			InputTokens:  0,
			OutputTokens: 0,
		},
		{
			ID:           "E02",
			Seq:          2,
			TS:           "2026-07-12T10:30:05.000Z",
			Type:         EventModelInvoke,
			Model:        strPtr("deepseek-v4-pro"),
			LatencyMs:    1200,
			InputTokens:  4500,
			OutputTokens: 320,
		},
		{
			ID:           "E03",
			Seq:          3,
			TS:           "2026-07-12T10:30:06.000Z",
			Type:         EventCheckpoint,
			StepLabel:    strPtr("T1 complete"),
			LatencyMs:    0,
			InputTokens:  0,
			OutputTokens: 0,
		},
	}

	s := ComputeSummary(entries)

	if s.EventCount != 3 {
		t.Errorf("EventCount = %d, want 3", s.EventCount)
	}
	if s.TotalLatencyMs != 3540 {
		t.Errorf("TotalLatencyMs = %d, want 3540", s.TotalLatencyMs)
	}
	if s.TotalInTokens != 4500 {
		t.Errorf("TotalInTokens = %d, want 4500", s.TotalInTokens)
	}
	if s.TotalOutTokens != 320 {
		t.Errorf("TotalOutTokens = %d, want 320", s.TotalOutTokens)
	}
	if s.FirstTS != "2026-07-12T10:30:00.000Z" {
		t.Errorf("FirstTS = %q", s.FirstTS)
	}
	if s.LastTS != "2026-07-12T10:30:06.000Z" {
		t.Errorf("LastTS = %q", s.LastTS)
	}
	// deepseek-v4-pro: $2.00/1M in, $8.00/1M out
	// 4500/1M * 2.00 + 320/1M * 8.00 = 0.009 + 0.00256 = 0.01156
	expectedCost := (4500.0/1000000.0)*2.00 + (320.0/1000000.0)*8.00
	if s.EstimatedCost < expectedCost-0.0001 || s.EstimatedCost > expectedCost+0.0001 {
		t.Errorf("EstimatedCost = %f, want ~%f", s.EstimatedCost, expectedCost)
	}
}

func TestCalcCostKnownModel(t *testing.T) {
	entries := []TraceEntry{
		{
			Model:        strPtr("glm-5.2"),
			InputTokens:  1000000,
			OutputTokens: 500000,
		},
	}
	cost := CalcCost(entries)
	// $0.80/1M in, $3.20/1M out
	expected := 0.80 + (500000.0/1000000.0)*3.20
	if cost < expected-0.0001 || cost > expected+0.0001 {
		t.Errorf("cost = %f, want ~%f", cost, expected)
	}
}

func TestCalcCostUnknownModel(t *testing.T) {
	entries := []TraceEntry{
		{
			Model:        strPtr("nonexistent-model"),
			InputTokens:  1000000,
			OutputTokens: 1000000,
		},
	}
	cost := CalcCost(entries)
	// unknown fallback: $1.00/1M in, $4.00/1M out = 1.00 + 4.00 = 5.00
	expected := 5.0
	if cost < expected-0.0001 || cost > expected+0.0001 {
		t.Errorf("cost = %f, want ~%f", cost, expected)
	}
}

func TestCalcCostMultipleModels(t *testing.T) {
	entries := []TraceEntry{
		{
			Model:        strPtr("deepseek-v4-pro"),
			InputTokens:  500000,
			OutputTokens: 100000,
		},
		{
			Model:        strPtr("deepseek-v4-flash"),
			InputTokens:  1000000,
			OutputTokens: 500000,
		},
	}
	cost := CalcCost(entries)
	// deepseek-v4-pro: 0.5/1M*$2.00 + 0.1/1M*$8.00 = $1.00 + $0.80 = $1.80
	// deepseek-v4-flash: 1.0/1M*$0.50 + 0.5/1M*$2.00 = $0.50 + $1.00 = $1.50
	// total = $3.30
	expected := 3.30
	if cost < expected-0.0001 || cost > expected+0.0001 {
		t.Errorf("cost = %f, want ~%f", cost, expected)
	}
}

func TestCalcCostNilModel(t *testing.T) {
	entries := []TraceEntry{
		{
			Model:       nil,
			InputTokens: 1000,
		},
	}
	cost := CalcCost(entries)
	// nil model shouldn't panic; treated as unknown fallback
	if cost < 0 {
		t.Errorf("cost should not be negative: %f", cost)
	}
}

func TestCostRoundTrip(t *testing.T) {
	// Round-trip: marshal, unmarshal, compute cost
	entries := []TraceEntry{
		{
			ID:           "E01",
			Seq:          1,
			TS:           "2026-07-12T10:30:00.000Z",
			Type:         EventModelInvoke,
			Model:        strPtr("kimi-k2.7-code"),
			InputTokens:  2000000,
			OutputTokens: 1000000,
			LatencyMs:    5000,
		},
	}

	data, err := json.Marshal(entries[0])
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var entry TraceEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	cost := CalcCost([]TraceEntry{entry})
	// kimi-k2.7-code: $1.50/1M in, $6.00/1M out
	// 2.0 * 1.50 + 1.0 * 6.00 = 9.00
	expected := 9.0
	if cost < expected-0.0001 || cost > expected+0.0001 {
		t.Errorf("cost = %f, want ~%f", cost, expected)
	}
}

func strPtr(s string) *string { return &s }
