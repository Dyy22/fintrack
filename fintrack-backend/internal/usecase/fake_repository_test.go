package usecase

import (
	"context"
	"time"

	"fintrack-backend/internal/domain"

	"github.com/google/uuid"
)

type fakeRepository struct {
	createUserFn                func(ctx context.Context, email, passwordHash string) (domain.User, error)
	findUserByEmailFn           func(ctx context.Context, email string) (domain.User, error)
	listAccountTypesFn          func(ctx context.Context) ([]domain.AccountType, error)
	listAccountsFn              func(ctx context.Context, userID uuid.UUID) ([]domain.Account, error)
	findAccountFn               func(ctx context.Context, userID, accountID uuid.UUID) (domain.Account, error)
	accountTypeNameFn           func(ctx context.Context, accountTypeID int) (string, error)
	createAccountFn             func(ctx context.Context, userID uuid.UUID, name string, accountTypeID int, balance float64, goldGrams *float64, goldPrice *float64, stockSymbol *string, stockLots *float64, stockPrice *float64) (domain.Account, error)
	updateAccountFn             func(ctx context.Context, userID, accountID uuid.UUID, name *string, isActive *bool) (domain.Account, error)
	softDeleteAccountFn         func(ctx context.Context, userID, accountID uuid.UUID) error
	hardDeleteAccountFn         func(ctx context.Context, userID, accountID uuid.UUID) error
	listCategoriesFn            func(ctx context.Context, userID uuid.UUID, typ string) ([]domain.Category, error)
	createCategoryFn            func(ctx context.Context, userID uuid.UUID, name, typ string) (domain.Category, error)
	updateCategoryFn            func(ctx context.Context, userID, categoryID uuid.UUID, name string) (domain.Category, error)
	deleteCategoryFn            func(ctx context.Context, userID, categoryID uuid.UUID) error
	createTransactionFn         func(ctx context.Context, tx domain.Transaction) (domain.Transaction, error)
	listTransactionsFn          func(ctx context.Context, userID uuid.UUID, start, end, accountID, categoryID, typ string, limit, offset int) ([]domain.Transaction, error)
	netWorthFn                  func(ctx context.Context, userID uuid.UUID) (float64, []domain.Account, error)
	spendingByCategoryFn        func(ctx context.Context, userID uuid.UUID, start, end time.Time) (float64, []domain.SpendingCategory, float64, error)
	latestGoldPriceFn           func(ctx context.Context) (domain.GoldPrice, error)
	saveGoldPriceFn             func(ctx context.Context, price domain.GoldPrice) (domain.GoldPrice, error)
	listGoldPriceHistoryFn      func(ctx context.Context, days int) ([]domain.GoldPriceHistoryPoint, error)
	refreshGoldFn               func(ctx context.Context, price domain.GoldPrice) error
	refreshStockFn              func(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, quote domain.StockQuote) error
	latestStockQuoteFn          func(ctx context.Context, symbol string) (domain.StockQuote, error)
	saveStockQuoteFn            func(ctx context.Context, quote domain.StockQuote) (domain.StockQuote, error)
	latestMarketChartFn         func(ctx context.Context, symbol, rng, interval string) (domain.MarketChart, error)
	saveMarketChartFn           func(ctx context.Context, symbol, rng, interval string, chart domain.MarketChart) (domain.MarketChart, error)
	createBudgetFn              func(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, month, year int, amount float64) (domain.Budget, error)
	listBudgetsFn               func(ctx context.Context, userID uuid.UUID, month, year int) ([]domain.Budget, error)
	updateBudgetFn              func(ctx context.Context, userID, budgetID uuid.UUID, amount float64) (domain.Budget, error)
	deleteBudgetFn              func(ctx context.Context, userID, budgetID uuid.UUID) error
	spendingByCategoryInRangeFn func(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]domain.SpendingCategory, error)
}

