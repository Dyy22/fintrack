package gold

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"fintrack-backend/internal/domain"
)

var ErrNoPriceSource = errors.New("gold price source is not configured")

type Provider struct {
	client        *http.Client
	sourceURL     string
	fallbackPrice float64
}

type sourceResponse struct {
	PricePerGram float64     `json:"price_per_gram"`
	Price        float64     `json:"price"`
	SellPrice    float64     `json:"sell_price"`
	BuyPrice     float64     `json:"buy_price"`
	Harga        float64     `json:"harga"`
	Source       string      `json:"source"`
	FetchedAt    string      `json:"fetched_at"`
	Timestamp    string      `json:"timestamp"`
	Data         []priceItem `json:"data"`
}

type priceItem struct {
	Source       string  `json:"source"`
	Material     string  `json:"material"`
	MaterialType string  `json:"materialType"`
	Weight       float64 `json:"weight"`
	WeightUnit   string  `json:"weightUnit"`
	SellPrice    float64 `json:"sellPrice"`
	BuybackPrice float64 `json:"buybackPrice"`
	Currency     string  `json:"currency"`
	RecordedDate string  `json:"recordedDate"`
	DisplayName  string  `json:"displayName"`
}

func NewProvider(sourceURL string, fallbackPrice float64) *Provider {
	return &Provider{
		client:        &http.Client{Timeout: 10 * time.Second},
		sourceURL:     strings.TrimSpace(sourceURL),
		fallbackPrice: fallbackPrice,
	}
}

func (p *Provider) Latest(ctx context.Context) (domain.GoldPrice, error) {
	if p.sourceURL == "" {
		if p.fallbackPrice > 0 {
			return domain.GoldPrice{PricePerGram: p.fallbackPrice, Source: "env_fallback", FetchedAt: time.Now().UTC()}, nil
		}
		return domain.GoldPrice{}, ErrNoPriceSource
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.sourceURL, nil)
	if err != nil {
		return p.fallback(ctx, err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "fintrack/1.0")

	res, err := p.client.Do(req)
	if err != nil {
		return p.fallback(ctx, err)
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return p.fallback(ctx, errors.New(res.Status))
	}

	var payload sourceResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return p.fallback(ctx, err)
	}

	if len(payload.Data) > 0 {
		price, err := parseLogamMuliaAPI(payload)
		if err != nil {
			return p.fallback(ctx, err)
		}
		return price, nil
	}

	price := firstPositive(payload.PricePerGram, payload.Price, payload.SellPrice, payload.BuyPrice, payload.Harga)
	if price <= 0 {
		return p.fallback(ctx, errors.New("gold price response does not contain a positive price"))
	}

	fetchedAt := parseTimestamp(payload.FetchedAt, time.Now().UTC())
	source := payload.Source
	if source == "" {
		source = p.sourceURL
	}
	return domain.GoldPrice{PricePerGram: price, Source: source, FetchedAt: fetchedAt}, nil
}

func (p *Provider) fallback(_ context.Context, err error) (domain.GoldPrice, error) {
	if p.fallbackPrice > 0 {
		return domain.GoldPrice{PricePerGram: p.fallbackPrice, Source: "env_fallback", FetchedAt: time.Now().UTC()}, nil
	}
	return domain.GoldPrice{}, err
}

func parseLogamMuliaAPI(payload sourceResponse) (domain.GoldPrice, error) {
	var selected *priceItem
	for i := range payload.Data {
		item := payload.Data[i]
		if !isOneGramGold(item) {
			continue
		}
		if selected == nil {
			selected = &payload.Data[i]
		}
		if strings.EqualFold(strings.TrimSpace(item.MaterialType), "Emas Batangan") {
			selected = &payload.Data[i]
			break
		}
	}
	if selected == nil || selected.SellPrice <= 0 {
		return domain.GoldPrice{}, errors.New("logam mulia API response does not contain a 1 gram gold sell price")
	}

	source := selected.DisplayName
	if source == "" {
		source = selected.Source
	}
	if source == "" {
		source = "logam-mulia-api"
	}

	return domain.GoldPrice{
		PricePerGram: selected.SellPrice,
		Source:       source,
		FetchedAt:    parseTimestamp(payload.Timestamp, time.Now().UTC()),
	}, nil
}

func isOneGramGold(item priceItem) bool {
	return strings.EqualFold(item.Material, "gold") &&
		strings.EqualFold(item.WeightUnit, "gr") &&
		item.Weight == 1 &&
		strings.EqualFold(item.Currency, "IDR")
}

func parseTimestamp(value string, fallback time.Time) time.Time {
	if value == "" {
		return fallback
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return fallback
	}
	return parsed.UTC()
}

func firstPositive(values ...float64) float64 {
	for _, value := range values {
		if value > 0 {
			return value
		}
	}
	return 0
}
