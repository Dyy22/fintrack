package usecase

import (
	"context"
	"math"
	"strings"
	"time"

	"fintrack-backend/internal/domain"
	"fintrack-backend/internal/platform/response"
	"fintrack-backend/internal/platform/security"

	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, email, passwordHash string) (domain.User, error)
	FindUserByEmail(ctx context.Context, email string) (domain.User, error)
	ListAccountTypes(ctx context.Context) ([]domain.AccountType, error)
	ListAccounts(ctx context.Context, userID uuid.UUID) ([]domain.Account, error)
	FindAccount(ctx context.Context, userID, accountID uuid.UUID) (domain.Account, error)
	AccountTypeName(ctx context.Context, accountTypeID int) (string, error)
	CreateAccount(ctx context.Context, userID uuid.UUID, name string, accountTypeID int, balance float64, goldGrams *float64, goldPrice *float64, stockSymbol *string, stockLots *float64, stockPrice *float64) (domain.Account, error)
	UpdateAccount(ctx context.Context, userID, accountID uuid.UUID, name *string, isActive *bool) (domain.Account, error)
	SoftDeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error
	HardDeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error
	ListCategories(ctx context.Context, userID uuid.UUID, typ string) ([]domain.Category, error)
	CreateCategory(ctx context.Context, userID uuid.UUID, name, typ string) (domain.Category, error)
	UpdateCategory(ctx context.Context, userID, categoryID uuid.UUID, name string) (domain.Category, error)
	DeleteCategory(ctx context.Context, userID, categoryID uuid.UUID) error
	CreateTransaction(ctx context.Context, tx domain.Transaction) (domain.Transaction, error)
	ListTransactions(ctx context.Context, userID uuid.UUID, start, end, accountID, categoryID, typ string, limit, offset int) ([]domain.Transaction, error)
	NetWorth(ctx context.Context, userID uuid.UUID) (float64, []domain.Account, error)
	SpendingByCategory(ctx context.Context, userID uuid.UUID, start, end time.Time) (float64, []domain.SpendingCategory, float64, error)
	LatestGoldPrice(ctx context.Context) (domain.GoldPrice, error)
	SaveGoldPrice(ctx context.Context, price domain.GoldPrice) (domain.GoldPrice, error)
	ListGoldPriceHistory(ctx context.Context, days int) ([]domain.GoldPriceHistoryPoint, error)
	RefreshGoldAccountBalances(ctx context.Context, price domain.GoldPrice) error
	RefreshStockAccountBalance(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, quote domain.StockQuote) error
	LatestStockQuote(ctx context.Context, symbol string) (domain.StockQuote, error)
	SaveStockQuote(ctx context.Context, quote domain.StockQuote) (domain.StockQuote, error)
	LatestMarketChart(ctx context.Context, symbol, rng, interval string) (domain.MarketChart, error)
	SaveMarketChart(ctx context.Context, symbol, rng, interval string, chart domain.MarketChart) (domain.MarketChart, error)
	CreateBudget(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, month, year int, amount float64) (domain.Budget, error)
	ListBudgets(ctx context.Context, userID uuid.UUID, month, year int) ([]domain.Budget, error)
	UpdateBudget(ctx context.Context, userID, budgetID uuid.UUID, amount float64) (domain.Budget, error)
	DeleteBudget(ctx context.Context, userID, budgetID uuid.UUID) error
	SpendingByCategoryInRange(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]domain.SpendingCategory, error)
}

type GoldPriceProvider interface {
	Latest(ctx context.Context) (domain.GoldPrice, error)
}

type StockMarketProvider interface {
	Quote(ctx context.Context, symbol string) (domain.StockQuote, error)
	Chart(ctx context.Context, symbol, rng, interval string) (domain.MarketChart, error)
}

const (
	stockQuoteRefreshInterval  = 5 * time.Minute
	marketChartRefreshInterval = 15 * time.Minute
)

