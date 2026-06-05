package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AccountType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Account struct {
	ID                 uuid.UUID `json:"id"`
	UserID             uuid.UUID `json:"-"`
	AccountTypeID      int       `json:"account_type_id,omitempty"`
	Type               string    `json:"type"`
	Name               string    `json:"name"`
	Balance            float64   `json:"balance"`
	Currency           string    `json:"currency"`
	GoldGrams          *float64  `json:"gold_grams,omitempty"`
	GoldPrice          *float64  `json:"gold_price_per_gram,omitempty"`
	StockSymbol        *string   `json:"stock_symbol,omitempty"`
	StockLots          *float64  `json:"stock_lots,omitempty"`
	StockPricePerShare *float64  `json:"stock_price_per_share,omitempty"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
	UpdatedAt          time.Time `json:"updated_at,omitempty"`
}

type Category struct {
	ID        uuid.UUID  `json:"id"`
	UserID    *uuid.UUID `json:"-"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	IsDefault bool       `json:"is_default"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
}

type Transaction struct {
	ID                uuid.UUID  `json:"id"`
	UserID            uuid.UUID  `json:"-"`
	AccountID         uuid.UUID  `json:"account_id,omitempty"`
	CategoryID        *uuid.UUID `json:"category_id,omitempty"`
	Type              string     `json:"type"`
	Amount            float64    `json:"amount"`
	GoldGrams         *float64   `json:"gold_grams,omitempty"`
	Description       string     `json:"description,omitempty"`
	Date              time.Time  `json:"date"`
	TransferAccountID *uuid.UUID `json:"transfer_account_id,omitempty"`
	CreatedAt         time.Time  `json:"created_at,omitempty"`
	UpdatedAt         time.Time  `json:"updated_at,omitempty"`
	Account           *Account   `json:"account,omitempty"`
	Category          *Category  `json:"category,omitempty"`
}

type SpendingCategory struct {
	CategoryID uuid.UUID `json:"category_id"`
	Name       string    `json:"name"`
	Amount     float64   `json:"amount"`
	Percentage float64   `json:"percentage"`
}

type GoldPrice struct {
	PricePerGram float64   `json:"price_per_gram"`
	Source       string    `json:"source"`
	FetchedAt    time.Time `json:"fetched_at"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

type GoldPriceHistoryPoint struct {
	Date         string  `json:"date"`
	PricePerGram float64 `json:"price_per_gram"`
	Source       string  `json:"source"`
}

type StockQuote struct {
	Symbol      string    `json:"symbol"`
	Name        string    `json:"name,omitempty"`
	Price       float64   `json:"price"`
	Currency    string    `json:"currency"`
	Source      string    `json:"source"`
	FetchedAt   time.Time `json:"fetched_at"`
	CacheStatus string    `json:"cache_status,omitempty"`
}

type MarketChartPoint struct {
	Time  time.Time `json:"time"`
	Close float64   `json:"close"`
}

type MarketChart struct {
	Symbol      string             `json:"symbol"`
	Name        string             `json:"name,omitempty"`
	Currency    string             `json:"currency"`
	Source      string             `json:"source"`
	FetchedAt   time.Time          `json:"fetched_at"`
	CacheStatus string             `json:"cache_status,omitempty"`
	Points      []MarketChartPoint `json:"points"`
}

type Budget struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"-"`
	CategoryID uuid.UUID `json:"category_id"`
	Category   *Category `json:"category,omitempty"`
	Month      int       `json:"month"`
	Year       int       `json:"year"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type BudgetWithSpending struct {
	Budget
	Spent     float64 `json:"spent"`
	Remaining float64 `json:"remaining"`
	Percent   float64 `json:"percent"`
}
