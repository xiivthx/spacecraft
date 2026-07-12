package trace

type TraceSummary struct {
	EventCount     int
	TotalLatencyMs int
	TotalInTokens  int
	TotalOutTokens int
	EstimatedCost  float64
	FirstTS        string
	LastTS         string
}

func ComputeSummary(entries []TraceEntry) TraceSummary {
	var s TraceSummary
	s.EventCount = len(entries)
	if s.EventCount == 0 {
		return s
	}
	s.FirstTS = entries[0].TS
	s.LastTS = entries[len(entries)-1].TS
	for _, e := range entries {
		s.TotalLatencyMs += e.LatencyMs
		s.TotalInTokens += e.InputTokens
		s.TotalOutTokens += e.OutputTokens
	}
	s.EstimatedCost = CalcCost(entries)
	return s
}
