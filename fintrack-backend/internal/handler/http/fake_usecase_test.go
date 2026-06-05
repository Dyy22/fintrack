package http

import (
	"context"
	"time"

	"fintrack-backend/internal/domain"
	"fintrack-backend/internal/usecase"

	"github.com/google/uuid"
)

type fakeUsecase struct {
	registerFn           func(ctx context.Context, email, password string) (domain.User, error)
	loginFn              func(ctx context.Context, email, password string) (usecase.LoginResult, error)
	listAccountTypesFn   func(ctx context.Context) ([]domain.AccountType, error)
	listAccountsFn       func(ctx context.Context, userID uuid.UUID) ([]domain.Account, error)
	createAccountFn      func(ctx context.Context, userID uuid.UUID, name string, accountTypeID int, balance float64, goldGrams *float64, stockSymbol *string, stockLots *float64) (domain.Account, error)
	updateAccountFn      func(ctx context.Context, userID, accountID uuid.UUID, name *string, isActive *bool) (domain.Account, error)
	softDeleteAccountFn  func(ctx context.Context, userID, accountID uuid.UUID) error
	hardDeleteAccountFn  func(ctx context.Context, userID, accountID uuid.UUID) error
	listCategoriesFn     func(ctx context.Context, userID uuid.UUID, typ string) ([]domain.Category, error)
	createCategoryFn     func(ctx context.Context, userID uuid.UUID, name, typ string) (domain.Category, error)
	updateCategoryFn     func(ctx context.Context, userID, categoryID uuid.UUID, name string) (domain.Category, error)
	deleteCategoryFn     func(ctx context.Context, userID, categoryID uuid.UUID) error
	listTransactionsFn   func(ctx context.Context, userID uuid.UUID, start, end, accountID, categoryID, typ string, limit, offset int) ([]domain.Transaction, error)
	createTransactionFn  func(ctx context.Context, tx domain.Transaction) (domain.Transaction, error)
	netWorthFn           func(ctx context.Context, userID uuid.UUID) (float64, []domain.Account, error)
	spendingByCategoryFn func(ctx context.Context, userID uuid.UUID, startDate, endDate string) (time.Time, time.Time, float64, []domain.SpendingCategory, float64, error)
	latestGoldPriceFn    func(ctx context.Context) (domain.GoldPrice, error)
	goldPriceHistoryFn   func(ctx context.Context, days int) ([]domain.GoldPriceHistoryPoint, error)
	marketChartFn        func(ctx context.Context, symbol, rng, interval string) (domain.MarketChart, error)
	createBudgetFn       func(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, month, year int, amount float64) (domain.BudgetWithSpending, error)
	listBudgetsFn        func(ctx context.Context, userID uuid.UUID, month, year int) ([]domain.BudgetWithSpending, error)
	updateBudgetFn       func(ctx context.Context, userID, budgetID uuid.UUID, amount float64) (domain.BudgetWithSpending, error)
	deleteBudgetFn       func(ctx context.Context, userID, budgetID uuid.UUID) error
}