type Usecases struct {
	repo                     Repository
	jwt                      security.JWTService
	goldProvider             GoldPriceProvider
	stockProvider            StockMarketProvider
	goldPriceRefreshInterval time.Duration
}

func New(repo Repository, jwt security.JWTService) *Usecases {
	return &Usecases{repo: repo, jwt: jwt, goldPriceRefreshInterval: time.Hour}
}

func (u *Usecases) WithGoldPriceProvider(provider GoldPriceProvider, refreshInterval time.Duration) *Usecases {
	u.goldProvider = provider
	if refreshInterval > 0 {
		u.goldPriceRefreshInterval = refreshInterval
	}
	return u
}

func (u *Usecases) WithStockMarketProvider(provider StockMarketProvider) *Usecases {
	u.stockProvider = provider
	return u
}

type LoginResult struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}

func (u *Usecases) Register(ctx context.Context, email, password string) (domain.User, error) {
	email = normalizeEmail(email)
	if len(password) < 8 {
		return domain.User{}, response.ErrBadRequest
	}
	hash, err := security.HashPassword(password)
	if err != nil {
		return domain.User{}, err
	}
	user, err := u.repo.CreateUser(ctx, email, hash)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return domain.User{}, response.ErrConflict
		}
		return domain.User{}, err
	}
	return user, nil
}

func (u *Usecases) Login(ctx context.Context, email, password string) (LoginResult, error) {
	user, err := u.repo.FindUserByEmail(ctx, normalizeEmail(email))
	if err != nil {
		return LoginResult{}, response.ErrUnauthorized
	}
	if !security.CheckPassword(user.PasswordHash, password) {
		return LoginResult{}, response.ErrUnauthorized
	}
	token, err := u.jwt.Generate(user.ID, user.Email)
	if err != nil {
		return LoginResult{}, err
	}
	return LoginResult{Token: token, User: user}, nil
}

func (u *Usecases) ListAccountTypes(ctx context.Context) ([]domain.AccountType, error) {
	return u.repo.ListAccountTypes(ctx)
}

func (u *Usecases) ListAccounts(ctx context.Context, userID uuid.UUID) ([]domain.Account, error) {
	if price, err := u.LatestGoldPrice(ctx); err == nil {
		_ = u.repo.RefreshGoldAccountBalances(ctx, price)
	}
	_ = u.refreshUserStockAccounts(ctx, userID)
	return u.repo.ListAccounts(ctx, userID)
}
func (u *Usecases) CreateAccount(ctx context.Context, userID uuid.UUID, name string, accountTypeID int, balance float64, goldGrams *float64, stockSymbol *string, stockLots *float64) (domain.Account, error) {
	name = strings.TrimSpace(name)
	if name == "" || accountTypeID <= 0 {
		return domain.Account{}, response.ErrBadRequest
	}
	accountType, err := u.repo.AccountTypeName(ctx, accountTypeID)
	if err != nil {
		return domain.Account{}, err
	}
	var goldPrice *float64
	var normalizedStockSymbol *string
	var stockPrice *float64
	if accountType == "gold" {
		if goldGrams == nil || *goldGrams < 0 || balance < 0 {
			return domain.Account{}, response.ErrBadRequest
		}
		price, err := u.LatestGoldPrice(ctx)
		if err != nil {
			return domain.Account{}, err
		}
		if !goldValueMatches(balance, *goldGrams, price.PricePerGram) {
			return domain.Account{}, response.ErrBadRequest
		}
		goldPrice = &price.PricePerGram
	} else {
		goldGrams = nil
	}
	if accountType == "stock_broker" {
		if stockSymbol == nil || strings.TrimSpace(*stockSymbol) == "" || stockLots == nil || *stockLots <= 0 {
			return domain.Account{}, response.ErrBadRequest
		}
		symbol := normalizeIDXSymbol(*stockSymbol)
		if symbol == "" {
			return domain.Account{}, response.ErrBadRequest
		}
		quote, quoteErr := u.StockQuote(ctx, symbol)
		if quoteErr != nil || quote.Price <= 0 {
			if !isKnownIDXSymbol(symbol) {
				if quoteErr != nil {
					return domain.Account{}, quoteErr
				}
				return domain.Account{}, response.ErrBadRequest
			}
			balance = 0
		} else {
			if quote.Currency != "" && quote.Currency != "IDR" {
				return domain.Account{}, response.ErrBadRequest
			}
			balance = *stockLots * 100 * quote.Price
			stockPrice = &quote.Price
			symbol = quote.Symbol
		}
		normalizedStockSymbol = &symbol
		name = symbol
	} else {
		stockLots = nil
	}
	return u.repo.CreateAccount(ctx, userID, name, accountTypeID, balance, goldGrams, goldPrice, normalizedStockSymbol, stockLots, stockPrice)
}
func (u *Usecases) UpdateAccount(ctx context.Context, userID, accountID uuid.UUID, name *string, isActive *bool) (domain.Account, error) {
	if name != nil {
		trimmed := strings.TrimSpace(*name)
		if trimmed == "" {
			return domain.Account{}, response.ErrBadRequest
		}
		name = &trimmed
	}
	return u.repo.UpdateAccount(ctx, userID, accountID, name, isActive)
}
func (u *Usecases) SoftDeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	return u.repo.SoftDeleteAccount(ctx, userID, accountID)
}
func (u *Usecases) HardDeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	return u.repo.HardDeleteAccount(ctx, userID, accountID)
}

