//go:build integration

package postgres

import (
	"context"
	"database/sql"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"fintrack-backend/internal/domain"
	"fintrack-backend/internal/platform/database"
	"fintrack-backend/internal/platform/response"

	"github.com/google/uuid"
)

const testDatabaseURLEnv = "FINTRACK_TEST_DB_URL"

func setupIntegrationTestDB(t *testing.T) (*sql.DB, *Repositories) {
	t.Helper()

	databaseURL := os.Getenv(testDatabaseURLEnv)
	if databaseURL == "" {
		t.Skipf("%s is not set", testDatabaseURLEnv)
	}
	assertTestDatabaseURL(t, databaseURL)

	db, err := database.Open(databaseURL)
	if err != nil {
		t.Fatalf("failed to connect test database: %v", err)
	}

	migrationDir := filepath.Join("..", "..", "..", "migrations")
	downSQL, err := os.ReadFile(filepath.Join(migrationDir, "000001_init.down.sql"))
	if err != nil {
		db.Close()
		t.Fatalf("failed to read down migration: %v", err)
	}
	upSQL, err := os.ReadFile(filepath.Join(migrationDir, "000001_init.up.sql"))
	if err != nil {
		db.Close()
		t.Fatalf("failed to read up migration: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := db.ExecContext(ctx, string(downSQL)); err != nil {
		db.Close()
		t.Fatalf("failed to reset test database: %v", err)
	}
	if _, err := db.ExecContext(ctx, string(upSQL)); err != nil {
		db.Close()
		t.Fatalf("failed to migrate test database: %v", err)
	}

	return db, New(db)
}

func assertTestDatabaseURL(t *testing.T, databaseURL string) {
	t.Helper()

	parsed, err := url.Parse(databaseURL)
	if err != nil {
		t.Fatalf("invalid %s: %v", testDatabaseURLEnv, err)
	}
	databaseName := strings.TrimPrefix(parsed.Path, "/")
	if !strings.Contains(strings.ToLower(databaseName), "test") {
		t.Fatalf("refusing to run integration tests against database %q; database name must contain 'test'", databaseName)
	}
}

func TestRepositoriesIntegrationUserAccountCategoryTransactionAndReports(t *testing.T) {
	db, repo := setupIntegrationTestDB(t)
	defer db.Close()

	ctx := context.Background()

	user, err := repo.CreateUser(ctx, "user@example.com", "hashed-password")
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}
	if user.ID == uuid.Nil || user.Email != "user@example.com" {
		t.Fatalf("unexpected user: %+v", user)
	}

	foundUser, err := repo.FindUserByEmail(ctx, "user@example.com")
	if err != nil {
		t.Fatalf("FindUserByEmail returned error: %v", err)
	}
	if foundUser.ID != user.ID {
		t.Fatalf("expected found user id %s, got %s", user.ID, foundUser.ID)
	}

	accountTypes, err := repo.ListAccountTypes(ctx)
	if err != nil {
		t.Fatalf("ListAccountTypes returned error: %v", err)
	}
	if len(accountTypes) != 5 || accountTypes[0].Name != "bank" || accountTypes[4].Name != "stock_broker" {
		t.Fatalf("unexpected account types: %+v", accountTypes)
	}

	bankAccount, err := repo.CreateAccount(ctx, user.ID, "BCA Savings", 1, 5_000_000, nil, nil)
	if err != nil {
		t.Fatalf("CreateAccount bank returned error: %v", err)
	}
	cashAccount, err := repo.CreateAccount(ctx, user.ID, "Cash", 3, 500_000, nil, nil)
	if err != nil {
		t.Fatalf("CreateAccount cash returned error: %v", err)
	}
	if bankAccount.Type != "bank" || bankAccount.Currency != "IDR" || !bankAccount.IsActive {
		t.Fatalf("unexpected bank account: %+v", bankAccount)
	}

	accounts, err := repo.ListAccounts(ctx, user.ID)
	if err != nil {
		t.Fatalf("ListAccounts returned error: %v", err)
	}
	if len(accounts) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(accounts))
	}

	updatedName := "BCA Main"
	updatedAccount, err := repo.UpdateAccount(ctx, user.ID, bankAccount.ID, &updatedName, nil)
	if err != nil {
		t.Fatalf("UpdateAccount returned error: %v", err)
	}
	if updatedAccount.Name != updatedName {
		t.Fatalf("expected updated account name %q, got %q", updatedName, updatedAccount.Name)
	}

	categories, err := repo.ListCategories(ctx, user.ID, "expense")
	if err != nil {
		t.Fatalf("ListCategories default returned error: %v", err)
	}
	if len(categories) == 0 {
		t.Fatal("expected default expense categories")
	}

	customCategory, err := repo.CreateCategory(ctx, user.ID, "Coffee", "expense")
	if err != nil {
		t.Fatalf("CreateCategory returned error: %v", err)
	}
	if customCategory.IsDefault || customCategory.Type != "expense" {
		t.Fatalf("unexpected custom category: %+v", customCategory)
	}

	updatedCategory, err := repo.UpdateCategory(ctx, user.ID, customCategory.ID, "Cafe")
	if err != nil {
		t.Fatalf("UpdateCategory returned error: %v", err)
	}
	if updatedCategory.Name != "Cafe" {
		t.Fatalf("expected updated category name Cafe, got %q", updatedCategory.Name)
	}

	expenseTx, err := repo.CreateTransaction(ctx, domain.Transaction{
		UserID:      user.ID,
		AccountID:   bankAccount.ID,
		CategoryID:  &customCategory.ID,
		Type:        "expense",
		Amount:      50_000,
		Description: "Coffee",
		Date:        time.Date(2026, 6, 3, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("CreateTransaction expense returned error: %v", err)
	}
	if expenseTx.ID == uuid.Nil {
		t.Fatal("expected expense transaction id")
	}

	transferTx, err := repo.CreateTransaction(ctx, domain.Transaction{
		UserID:            user.ID,
		AccountID:         bankAccount.ID,
		TransferAccountID: &cashAccount.ID,
		Type:              "transfer",
		Amount:            100_000,
		Description:       "ATM withdrawal",
		Date:              time.Date(2026, 6, 4, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("CreateTransaction transfer returned error: %v", err)
	}
	if transferTx.ID == uuid.Nil {
		t.Fatal("expected transfer transaction id")
	}

	transactions, err := repo.ListTransactions(ctx, user.ID, "2026-06-01", "2026-06-30", "", "", "", 50, 0)
	if err != nil {
		t.Fatalf("ListTransactions returned error: %v", err)
	}
	if len(transactions) != 2 {
		t.Fatalf("expected 2 transactions, got %d", len(transactions))
	}

	pagedTransactions, err := repo.ListTransactions(ctx, user.ID, "2026-06-01", "2026-06-30", "", "", "", 1, 1)
	if err != nil {
		t.Fatalf("ListTransactions paged returned error: %v", err)
	}
	if len(pagedTransactions) != 1 || pagedTransactions[0].ID != expenseTx.ID {
		t.Fatalf("expected second page to contain expense transaction, got %+v", pagedTransactions)
	}

	netWorth, activeAccounts, err := repo.NetWorth(ctx, user.ID)
	if err != nil {
		t.Fatalf("NetWorth returned error: %v", err)
	}
	if netWorth != 5_450_000 {
		t.Fatalf("expected net worth 5450000, got %f", netWorth)
	}
	if len(activeAccounts) != 2 {
		t.Fatalf("expected 2 active accounts, got %d", len(activeAccounts))
	}

	totalSpending, spendingCategories, err := repo.SpendingByCategory(ctx, user.ID, time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("SpendingByCategory returned error: %v", err)
	}
	if totalSpending != 50_000 {
		t.Fatalf("expected total spending 50000, got %f", totalSpending)
	}
	if len(spendingCategories) != 1 || spendingCategories[0].Name != "Cafe" || spendingCategories[0].Percentage != 100 {
		t.Fatalf("unexpected spending categories: %+v", spendingCategories)
	}

	if err := repo.DeleteCategory(ctx, user.ID, customCategory.ID); err != nil {
		t.Fatalf("DeleteCategory returned error: %v", err)
	}
	if err := repo.DeleteCategory(ctx, user.ID, customCategory.ID); err != response.ErrNotFound {
		t.Fatalf("expected ErrNotFound when deleting category twice, got %v", err)
	}

	if err := repo.SoftDeleteAccount(ctx, user.ID, cashAccount.ID); err != nil {
		t.Fatalf("SoftDeleteAccount returned error: %v", err)
	}
	_, activeAccounts, err = repo.NetWorth(ctx, user.ID)
	if err != nil {
		t.Fatalf("NetWorth after soft delete returned error: %v", err)
	}
	if len(activeAccounts) != 1 {
		t.Fatalf("expected 1 active account after soft delete, got %d", len(activeAccounts))
	}
}

func TestRepositoriesIntegrationNotFoundPaths(t *testing.T) {
	db, repo := setupIntegrationTestDB(t)
	defer db.Close()

	ctx := context.Background()
	if _, err := repo.FindUserByEmail(ctx, "missing@example.com"); err != response.ErrNotFound {
		t.Fatalf("expected ErrNotFound for missing user, got %v", err)
	}

	user, err := repo.CreateUser(ctx, "user@example.com", "hashed-password")
	if err != nil {
		t.Fatalf("CreateUser returned error: %v", err)
	}
	if _, err := repo.UpdateAccount(ctx, user.ID, uuid.New(), nil, nil); err != response.ErrNotFound {
		t.Fatalf("expected ErrNotFound for missing account update, got %v", err)
	}
	if err := repo.SoftDeleteAccount(ctx, user.ID, uuid.New()); err != response.ErrNotFound {
		t.Fatalf("expected ErrNotFound for missing account delete, got %v", err)
	}
}
