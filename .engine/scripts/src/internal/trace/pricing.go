package trace

type Price struct {
	Input  float64
	Output float64
}

var pricingTable = map[string]Price{
	"deepseek-v4-pro":  {Input: 2.00, Output: 8.00},
	"deepseek-v4-flash": {Input: 0.50, Output: 2.00},
	"kimi-k2.7-code":   {Input: 1.50, Output: 6.00},
	"qwen3.7-plus":     {Input: 1.00, Output: 4.00},
	"glm-5.2":          {Input: 0.80, Output: 3.20},
	"unknown":          {Input: 1.00, Output: 4.00},
}

func ModelPrice(model string) Price {
	if p, ok := pricingTable[model]; ok {
		return p
	}
	return pricingTable["unknown"]
}