func (u *Usecases) ListCategories(ctx context.Context, userID uuid.UUID, typ string) ([]domain.Category, error) {
	if typ != "" && typ != "income" && typ != "expense" {
		return nil, response.ErrBadRequest
	}
	return u.repo.ListCategories(ctx, userID, typ)
}
func (u *Usecases) CreateCategory(ctx context.Context, userID uuid.UUID, name, typ string) (domain.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" || (typ != "income" && typ != "expense") {
		return domain.Category{}, response.ErrBadRequest
	}
	cat, err := u.repo.CreateCategory(ctx, userID, name, typ)
	if err != nil && (strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique")) {
		return domain.Category{}, response.ErrConflict
	}
	return cat, err
}
func (u *Usecases) UpdateCategory(ctx context.Context, userID, categoryID uuid.UUID, name string) (domain.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.Category{}, response.ErrBadRequest
	}
	cat, err := u.repo.UpdateCategory(ctx, userID, categoryID, name)
	if err != nil && (strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique")) {
		return domain.Category{}, response.ErrConflict
	}
	return cat, err
}
func (u *Usecases) DeleteCategory(ctx context.Context, userID, categoryID uuid.UUID) error {
	return u.repo.DeleteCategory(ctx, userID, categoryID)
}

func (u *Usecases) CreateTransaction(ctx context.Context, tx domain.Transaction) (domain.Transaction, error) {
	if tx.Amount <= 0 {
		return domain.Transaction{}, response.ErrBadRequest
	}
	if tx.Date.IsZero() {
		tx.Date = time.Now().UTC()
	}
	if tx.Type == "transfer" {
		if tx.TransferAccountID == nil || tx.CategoryID != nil || *tx.TransferAccountID == tx.AccountID {
			return domain.Transaction{}, response.ErrBadRequest
		}
	} else if tx.Type == "income" || tx.Type == "expense" {
		if tx.TransferAccountID != nil {
			return domain.Transaction{}, response.ErrBadRequest
		}
		if tx.CategoryID == nil {
			return domain.Transaction{}, response.ErrBadRequest
		}
	} else {
		return domain.Transaction{}, response.ErrBadRequest
	}
	if err := u.validateGoldTransactionAmount(ctx, tx); err != nil {
		return domain.Transaction{}, err
	}
	return u.repo.CreateTransaction(ctx, tx)
}

