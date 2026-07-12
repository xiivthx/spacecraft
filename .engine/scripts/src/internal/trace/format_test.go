package trace

import (
	"strings"
	"testing"
)

func TestFormatTraceTableBasic(t *testing.T) {
	entries := []TraceEntry{
		{
			ID:           "E01",
			Seq:          1,
			TS:           "2026-07-12T10:30:00.123Z",
			Type:         EventToolCall,
			Tool:         strPtr("bash"),
			LatencyMs:    2340,
			ExitCode:     intPtr(0),
			InputTokens:  0,
			OutputTokens: 0,
		},
		{
			ID:           "E02",
			Seq:          2,
			TS:           "2026-07-12T10:30:02.500Z",
			Type:         EventModelInvoke,
			Model:        strPtr("glm-5.2"),
			LatencyMs:    1200,
			InputTokens:  4500,
			OutputTokens: 320,
		},
		{
			ID:           "E03",
			Seq:          3,
			TS:           "2026-07-12T10:30:45.200Z",
			Type:         EventCheckpoint,
			StepLabel:    strPtr("T1 complete"),
			LatencyMs:    0,
			InputTokens:  0,
			OutputTokens: 0,
		},
	}

	opts := TraceDisplayOptions{}
	result := FormatTraceTable(entries, "M07N6P7I4", "Observability: token tracking", opts)

	if !strings.Contains(result, "=== Trace: M07N6P7I4") {
		t.Error("missing header with mission ID")
	}
	if !strings.Contains(result, "T1 complete") {
		t.Error("missing step label")
	}
	if !strings.Contains(result, "model_invoke") {
		t.Error("missing model_invoke type")
	}
	if !strings.Contains(result, "Summary:") {
		t.Error("missing summary line")
	}
}

func TestFormatCostTable(t *testing.T) {
	rows := []CostRow{
		{Mission: "M07N6P7I4 — Observability...", TokensIn: 16500, TokensOut: 2720, Cost: 0.042},
		{Mission: "M07FYB5W5 — Another mission", TokensIn: 45000, TokensOut: 8100, Cost: 0.115},
	}
	total := CostRow{Mission: "Total (2 missions)", TokensIn: 61500, TokensOut: 10820, Cost: 0.157}

	result := FormatCostTable(rows, total)

	if !strings.Contains(result, "M07N6P7I4 — Observability...") {
		t.Error("missing first mission row")
	}
	if !strings.Contains(result, "M07FYB5W5 — Another mission") {
		t.Error("missing second mission row")
	}
	if !strings.Contains(result, "Total (2 missions)") {
		t.Error("missing total row")
	}
	if !strings.Contains(result, "$0.042") {
		t.Error("missing cost for first mission")
	}
	if !strings.Contains(result, "─") {
		t.Error("missing separator line")
	}
}

func TestFormatCostTableSingleMission(t *testing.T) {
	rows := []CostRow{
		{Mission: "M07N6P7I4 — Test", TokensIn: 1000, TokensOut: 500, Cost: 0.003},
	}
	total := CostRow{Mission: "Total (1 mission)", TokensIn: 1000, TokensOut: 500, Cost: 0.003}

	result := FormatCostTable(rows, total)
	if !strings.Contains(result, "M07N6P7I4 — Test") {
		t.Error("missing mission row")
	}
}

func TestFormatCostBreakdown(t *testing.T) {
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

	result := FormatCostBreakdown(entries)

	if !strings.Contains(result, "deepseek-v4-pro") {
		t.Error("missing deepseek-v4-pro in breakdown")
	}
	if !strings.Contains(result, "deepseek-v4-flash") {
		t.Error("missing deepseek-v4-flash in breakdown")
	}
}

func TestFormatTraceTableVerbose(t *testing.T) {
	entries := []TraceEntry{
		{
			ID:           "E01",
			Seq:          1,
			TS:           "2026-07-12T10:30:00.000Z",
			Type:         EventCheckpoint,
			StepLabel:    strPtr("done"),
			LatencyMs:    0,
			InputTokens:  0,
			OutputTokens: 0,
		},
	}

	opts := TraceDisplayOptions{Verbose: true}
	result := FormatTraceTable(entries, "M07N6P7I4", "Test", opts)

	if !strings.Contains(result, "=== Trace:") {
		t.Error("missing header in verbose mode")
	}
}

func TestFormatTraceTableFlat(t *testing.T) {
	entries := []TraceEntry{
		{
			ID:           "E01",
			Seq:          1,
			TS:           "2026-07-12T10:30:00.000Z",
			Type:         EventSubagentSpawn,
			Subagent:     strPtr("sc-coder"),
			LatencyMs:    0,
			InputTokens:  0,
			OutputTokens: 0,
		},
		{
			ID:           "E02",
			Seq:          2,
			TS:           "2026-07-12T10:30:45.000Z",
			Type:         EventSubagentResult,
			Subagent:     strPtr("sc-coder"),
			TraceID:      strPtr("E01"),
			LatencyMs:    41000,
			InputTokens:  12000,
			OutputTokens: 2400,
		},
	}

	opts := TraceDisplayOptions{Flat: true}
	result := FormatTraceTable(entries, "M07N6P7I4", "Test", opts)

	if !strings.Contains(result, "subagent_spawn") {
		t.Error("missing subagent_spawn in flat output")
	}
	if !strings.Contains(result, "subagent_result") {
		t.Error("missing subagent_result in flat output")
	}
}

func intPtr(i int) *int { return &i }
