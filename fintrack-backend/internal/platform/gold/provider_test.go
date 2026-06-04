package gold

import "testing"

func TestParseLogamMuliaAPIPicksRegularOneGramGold(t *testing.T) {
	price, err := parseLogamMuliaAPI(sourceResponse{
		Timestamp: "2026-06-04T01:04:48.161Z",
		Data: []priceItem{
			{Source: "logammulia", Material: "gold", MaterialType: "Emas Batangan Gift Series", Weight: 1, WeightUnit: "gr", SellPrice: 2909000, Currency: "IDR", DisplayName: "Logam Mulia"},
			{Source: "logammulia", Material: "gold", MaterialType: "Emas Batangan", Weight: 1, WeightUnit: "gr", SellPrice: 2759000, Currency: "IDR", DisplayName: "Logam Mulia"},
		},
	})
	if err != nil {
		t.Fatalf("parseLogamMuliaAPI returned error: %v", err)
	}
	if price.PricePerGram != 2759000 {
		t.Fatalf("expected regular 1 gram price 2759000, got %f", price.PricePerGram)
	}
	if price.Source != "Logam Mulia" {
		t.Fatalf("expected source Logam Mulia, got %q", price.Source)
	}
}
