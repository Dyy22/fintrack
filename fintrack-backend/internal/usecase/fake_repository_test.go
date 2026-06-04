package usecase

import (
	"context"
	"time"

	"fintrack-backend/internal/domain"

	"github.com/google/uuid"
)

type fakeRepository struct {
	createUserFn           func(ctx context.Context, email, passwordHash string) (domain.User, error)
	findUserByEmailFn      func(ctx context.Context, email string) (domain.User, error)
	listAccountTypesFn     func(ctx context.Context) ([]domain.AccountType, error)
	listAccountsFn         func(ctx context.Context, userID uuid.UUID) ([]domain.Account, error)
	findAccountFn          func(ctx context.Context, userID, accountID uuid.UUID) (domain.Account, error)
	accountTypeNameFn      func(ctx context.Context, accountTypeID int) (string, error)
	createAccountFn        func(ctx context.Context, userID uuid.UUID, name string, accountTypeID int, balance float64, goldGrams *float64, goldPrice *float64) (domain.Account, error)
	updateAccountFn        func(ctx context.Context, userID, accountID uuid.UUID, name *string, isActive *bool) (domain.Account, error)
	softDeleteAccountFn    func(ctx context.Context, userID, accountID uuid.UUID) error
	hardDeleteAccountFn    func(ctx context.Context, userID, accountID uuid.UUID) error
	listCategoriesFn       func(ctx context.Context, userID uuid.UUID, typ string) ([]domain.Category, error)
	createCategoryFn       func(ctx context.Context, userID uuid.UUID, name, typ string) (domain.Category, error)
	updateCategoryFn       func(ctx context.Context, userID, categoryID uuid.UUID, name string) (domain.Category, error)
	deleteCategoryFn       func(ctx context.Context, userID, categoryID uuid.UUID) error
	createTransactionFn    func(ctx context.Context, tx domain.Transaction) (domain.Transaction, error)
	listTransactionsFn     func(ctx context.Context, userID uuid.UUID, start, end, accountID, categoryID, typ string, limit, offset int) ([]domain.Transaction, error)
	netWorthFn             func(ctx context.Context, userID uuid.UUID) (float64, []domain.Account, error)
	spendingByCategoryFn   func(ctx context.Context, userID uuid.UUID, start, end time.Time) (float64, []domain.SpendingCategory, float64, error)
	latestGoldPriceFn      func(ctx context.Context) (domain.GoldPrice, error)
	saveGoldPriceFn        func(ctx context.Context, price domain.GoldPrice) (domain.GoldPrice, error)
	listGoldPriceHistoryFn func(ctx context.Context, days int) ([]domain.GoldPriceHistoryPoint, error)
	refreshGoldFn          func(ctx context.Context, price domain.GoldPrice) error
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
func (f fakeRepository) CreateAccount(ctx context.Context, userID uuid.UUID, name string, accountTypeID int, balance float64, goldGrams *float64, goldPrice *float64) (domain.Account, error) {
	return f.createAccountFn(ctx, userID, name, accountTypeID, balance, goldGrams, goldPrice)
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
