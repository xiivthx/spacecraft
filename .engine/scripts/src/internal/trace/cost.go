package trace

func CalcCost(entries []TraceEntry) float64 {
	var total float64
	for _, e := range entries {
		model := ""
		if e.Model != nil {
			model = *e.Model
		}
		price := ModelPrice(model)
		total += (float64(e.InputTokens) / 1_000_000.0) * price.Input
		total += (float64(e.OutputTokens) / 1_000_000.0) * price.Output
	}
	return total
}