func (u *Usecases) ListTransactions(ctx context.Context, userID uuid.UUID, start, end, accountID, categoryID, typ string, limit, offset int) ([]domain.Transaction, error) {
	if typ != "" && typ != "income" && typ != "expense" && typ != "transfer" {
		return nil, response.ErrBadRequest
	}
	limit, offset, err := normalizePagination(limit, offset)
	if err != nil {
		return nil, err
	}
	return u.repo.ListTransactions(ctx, userID, start, end, accountID, categoryID, typ, limit, offset)
}

func (u *Usecases) NetWorth(ctx context.Context, userID uuid.UUID) (float64, []domain.Account, error) {
	if price, err := u.LatestGoldPrice(ctx); err == nil {
		_ = u.repo.RefreshGoldAccountBalances(ctx, price)
	}
	_ = u.refreshUserStockAccounts(ctx, userID)
	return u.repo.NetWorth(ctx, userID)
}

func (u *Usecases) LatestGoldPrice(ctx context.Context) (domain.GoldPrice, error) {
	cached, err := u.repo.LatestGoldPrice(ctx)
	if err == nil && time.Since(cached.FetchedAt) < u.goldPriceRefreshInterval {
		return cached, nil
	}
	if u.goldProvider == nil {
		return cached, err
	}
	fresh, fetchErr := u.goldProvider.Latest(ctx)
	if fetchErr != nil {
		if err == nil {
			return cached, nil
		}
		return domain.GoldPrice{}, fetchErr
	}
	fresh, err = u.repo.SaveGoldPrice(ctx, fresh)
	if err != nil {
		return domain.GoldPrice{}, err
	}
	_ = u.repo.RefreshGoldAccountBalances(ctx, fresh)
	return fresh, nil
}
func (u *Usecases) GoldPriceHistory(ctx context.Context, days int) ([]domain.GoldPriceHistoryPoint, error) {
	if days <= 0 {
		days = 7
	}
	if days > 30 {
		days = 30
	}
	_, _ = u.LatestGoldPrice(ctx)
	return u.repo.ListGoldPriceHistory(ctx, days)
}

func (u *Usecases) MarketChart(ctx context.Context, symbol, rng, interval string) (domain.MarketChart, error) {
	if rng != "" && rng != "5d" && rng != "1mo" && rng != "3mo" && rng != "6mo" && rng != "1y" {
		return domain.MarketChart{}, response.ErrBadRequest
	}
	if interval != "" && interval != "1d" && interval != "1wk" && interval != "1mo" {
		return domain.MarketChart{}, response.ErrBadRequest
	}
	if rng == "" {
		rng = "1mo"
	}
	if interval == "" {
		interval = "1d"
	}
	symbol = normalizeMarketSymbol(symbol)
	cached, cacheErr := u.repo.LatestMarketChart(ctx, symbol, rng, interval)
	if cacheErr == nil && time.Since(cached.FetchedAt) < marketChartRefreshInterval {
		return cached, nil
	}
	if u.stockProvider == nil {
		return cached, cacheErr
	}
	fresh, fetchErr := u.stockProvider.Chart(ctx, symbol, rng, interval)
	if fetchErr != nil {
		if cacheErr == nil {
			return cached, nil
		}
		return domain.MarketChart{}, fetchErr
	}
	fresh.Symbol = normalizeMarketSymbol(fresh.Symbol)
	if fresh.FetchedAt.IsZero() {
		fresh.FetchedAt = time.Now().UTC()
	}
	stored, err := u.repo.SaveMarketChart(ctx, symbol, rng, interval, fresh)
	if err != nil {
		return domain.MarketChart{}, err
	}
	return stored, nil
}

func (u *Usecases) StockQuote(ctx context.Context, symbol string) (domain.StockQuote, error) {
	symbol = normalizeIDXSymbol(symbol)
	cached, cacheErr := u.repo.LatestStockQuote(ctx, symbol)
	if cacheErr == nil && time.Since(cached.FetchedAt) < stockQuoteRefreshInterval {
		return cached, nil
	}
	if u.stockProvider == nil {
		return cached, cacheErr
	}
	fresh, fetchErr := u.stockProvider.Quote(ctx, symbol)
	if fetchErr != nil {
		if cacheErr == nil {
			return cached, nil
		}
		return domain.StockQuote{}, fetchErr
	}
	fresh.Symbol = normalizeIDXSymbol(fresh.Symbol)
	stored, err := u.repo.SaveStockQuote(ctx, fresh)
	if err != nil {
		return domain.StockQuote{}, err
	}
	return stored, nil
}

