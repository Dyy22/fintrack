package http

import (
	"context"
	stdhttp "net/http"
	"strconv"
	"time"

	"fintrack-backend/internal/domain"
	"fintrack-backend/internal/middleware"
	"fintrack-backend/internal/platform/response"
	"fintrack-backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Usecase interface {
	Register(ctx context.Context, email, password string) (domain.User, error)
	Login(ctx context.Context, email, password string) (usecase.LoginResult, error)
	ListAccountTypes(ctx context.Context) ([]domain.AccountType, error)
	ListAccounts(ctx context.Context, userID uuid.UUID) ([]domain.Account, error)
	CreateAccount(ctx context.Context, userID uuid.UUID, name string, accountTypeID int, balance float64, goldGrams *float64, stockSymbol *string, stockLots *float64) (domain.Account, error)
	UpdateAccount(ctx context.Context, userID, accountID uuid.UUID, name *string, isActive *bool) (domain.Account, error)
	SoftDeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error
	HardDeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error
	ListCategories(ctx context.Context, userID uuid.UUID, typ string) ([]domain.Category, error)
	CreateCategory(ctx context.Context, userID uuid.UUID, name, typ string) (domain.Category, error)
	UpdateCategory(ctx context.Context, userID, categoryID uuid.UUID, name string) (domain.Category, error)
	DeleteCategory(ctx context.Context, userID, categoryID uuid.UUID) error
	ListTransactions(ctx context.Context, userID uuid.UUID, start, end, accountID, categoryID, typ string, limit, offset int) ([]domain.Transaction, error)
	CreateTransaction(ctx context.Context, tx domain.Transaction) (domain.Transaction, error)
	NetWorth(ctx context.Context, userID uuid.UUID) (float64, []domain.Account, error)
	SpendingByCategory(ctx context.Context, userID uuid.UUID, startDate, endDate string) (time.Time, time.Time, float64, []domain.SpendingCategory, float64, error)
	LatestGoldPrice(ctx context.Context) (domain.GoldPrice, error)
	GoldPriceHistory(ctx context.Context, days int) ([]domain.GoldPriceHistoryPoint, error)
	MarketChart(ctx context.Context, symbol, rng, interval string) (domain.MarketChart, error)
	CreateBudget(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, month, year int, amount float64) (domain.BudgetWithSpending, error)
	ListBudgets(ctx context.Context, userID uuid.UUID, month, year int) ([]domain.BudgetWithSpending, error)
	UpdateBudget(ctx context.Context, userID, budgetID uuid.UUID, amount float64) (domain.BudgetWithSpending, error)
	DeleteBudget(ctx context.Context, userID, budgetID uuid.UUID) error
}

type Handler struct{ uc Usecase }

func New(uc Usecase) *Handler { return &Handler{uc: uc} }

func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	user, err := h.uc.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(stdhttp.StatusCreated, user)
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	result, err := h.uc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, result)
}

func (h *Handler) ListAccountTypes(c *gin.Context) {
	accountTypes, err := h.uc.ListAccountTypes(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, gin.H{"account_types": accountTypes})
}

func (h *Handler) ListAccounts(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized)
		return
	}
	accounts, err := h.uc.ListAccounts(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, gin.H{"accounts": accounts})
}

func (h *Handler) CreateAccount(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized)
		return
	}
	var req createAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	account, err := h.uc.CreateAccount(c.Request.Context(), userID, req.Name, req.AccountTypeID, req.Balance, req.GoldGrams, req.StockSymbol, req.StockLots)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(stdhttp.StatusCreated, account)
}

func (h *Handler) UpdateAccount(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized)
		return
	}
	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, response.ErrBadRequest)
		return
	}
	var req updateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	account, err := h.uc.UpdateAccount(c.Request.Context(), userID, accountID, req.Name, req.IsActive)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, account)
}

func (h *Handler) DeleteAccount(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized)
		return
	}
	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, response.ErrBadRequest)
		return
	}
	if c.Query("hard") == "true" {
		if err := h.uc.HardDeleteAccount(c.Request.Context(), userID, accountID); err != nil {
			response.Error(c, err)
			return
		}
		c.Status(stdhttp.StatusNoContent)
		return
	}
	if err := h.uc.SoftDeleteAccount(c.Request.Context(), userID, accountID); err != nil {
		response.Error(c, err)
		return
	}
	c.Status(stdhttp.StatusNoContent)
}

func (h *Handler) ListCategories(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized)
		return
	}
	categories, err := h.uc.ListCategories(c.Request.Context(), userID, c.Query("type"))
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, gin.H{"categories": categories})
}

func (h *Handler) CreateCategory(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized)
		return
	}
	var req createCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	cat, err := h.uc.CreateCategory(c.Request.Context(), userID, req.Name, req.Type)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(stdhttp.StatusCreated, cat)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized)
		return
	}
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, response.ErrBadRequest)
		return
	}
	var req updateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	cat, err := h.uc.UpdateCategory(c.Request.Context(), userID, categoryID, req.Name)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, cat)
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized)
		return
	}
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, response.ErrBadRequest)
		return
	}
	if err := h.uc.DeleteCategory(c.Request.Context(), userID, categoryID); err != nil {
		response.Error(c, err)
		return
	}
	c.Status(stdhttp.StatusNoContent)
}