func (f fakeUsecase) Register(ctx context.Context, email, password string) (domain.User, error) {
	return f.registerFn(ctx, email, password)
}
func (f fakeUsecase) Login(ctx context.Context, email, password string) (usecase.LoginResult, error) {
	return f.loginFn(ctx, email, password)
}
func (f fakeUsecase) ListAccountTypes(ctx context.Context) ([]domain.AccountType, error) {
	return f.listAccountTypesFn(ctx)
}
func (f fakeUsecase) ListAccounts(ctx context.Context, userID uuid.UUID) ([]domain.Account, error) {
	return f.listAccountsFn(ctx, userID)
}
func (f fakeUsecase) CreateAccount(ctx context.Context, userID uuid.UUID, name string, accountTypeID int, balance float64, goldGrams *float64, stockSymbol *string, stockLots *float64) (domain.Account, error) {
	return f.createAccountFn(ctx, userID, name, accountTypeID, balance, goldGrams, stockSymbol, stockLots)
}
func (f fakeUsecase) UpdateAccount(ctx context.Context, userID, accountID uuid.UUID, name *string, isActive *bool) (domain.Account, error) {
	return f.updateAccountFn(ctx, userID, accountID, name, isActive)
}
func (f fakeUsecase) SoftDeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	return f.softDeleteAccountFn(ctx, userID, accountID)
}
func (f fakeUsecase) HardDeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	return f.hardDeleteAccountFn(ctx, userID, accountID)
}
func (f fakeUsecase) ListCategories(ctx context.Context, userID uuid.UUID, typ string) ([]domain.Category, error) {
	return f.listCategoriesFn(ctx, userID, typ)
}
func (f fakeUsecase) CreateCategory(ctx context.Context, userID uuid.UUID, name, typ string) (domain.Category, error) {
	return f.createCategoryFn(ctx, userID, name, typ)
}
func (f fakeUsecase) UpdateCategory(ctx context.Context, userID, categoryID uuid.UUID, name string) (domain.Category, error) {
	return f.updateCategoryFn(ctx, userID, categoryID, name)
}
func (f fakeUsecase) DeleteCategory(ctx context.Context, userID, categoryID uuid.UUID) error {
	return f.deleteCategoryFn(ctx, userID, categoryID)
}
func (f fakeUsecase) ListTransactions(ctx context.Context, userID uuid.UUID, start, end, accountID, categoryID, typ string, limit, offset int) ([]domain.Transaction, error) {
	return f.listTransactionsFn(ctx, userID, start, end, accountID, categoryID, typ, limit, offset)
}
func (f fakeUsecase) CreateTransaction(ctx context.Context, tx domain.Transaction) (domain.Transaction, error) {
	return f.createTransactionFn(ctx, tx)
}
func (f fakeUsecase) NetWorth(ctx context.Context, userID uuid.UUID) (float64, []domain.Account, error) {
	return f.netWorthFn(ctx, userID)
}
func (f fakeUsecase) SpendingByCategory(ctx context.Context, userID uuid.UUID, startDate, endDate string) (time.Time, time.Time, float64, []domain.SpendingCategory, float64, error) {
	return f.spendingByCategoryFn(ctx, userID, startDate, endDate)
}
func (f fakeUsecase) LatestGoldPrice(ctx context.Context) (domain.GoldPrice, error) {
	if f.latestGoldPriceFn != nil {
		return f.latestGoldPriceFn(ctx)
	}
	return domain.GoldPrice{}, nil
}
func (f fakeUsecase) GoldPriceHistory(ctx context.Context, days int) ([]domain.GoldPriceHistoryPoint, error) {
	if f.goldPriceHistoryFn != nil {
		return f.goldPriceHistoryFn(ctx, days)
	}
	return nil, nil
}
func (f fakeUsecase) MarketChart(ctx context.Context, symbol, rng, interval string) (domain.MarketChart, error) {
	if f.marketChartFn != nil {
		return f.marketChartFn(ctx, symbol, rng, interval)
	}
	return domain.MarketChart{}, nil
}
func (f fakeUsecase) CreateBudget(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, month, year int, amount float64) (domain.BudgetWithSpending, error) {
	if f.createBudgetFn != nil {
		return f.createBudgetFn(ctx, userID, categoryID, month, year, amount)
	}
	return domain.BudgetWithSpending{}, nil
}
func (f fakeUsecase) ListBudgets(ctx context.Context, userID uuid.UUID, month, year int) ([]domain.BudgetWithSpending, error) {
	if f.listBudgetsFn != nil {
		return f.listBudgetsFn(ctx, userID, month, year)
	}
	return nil, nil
}
func (f fakeUsecase) UpdateBudget(ctx context.Context, userID, budgetID uuid.UUID, amount float64) (domain.BudgetWithSpending, error) {
	if f.updateBudgetFn != nil {
		return f.updateBudgetFn(ctx, userID, budgetID, amount)
	}
	return domain.BudgetWithSpending{}, nil
}
func (f fakeUsecase) DeleteBudget(ctx context.Context, userID, budgetID uuid.UUID) error {
	if f.deleteBudgetFn != nil {
		return f.deleteBudgetFn(ctx, userID, budgetID)
	}
	return nil
}
