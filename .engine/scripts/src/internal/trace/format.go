package trace

import (
	"fmt"
	"strings"
	"time"
)

type TraceDisplayOptions struct {
	Verbose bool
	Flat    bool
}

type CostRow struct {
	Mission   string
	TokensIn  int
	TokensOut int
	Cost      float64
}

func FormatTraceTable(entries []TraceEntry, missionID, missionTitle string, opts TraceDisplayOptions) string {
	var b strings.Builder

	// Header
	b.WriteString(fmt.Sprintf("=== Trace: %s — %q ===\n\n", missionID, missionTitle))

	if len(entries) == 0 {
		b.WriteString("No trace events recorded.\n")
		return b.String()
	}

	// Determine base time for relative timestamps
	baseTS, err := time.Parse(time.RFC3339Nano, entries[0].TS)
	if err != nil {
		baseTS = time.Time{}
	}

	// Indent tracking for subagent nesting
	indent := 0

	for _, e := range entries {
		relTS := formatRelativeTS(e.TS, baseTS)
		latStr := formatLatency(e.LatencyMs)
		typeAndDetail := formatEventTypeDetail(e, opts)

		if !opts.Flat && e.TraceID != nil {
			indent = 2
		} else if !opts.Flat && e.Type == EventSubagentResult {
			indent = 2
		}

		pad := ""
		if indent > 0 {
			pad = strings.Repeat("  ", indent)
		}

		b.WriteString(fmt.Sprintf("%s#%d  %s  %-24s %s  %s",
			pad, e.Seq, relTS, typeAndDetail, latStr, formatTokens(e)))

		if e.ExitCode != nil {
			b.WriteString(fmt.Sprintf("  exit=%d", *e.ExitCode))
		}

		b.WriteString("\n")

		// Reset indent for non-spawn events
		if e.Type != EventSubagentSpawn && indent > 0 {
			indent = 0
		}
	}

	// Summary
	s := ComputeSummary(entries)
	b.WriteString(fmt.Sprintf("\nSummary: %d events | %.1fs total | %s in | %s out | $%.3f\n",
		s.EventCount,
		float64(s.TotalLatencyMs)/1000.0,
		formatNumber(s.TotalInTokens),
		formatNumber(s.TotalOutTokens),
		s.EstimatedCost))

	return b.String()
}

func FormatCostTable(rows []CostRow, total CostRow) string {
	var b strings.Builder

	// Header
	b.WriteString(fmt.Sprintf("%-50s %12s %12s %12s\n", "Mission", "Tokens In", "Tokens Out", "Est. Cost"))

	// Rows
	for _, r := range rows {
		b.WriteString(fmt.Sprintf("%-50s %12s %12s $%.3f\n",
			truncateString(r.Mission, 48),
			formatNumber(r.TokensIn),
			formatNumber(r.TokensOut),
			r.Cost))
	}

	// Separator
	b.WriteString(strings.Repeat("─", 88) + "\n")

	// Total
	b.WriteString(fmt.Sprintf("%-50s %12s %12s $%.3f\n",
		total.Mission,
		formatNumber(total.TokensIn),
		formatNumber(total.TokensOut),
		total.Cost))

	return b.String()
}

func FormatCostBreakdown(entries []TraceEntry) string {
	type modelAgg struct {
		inTokens  int
		outTokens int
	}

	agg := make(map[string]*modelAgg)
	seen := make(map[string]bool)
	var order []string

	for _, e := range entries {
		if e.Model == nil {
			continue
		}
		model := *e.Model
		if !seen[model] {
			seen[model] = true
			order = append(order, model)
		}
		if a, ok := agg[model]; ok {
			a.inTokens += e.InputTokens
			a.outTokens += e.OutputTokens
		} else {
			agg[model] = &modelAgg{inTokens: e.InputTokens, outTokens: e.OutputTokens}
		}
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%-30s %12s %12s %12s\n", "Model", "Tokens In", "Tokens Out", "Est. Cost"))

	for _, model := range order {
		a := agg[model]
		price := ModelPrice(model)
		cost := (float64(a.inTokens)/1_000_000.0)*price.Input + (float64(a.outTokens)/1_000_000.0)*price.Output
		b.WriteString(fmt.Sprintf("%-30s %12s %12s $%.3f\n",
			model,
			formatNumber(a.inTokens),
			formatNumber(a.outTokens),
			cost))
	}

	return b.String()
}

func formatRelativeTS(ts string, base time.Time) string {
	t, err := time.Parse(time.RFC3339Nano, ts)
	if err != nil {
		return "??:??:??.???"
	}
	d := t.Sub(base)
	totalMs := int(d.Milliseconds())
	if totalMs < 0 {
		totalMs = 0
	}

	hours := totalMs / 3600000
	totalMs %= 3600000
	minutes := totalMs / 60000
	totalMs %= 60000
	seconds := totalMs / 1000
	millis := totalMs % 1000

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, millis)
}

func formatLatency(ms int) string {
	if ms == 0 {
		return "     --"
	}
	if ms < 1000 {
		return fmt.Sprintf("%6dms", ms)
	}
	return fmt.Sprintf("%5.1fs ", float64(ms)/1000.0)
}

func formatEventTypeDetail(e TraceEntry, opts TraceDisplayOptions) string {
	switch e.Type {
	case EventToolCall:
		tool := "?"
		if e.Tool != nil {
			tool = *e.Tool
		}
		args := formatArgs(e, opts.Verbose)
		return fmt.Sprintf("%-12s %-8s %s", e.Type, tool, args)
	case EventModelInvoke:
		model := "?"
		if e.Model != nil {
			model = *e.Model
		}
		return fmt.Sprintf("%-12s %-8s", e.Type, model)
	case EventSubagentSpawn:
		sub := "?"
		if e.Subagent != nil {
			sub = *e.Subagent
		}
		return fmt.Sprintf("%-14s %-6s", e.Type, sub)
	case EventSubagentResult:
		sub := "?"
		if e.Subagent != nil {
			sub = *e.Subagent
		}
		return fmt.Sprintf("%-14s %-6s", e.Type, sub)
	case EventCheckpoint:
		label := ""
		if e.StepLabel != nil {
			label = *e.StepLabel
		}
		return fmt.Sprintf("%-12s %s", e.Type, label)
	default:
		return string(e.Type)
	}
}

func formatArgs(e TraceEntry, verbose bool) string {
	if e.Args == nil || len(e.Args) == 0 || string(e.Args) == "null" {
		return ""
	}
	s := string(e.Args)
	s = strings.TrimSpace(s)
	if !verbose && len(s) > 40 {
		s = s[:37] + "..."
	}
	return s
}

func formatTokens(e TraceEntry) string {
	in := formatNumber(e.InputTokens)
	out := formatNumber(e.OutputTokens)
	if e.InputTokens == 0 && e.OutputTokens == 0 {
		return ""
	}
	return fmt.Sprintf("in:%s out:%s", in, out)
}

func formatNumber(n int) string {
	if n == 0 {
		return "0"
	}
	s := fmt.Sprintf("%d", n)
	var result []byte
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}

func truncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