func (u *Usecases) SpendingByCategory(ctx context.Context, userID uuid.UUID, startDate, endDate string) (time.Time, time.Time, float64, []domain.SpendingCategory, float64, error) {
	start, end, err := reportRange(startDate, endDate)
	if err != nil {
		return time.Time{}, time.Time{}, 0, nil, 0, err
	}
	total, items, totalIncome, err := u.repo.SpendingByCategory(ctx, userID, start, end)
	return start, end.AddDate(0, 0, -1), total, items, totalIncome, err
}

func (u *Usecases) refreshUserStockAccounts(ctx context.Context, userID uuid.UUID) error {
	if u.stockProvider == nil {
		return nil
	}
	accounts, err := u.repo.ListAccounts(ctx, userID)
	if err != nil {
		return err
	}
	for _, account := range accounts {
		if account.Type != "stock_broker" || account.StockSymbol == nil || account.StockLots == nil || !account.IsActive {
			continue
		}
		quote, err := u.StockQuote(ctx, *account.StockSymbol)
		if err != nil {
			continue
		}
		_ = u.repo.RefreshStockAccountBalance(ctx, userID, account.ID, quote)
	}
	return nil
}

func (u *Usecases) validateGoldTransactionAmount(ctx context.Context, tx domain.Transaction) error {
	primaryAccount, err := u.repo.FindAccount(ctx, tx.UserID, tx.AccountID)
	if err != nil {
		return err
	}
	goldInvolved := primaryAccount.Type == "gold"
	if tx.TransferAccountID != nil {
		transferAccount, err := u.repo.FindAccount(ctx, tx.UserID, *tx.TransferAccountID)
		if err != nil {
			return err
		}
		goldInvolved = goldInvolved || transferAccount.Type == "gold"
	}
	if !goldInvolved {
		tx.GoldGrams = nil
		return nil
	}
	if tx.GoldGrams == nil || *tx.GoldGrams <= 0 {
		return response.ErrBadRequest
	}
	price, err := u.LatestGoldPrice(ctx)
	if err != nil {
		return err
	}
	if !goldValueMatches(tx.Amount, *tx.GoldGrams, price.PricePerGram) {
		return response.ErrBadRequest
	}
	return nil
}

func goldValueMatches(amount, grams, pricePerGram float64) bool {
	if pricePerGram <= 0 || amount < 0 || grams < 0 {
		return false
	}
	return math.Abs(amount-(grams*pricePerGram)) <= 1
}

func normalizeEmail(email string) string { return strings.ToLower(strings.TrimSpace(email)) }

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

func isKnownIDXSymbol(symbol string) bool {
	_, ok := knownIDXSymbols[normalizeIDXSymbol(symbol)]
	return ok
}

var knownIDXSymbols = map[string]struct{}{
	"AALI": {},
	"ACES": {},
	"ADRO": {},
	"AKRA": {},
	"ANTM": {},
	"ARTO": {},
	"ASII": {},
	"BBCA": {},
	"BBNI": {},
	"BBRI": {},
	"BBTN": {},
	"BMRI": {},
	"BRIS": {},
	"BRPT": {},
	"BUKA": {},
	"CPIN": {},
	"EMTK": {},
	"ERAA": {},
	"EXCL": {},
	"GOTO": {},
	"ICBP": {},
	"INCO": {},
	"INDF": {},
	"INKP": {},
	"INTP": {},
	"ITMG": {},
	"KLBF": {},
	"MDKA": {},
	"MEDC": {},
	"PGAS": {},
	"PTBA": {},
	"SMGR": {},
	"TLKM": {},
	"TOWR": {},
	"UNTR": {},
	"UNVR": {},
}