func (h *Handler) ListTransactions(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized)
		return
	}
	limit, offset, err := paginationQuery(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	transactions, err := h.uc.ListTransactions(c.Request.Context(), userID, c.Query("start_date"), c.Query("end_date"), c.Query("account_id"), c.Query("category_id"), c.Query("type"), limit, offset)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, gin.H{"transactions": transactions, "limit": limit, "offset": offset})
}

func (h *Handler) CreateTransaction(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized)
		return
	}
	var req createTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		response.Error(c, response.ErrBadRequest)
		return
	}
	var categoryID *uuid.UUID
	if req.CategoryID != nil {
		id, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			response.Error(c, response.ErrBadRequest)
			return
		}
		categoryID = &id
	}
	var transferAccountID *uuid.UUID
	if req.TransferAccountID != nil {
		id, err := uuid.Parse(*req.TransferAccountID)
		if err != nil {
			response.Error(c, response.ErrBadRequest)
			return
		}
		transferAccountID = &id
	}
	date := req.Date
	txDate := zeroTime(date)
	tx, err := h.uc.CreateTransaction(c.Request.Context(), domain.Transaction{UserID: userID, AccountID: accountID, CategoryID: categoryID, Type: req.Type, Amount: req.Amount, GoldGrams: req.GoldGrams, Description: req.Description, Date: txDate, TransferAccountID: transferAccountID})
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(stdhttp.StatusCreated, tx)
}

func (h *Handler) NetWorth(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized)
		return
	}
	total, accounts, err := h.uc.NetWorth(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, gin.H{"net_worth": total, "accounts": accounts})
}

func (h *Handler) GoldPrice(c *gin.Context) {
	price, err := h.uc.LatestGoldPrice(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, price)
}

func (h *Handler) GoldPriceHistory(c *gin.Context) {
	days, err := strconv.Atoi(c.DefaultQuery("days", "7"))
	if err != nil {
		response.Error(c, response.ErrBadRequest)
		return
	}
	history, err := h.uc.GoldPriceHistory(c.Request.Context(), days)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, gin.H{"history": history})
}

func (h *Handler) MarketChart(c *gin.Context) {
	chart, err := h.uc.MarketChart(c.Request.Context(), c.DefaultQuery("symbol", "IHSG"), c.DefaultQuery("range", "1mo"), c.DefaultQuery("interval", "1d"))
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, chart)
}

func (h *Handler) ListBudgets(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized)
		return
	}
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		response.Error(c, response.ErrBadRequest)
		return
	}
	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		response.Error(c, response.ErrBadRequest)
		return
	}
	budgets, err := h.uc.ListBudgets(c.Request.Context(), userID, month, year)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, gin.H{"budgets": budgets})
}

func (h *Handler) CreateBudget(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized)
		return
	}
	var req createBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		response.Error(c, response.ErrBadRequest)
		return
	}
	budget, err := h.uc.CreateBudget(c.Request.Context(), userID, categoryID, req.Month, req.Year, req.Amount)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(stdhttp.StatusCreated, budget)
}

func (h *Handler) UpdateBudget(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized)
		return
	}
	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, response.ErrBadRequest)
		return
	}
	var req updateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}
	budget, err := h.uc.UpdateBudget(c.Request.Context(), userID, budgetID, req.Amount)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, budget)
}

func (h *Handler) DeleteBudget(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized)
		return
	}
	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, response.ErrBadRequest)
		return
	}
	if err := h.uc.DeleteBudget(c.Request.Context(), userID, budgetID); err != nil {
		response.Error(c, err)
		return
	}
	c.Status(stdhttp.StatusNoContent)
}

func (h *Handler) SpendingByCategory(c *gin.Context) {
	userID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, response.ErrUnauthorized)
		return
	}
	start, end, total, categories, totalIncome, err := h.uc.SpendingByCategory(c.Request.Context(), userID, c.Query("start_date"), c.Query("end_date"))
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(stdhttp.StatusOK, gin.H{"start_date": start.Format("2006-01-02"), "end_date": end.Format("2006-01-02"), "total_spending": total, "total_income": totalIncome, "categories": categories})
}

func zeroTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func paginationQuery(c *gin.Context) (int, int, error) {
	limit, err := intQuery(c, "limit", 50)
	if err != nil {
		return 0, 0, err
	}
	offset, err := intQuery(c, "offset", 0)
	if err != nil {
		return 0, 0, err
	}
	if limit < 0 || offset < 0 {
		return 0, 0, response.ErrBadRequest
	}
	if limit == 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	return limit, offset, nil
}

func intQuery(c *gin.Context, key string, fallback int) (int, error) {
	value := c.Query(key)
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, response.ErrBadRequest
	}
	return parsed, nil
}
