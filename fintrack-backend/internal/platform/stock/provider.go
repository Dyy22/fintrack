package stock

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"fintrack-backend/internal/domain"
	"fintrack-backend/internal/platform/response"
)

const (
	yahooChartURL = "https://query2.finance.yahoo.com/v8/finance/chart/%s"
	quoteCacheTTL = 5 * time.Minute
	chartCacheTTL = 15 * time.Minute
)

type Provider struct {
	client *http.Client
	mu     sync.RWMutex
	quotes map[string]quoteCacheEntry
	charts map[string]chartCacheEntry
}

type quoteCacheEntry struct {
	value     domain.StockQuote
	expiresAt time.Time
}

type chartCacheEntry struct {
	value     domain.MarketChart
	expiresAt time.Time
}

func NewProvider() *Provider {
	return &Provider{
		client: &http.Client{Timeout: 8 * time.Second},
		quotes: make(map[string]quoteCacheEntry),
		charts: make(map[string]chartCacheEntry),
	}
}

func (p *Provider) Quote(ctx context.Context, symbol string) (domain.StockQuote, error) {
	symbol = normalizeIDXSymbol(symbol)
	if cached, ok := p.cachedQuote(symbol, false); ok {
		return cached, nil
	}
	chart, err := p.fetchChart(ctx, yahooSymbol(symbol), "5d", "1d")
	if err != nil {
		if cached, ok := p.cachedQuote(symbol, true); ok {
			return cached, nil
		}
		return domain.StockQuote{}, err
	}
	if len(chart.Chart.Result) == 0 {
		return domain.StockQuote{}, response.ErrNotFound
	}
	result := chart.Chart.Result[0]
	if !isIDXExchange(result.Meta.ExchangeName, result.Meta.FullExchangeName) {
		return domain.StockQuote{}, response.ErrBadRequest
	}
	price := result.Meta.RegularMarketPrice
	if price <= 0 && len(result.Indicators.Quote) > 0 {
		closes := result.Indicators.Quote[0].Close
		for i := len(closes) - 1; i >= 0; i-- {
			if closes[i] != nil && *closes[i] > 0 {
				price = *closes[i]
				break
			}
		}
	}
	if price <= 0 {
		return domain.StockQuote{}, response.ErrNotFound
	}
	quote := domain.StockQuote{
		Symbol:    symbol,
		Name:      result.Meta.ShortName,
		Price:     price,
		Currency:  result.Meta.Currency,
		Source:    "Yahoo Finance",
		FetchedAt: time.Now().UTC(),
	}
	p.storeQuote(symbol, quote)
	return quote, nil
}

func (p *Provider) Chart(ctx context.Context, symbol, rng, interval string) (domain.MarketChart, error) {
	symbol = normalizeMarketSymbol(symbol)
	if rng == "" {
		rng = "1mo"
	}
	if interval == "" {
		interval = "1d"
	}
	cacheKey := chartCacheKey(symbol, rng, interval)
	if cached, ok := p.cachedChart(cacheKey, false); ok {
		return cached, nil
	}
	chart, err := p.fetchChart(ctx, yahooMarketSymbol(symbol), rng, interval)
	if err != nil {
		if cached, ok := p.cachedChart(cacheKey, true); ok {
			return cached, nil
		}
		return domain.MarketChart{}, err
	}
	if len(chart.Chart.Result) == 0 {
		return domain.MarketChart{}, response.ErrNotFound
	}
	result := chart.Chart.Result[0]
	if len(result.Timestamp) == 0 || len(result.Indicators.Quote) == 0 {
		return domain.MarketChart{}, response.ErrNotFound
	}
	closes := result.Indicators.Quote[0].Close
	points := make([]domain.MarketChartPoint, 0, len(result.Timestamp))
	for i, unix := range result.Timestamp {
		if i >= len(closes) || closes[i] == nil || *closes[i] <= 0 {
			continue
		}
		points = append(points, domain.MarketChartPoint{
			Time:  time.Unix(unix, 0).UTC(),
			Close: *closes[i],
		})
	}
	if len(points) == 0 {
		return domain.MarketChart{}, response.ErrNotFound
	}
	marketChart := domain.MarketChart{
		Symbol:    symbol,
		Name:      result.Meta.ShortName,
		Currency:  result.Meta.Currency,
		Source:    "Yahoo Finance",
		FetchedAt: time.Now().UTC(),
		Points:    points,
	}
	p.storeChart(cacheKey, marketChart)
	return marketChart, nil
}

