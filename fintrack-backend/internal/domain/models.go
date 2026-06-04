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
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"-"`
	AccountTypeID int       `json:"account_type_id,omitempty"`
	Type          string    `json:"type"`
	Name          string    `json:"name"`
	Balance       float64   `json:"balance"`
	Currency      string    `json:"currency"`
	GoldGrams     *float64  `json:"gold_grams,omitempty"`
	GoldPrice     *float64  `json:"gold_price_per_gram,omitempty"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
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
	Name       string  `json:"name"`
	Amount     float64 `json:"amount"`
	Percentage float64 `json:"percentage"`
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
