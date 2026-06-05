package stock

import (
	"testing"
	"time"

	"fintrack-backend/internal/domain"
)

func TestQuoteCache(t *testing.T) {
	provider := NewProvider()
	quote := domain.StockQuote{Symbol: "TLKM", Price: 2900, Currency: "IDR"}
	provider.storeQuote("TLKM", quote)

	cached, ok := provider.cachedQuote("TLKM", false)
	if !ok || cached.Price != quote.Price {
		t.Fatalf("expected fresh cached quote, got quote=%+v ok=%v", cached, ok)
	}

	provider.mu.Lock()
	entry := provider.quotes["TLKM"]
	entry.expiresAt = time.Now().Add(-time.Minute)
	provider.quotes["TLKM"] = entry
	provider.mu.Unlock()

	if _, ok := provider.cachedQuote("TLKM", false); ok {
		t.Fatal("expected expired cache to be ignored when allowExpired=false")
	}
	if cached, ok := provider.cachedQuote("TLKM", true); !ok || cached.Price != quote.Price {
		t.Fatalf("expected stale cached quote when allowExpired=true, got quote=%+v ok=%v", cached, ok)
	}
}

func TestChartCacheKey(t *testing.T) {
	if got, want := chartCacheKey("tlkm.jk", "3mo", "1d"), "TLKM|3mo|1d"; got != want {
		t.Fatalf("chartCacheKey() = %q, want %q", got, want)
	}
	if got, want := chartCacheKey("^JKSE", "3mo", "1d"), "IHSG|3mo|1d"; got != want {
		t.Fatalf("chartCacheKey() = %q, want %q", got, want)
	}
}

func TestIsIDXExchange(t *testing.T) {
	tests := []struct {
		name             string
		exchangeName     string
		fullExchangeName string
		want             bool
	}{
		{name: "yahoo jkt code", exchangeName: "JKT", fullExchangeName: "Jakarta", want: true},
		{name: "jakarta full name", exchangeName: "", fullExchangeName: "Jakarta", want: true},
		{name: "missing metadata allowed", exchangeName: "", fullExchangeName: "", want: true},
		{name: "non idx", exchangeName: "NASDAQ", fullExchangeName: "NasdaqGS", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isIDXExchange(tt.exchangeName, tt.fullExchangeName); got != tt.want {
				t.Fatalf("isIDXExchange() = %v, want %v", got, tt.want)
			}
		})
	}
}