func (p *Provider) cachedQuote(symbol string, allowExpired bool) (domain.StockQuote, bool) {
	p.mu.RLock()
	entry, ok := p.quotes[symbol]
	p.mu.RUnlock()
	if !ok {
		return domain.StockQuote{}, false
	}
	if !allowExpired && time.Now().After(entry.expiresAt) {
		return domain.StockQuote{}, false
	}
	return entry.value, true
}

func (p *Provider) storeQuote(symbol string, quote domain.StockQuote) {
	p.mu.Lock()
	p.quotes[symbol] = quoteCacheEntry{value: quote, expiresAt: time.Now().Add(quoteCacheTTL)}
	p.mu.Unlock()
}

func (p *Provider) cachedChart(key string, allowExpired bool) (domain.MarketChart, bool) {
	p.mu.RLock()
	entry, ok := p.charts[key]
	p.mu.RUnlock()
	if !ok {
		return domain.MarketChart{}, false
	}
	if !allowExpired && time.Now().After(entry.expiresAt) {
		return domain.MarketChart{}, false
	}
	return entry.value, true
}

func (p *Provider) storeChart(key string, chart domain.MarketChart) {
	p.mu.Lock()
	p.charts[key] = chartCacheEntry{value: chart, expiresAt: time.Now().Add(chartCacheTTL)}
	p.mu.Unlock()
}

func chartCacheKey(symbol, rng, interval string) string {
	return strings.Join([]string{normalizeMarketSymbol(symbol), rng, interval}, "|")
}

func (p *Provider) fetchChart(ctx context.Context, symbol, rng, interval string) (yahooChartResponse, error) {
	endpoint := fmt.Sprintf(yahooChartURL, url.PathEscape(symbol))
	query := url.Values{}
	query.Set("range", rng)
	query.Set("interval", interval)
	query.Set("includePrePost", "false")
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+query.Encode(), nil)
	if err != nil {
		return yahooChartResponse{}, err
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Fintrack/1.0)")
	res, err := p.client.Do(request)
	if err != nil {
		return yahooChartResponse{}, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return yahooChartResponse{}, response.ErrNotFound
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return yahooChartResponse{}, response.ErrBadRequest
	}
	var payload yahooChartResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return yahooChartResponse{}, err
	}
	if payload.Chart.Error != nil {
		return yahooChartResponse{}, response.ErrBadRequest
	}
	return payload, nil
}

func normalizeIDXSymbol(symbol string) string {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	symbol = strings.TrimSuffix(symbol, ".JK")
	return symbol
}

func normalizeMarketSymbol(symbol string) string {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" || symbol == "IHSG" || symbol == "JKSE" || symbol == "^JKSE" {
		return "IHSG"
	}
	return normalizeIDXSymbol(symbol)
}

func yahooSymbol(symbol string) string {
	return normalizeIDXSymbol(symbol) + ".JK"
}

func isIDXExchange(exchangeName, fullExchangeName string) bool {
	exchangeName = strings.ToUpper(strings.TrimSpace(exchangeName))
	fullExchangeName = strings.ToUpper(strings.TrimSpace(fullExchangeName))
	if exchangeName == "" && fullExchangeName == "" {
		return true
	}
	return exchangeName == "JKT" ||
		exchangeName == "JKSE" ||
		strings.Contains(exchangeName, "JAKARTA") ||
		strings.Contains(fullExchangeName, "JAKARTA")
}

func yahooMarketSymbol(symbol string) string {
	if normalizeMarketSymbol(symbol) == "IHSG" {
		return "^JKSE"
	}
	return yahooSymbol(symbol)
}

type yahooChartResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency           string  `json:"currency"`
				ShortName          string  `json:"shortName"`
				ExchangeName       string  `json:"exchangeName"`
				FullExchangeName   string  `json:"fullExchangeName"`
				RegularMarketPrice float64 `json:"regularMarketPrice"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Close []*float64 `json:"close"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error any `json:"error"`
	} `json:"chart"`
}