func normalizePagination(limit, offset int) (int, int, error) {
	if limit == 0 {
		limit = 50
	}
	if limit < 0 || offset < 0 {
		return 0, 0, response.ErrBadRequest
	}
	if limit > 100 {
		limit = 100
	}
	return limit, offset, nil
}

func reportRange(startDate, endDate string) (time.Time, time.Time, error) {
	loc := time.UTC
	if startDate == "" && endDate == "" {
		now := time.Now().UTC()
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		return start, start.AddDate(0, 1, 0), nil
	}
	if startDate == "" || endDate == "" {
		return time.Time{}, time.Time{}, response.ErrBadRequest
	}
	start, err := time.ParseInLocation("2006-01-02", startDate, loc)
	if err != nil {
		return time.Time{}, time.Time{}, response.ErrBadRequest
	}
	endInclusive, err := time.ParseInLocation("2006-01-02", endDate, loc)
	if err != nil {
		return time.Time{}, time.Time{}, response.ErrBadRequest
	}
	if endInclusive.Before(start) {
		return time.Time{}, time.Time{}, response.ErrBadRequest
	}
	return start, endInclusive.AddDate(0, 0, 1), nil
}

func (u *Usecases) CreateBudget(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, month, year int, amount float64) (domain.BudgetWithSpending, error) {
	if month < 1 || month > 12 || year < 2000 || amount <= 0 {
		return domain.BudgetWithSpending{}, response.ErrBadRequest
	}
	budget, err := u.repo.CreateBudget(ctx, userID, categoryID, month, year, amount)
	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			return domain.BudgetWithSpending{}, response.ErrConflict
		}
		return domain.BudgetWithSpending{}, err
	}
	return u.enrichBudgetWithSpending(ctx, budget, month, year), nil
}

func (u *Usecases) ListBudgets(ctx context.Context, userID uuid.UUID, month, year int) ([]domain.BudgetWithSpending, error) {
	if month < 1 || month > 12 || year < 2000 {
		return nil, response.ErrBadRequest
	}
	budgets, err := u.repo.ListBudgets(ctx, userID, month, year)
	if err != nil {
		return nil, err
	}
	result := make([]domain.BudgetWithSpending, len(budgets))
	for i, b := range budgets {
		result[i] = u.enrichBudgetWithSpending(ctx, b, month, year)
	}
	return result, nil
}

func (u *Usecases) UpdateBudget(ctx context.Context, userID, budgetID uuid.UUID, amount float64) (domain.BudgetWithSpending, error) {
	if amount <= 0 {
		return domain.BudgetWithSpending{}, response.ErrBadRequest
	}
	budget, err := u.repo.UpdateBudget(ctx, userID, budgetID, amount)
	if err != nil {
		return domain.BudgetWithSpending{}, err
	}
	return u.enrichBudgetWithSpending(ctx, budget, budget.Month, budget.Year), nil
}

func (u *Usecases) DeleteBudget(ctx context.Context, userID, budgetID uuid.UUID) error {
	return u.repo.DeleteBudget(ctx, userID, budgetID)
}

func (u *Usecases) enrichBudgetWithSpending(ctx context.Context, budget domain.Budget, month, year int) domain.BudgetWithSpending {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	spending, err := u.repo.SpendingByCategoryInRange(ctx, budget.UserID, start, end)
	if err != nil {
		return domain.BudgetWithSpending{Budget: budget, Spent: 0, Remaining: budget.Amount, Percent: 0}
	}
	var spent float64
	for _, s := range spending {
		if s.CategoryID == budget.CategoryID {
			spent = s.Amount
			break
		}
	}
	remaining := budget.Amount - spent
	percent := 0.0
	if budget.Amount > 0 {
		percent = math.Round((spent/budget.Amount)*10000) / 100
	}
	return domain.BudgetWithSpending{Budget: budget, Spent: spent, Remaining: remaining, Percent: percent}
}