func (f fakeRepository) CreateUser(ctx context.Context, email, passwordHash string) (domain.User, error) {
	return f.createUserFn(ctx, email, passwordHash)
}
func (f fakeRepository) FindUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return f.findUserByEmailFn(ctx, email)
}
func (f fakeRepository) ListAccountTypes(ctx context.Context) ([]domain.AccountType, error) {
	return f.listAccountTypesFn(ctx)
}
func (f fakeRepository) ListAccounts(ctx context.Context, userID uuid.UUID) ([]domain.Account, error) {
	return f.listAccountsFn(ctx, userID)
}
func (f fakeRepository) FindAccount(ctx context.Context, userID, accountID uuid.UUID) (domain.Account, error) {
	if f.findAccountFn != nil {
		return f.findAccountFn(ctx, userID, accountID)
	}
	return domain.Account{ID: accountID, UserID: userID, Type: "bank", IsActive: true}, nil
}
func (f fakeRepository) AccountTypeName(ctx context.Context, accountTypeID int) (string, error) {
	if f.accountTypeNameFn != nil {
		return f.accountTypeNameFn(ctx, accountTypeID)
	}
	return "bank", nil
}
func (f fakeRepository) CreateAccount(ctx context.Context, userID uuid.UUID, name string, accountTypeID int, balance float64, goldGrams *float64, goldPrice *float64, stockSymbol *string, stockLots *float64, stockPrice *float64) (domain.Account, error) {
	return f.createAccountFn(ctx, userID, name, accountTypeID, balance, goldGrams, goldPrice, stockSymbol, stockLots, stockPrice)
}
func (f fakeRepository) UpdateAccount(ctx context.Context, userID, accountID uuid.UUID, name *string, isActive *bool) (domain.Account, error) {
	return f.updateAccountFn(ctx, userID, accountID, name, isActive)
}
func (f fakeRepository) SoftDeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	return f.softDeleteAccountFn(ctx, userID, accountID)
}
func (f fakeRepository) HardDeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	return f.hardDeleteAccountFn(ctx, userID, accountID)
}
func (f fakeRepository) ListCategories(ctx context.Context, userID uuid.UUID, typ string) ([]domain.Category, error) {
	return f.listCategoriesFn(ctx, userID, typ)
}
func (f fakeRepository) CreateCategory(ctx context.Context, userID uuid.UUID, name, typ string) (domain.Category, error) {
	return f.createCategoryFn(ctx, userID, name, typ)
}
func (f fakeRepository) UpdateCategory(ctx context.Context, userID, categoryID uuid.UUID, name string) (domain.Category, error) {
	return f.updateCategoryFn(ctx, userID, categoryID, name)
}
func (f fakeRepository) DeleteCategory(ctx context.Context, userID, categoryID uuid.UUID) error {
	return f.deleteCategoryFn(ctx, userID, categoryID)
}
func (f fakeRepository) CreateTransaction(ctx context.Context, tx domain.Transaction) (domain.Transaction, error) {
	return f.createTransactionFn(ctx, tx)
}
func (f fakeRepository) ListTransactions(ctx context.Context, userID uuid.UUID, start, end, accountID, categoryID, typ string, limit, offset int) ([]domain.Transaction, error) {
	return f.listTransactionsFn(ctx, userID, start, end, accountID, categoryID, typ, limit, offset)
}
func (f fakeRepository) NetWorth(ctx context.Context, userID uuid.UUID) (float64, []domain.Account, error) {
	return f.netWorthFn(ctx, userID)
}
func (f fakeRepository) SpendingByCategory(ctx context.Context, userID uuid.UUID, start, end time.Time) (float64, []domain.SpendingCategory, float64, error) {
	return f.spendingByCategoryFn(ctx, userID, start, end)
}
func (f fakeRepository) LatestGoldPrice(ctx context.Context) (domain.GoldPrice, error) {
	if f.latestGoldPriceFn != nil {
		return f.latestGoldPriceFn(ctx)
	}
	return domain.GoldPrice{}, nil
}
func (f fakeRepository) SaveGoldPrice(ctx context.Context, price domain.GoldPrice) (domain.GoldPrice, error) {
	if f.saveGoldPriceFn != nil {
		return f.saveGoldPriceFn(ctx, price)
	}
	return price, nil
}
func (f fakeRepository) ListGoldPriceHistory(ctx context.Context, days int) ([]domain.GoldPriceHistoryPoint, error) {
	if f.listGoldPriceHistoryFn != nil {
		return f.listGoldPriceHistoryFn(ctx, days)
	}
	return nil, nil
}
func (f fakeRepository) RefreshGoldAccountBalances(ctx context.Context, price domain.GoldPrice) error {
	if f.refreshGoldFn != nil {
		return f.refreshGoldFn(ctx, price)
	}
	return nil
}
func (f fakeRepository) RefreshStockAccountBalance(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, quote domain.StockQuote) error {
	if f.refreshStockFn != nil {
		return f.refreshStockFn(ctx, userID, accountID, quote)
	}
	return nil
}
func (f fakeRepository) LatestStockQuote(ctx context.Context, symbol string) (domain.StockQuote, error) {
	if f.latestStockQuoteFn != nil {
		return f.latestStockQuoteFn(ctx, symbol)
	}
	return domain.StockQuote{}, nil
}
func (f fakeRepository) SaveStockQuote(ctx context.Context, quote domain.StockQuote) (domain.StockQuote, error) {
	if f.saveStockQuoteFn != nil {
		return f.saveStockQuoteFn(ctx, quote)
	}
	return quote, nil
}
func (f fakeRepository) LatestMarketChart(ctx context.Context, symbol, rng, interval string) (domain.MarketChart, error) {
	if f.latestMarketChartFn != nil {
		return f.latestMarketChartFn(ctx, symbol, rng, interval)
	}
	return domain.MarketChart{}, nil
}
func (f fakeRepository) SaveMarketChart(ctx context.Context, symbol, rng, interval string, chart domain.MarketChart) (domain.MarketChart, error) {
	if f.saveMarketChartFn != nil {
		return f.saveMarketChartFn(ctx, symbol, rng, interval, chart)
	}
	return chart, nil
}
func (f fakeRepository) CreateBudget(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, month, year int, amount float64) (domain.Budget, error) {
	if f.createBudgetFn != nil {
		return f.createBudgetFn(ctx, userID, categoryID, month, year, amount)
	}
	return domain.Budget{}, nil
}
func (f fakeRepository) ListBudgets(ctx context.Context, userID uuid.UUID, month, year int) ([]domain.Budget, error) {
	if f.listBudgetsFn != nil {
		return f.listBudgetsFn(ctx, userID, month, year)
	}
	return nil, nil
}
func (f fakeRepository) UpdateBudget(ctx context.Context, userID, budgetID uuid.UUID, amount float64) (domain.Budget, error) {
	if f.updateBudgetFn != nil {
		return f.updateBudgetFn(ctx, userID, budgetID, amount)
	}
	return domain.Budget{}, nil
}
func (f fakeRepository) DeleteBudget(ctx context.Context, userID, budgetID uuid.UUID) error {
	if f.deleteBudgetFn != nil {
		return f.deleteBudgetFn(ctx, userID, budgetID)
	}
	return nil
}
func (f fakeRepository) SpendingByCategoryInRange(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]domain.SpendingCategory, error) {
	if f.spendingByCategoryInRangeFn != nil {
		return f.spendingByCategoryInRangeFn(ctx, userID, start, end)
	}
	return nil, nil
}
