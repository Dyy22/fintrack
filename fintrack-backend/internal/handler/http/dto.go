package http

import "time"

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type createAccountRequest struct {
	Name          string   `json:"name" binding:"required"`
	AccountTypeID int      `json:"account_type_id" binding:"required"`
	Balance       float64  `json:"balance"`
	GoldGrams     *float64 `json:"gold_grams"`
	StockSymbol   *string  `json:"stock_symbol"`
	StockLots     *float64 `json:"stock_lots"`
}

type updateAccountRequest struct {
	Name     *string `json:"name"`
	IsActive *bool   `json:"is_active"`
}

type createCategoryRequest struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required,oneof=income expense"`
}

type updateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type createTransactionRequest struct {
	AccountID         string     `json:"account_id" binding:"required"`
	CategoryID        *string    `json:"category_id"`
	Type              string     `json:"type" binding:"required,oneof=income expense transfer"`
	Amount            float64    `json:"amount" binding:"required,gt=0"`
	GoldGrams         *float64   `json:"gold_grams"`
	Description       string     `json:"description"`
	Date              *time.Time `json:"date"`
	TransferAccountID *string    `json:"transfer_account_id"`
}

type createBudgetRequest struct {
	CategoryID string  `json:"category_id" binding:"required"`
	Month      int     `json:"month" binding:"required,min=1,max=12"`
	Year       int     `json:"year" binding:"required,min=2020"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
}

type updateBudgetRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}
