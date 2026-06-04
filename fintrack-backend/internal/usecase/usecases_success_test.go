package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"fintrack-backend/internal/domain"
	"fintrack-backend/internal/platform/response"
	"fintrack-backend/internal/platform/security"

	"github.com/google/uuid"
)

func TestRegisterNormalizesEmailAndHashesPassword(t *testing.T) {
	userID := uuid.New()
	repo := fakeRepository{
		createUserFn: func(_ context.Context, email, passwordHash string) (domain.User, error) {
			if email != "user@example.com" {
				t.Fatalf("expected normalized email, got %q", email)
			}
			if passwordHash == "securepassword" || passwordHash == "" {
				t.Fatalf("expected bcrypt hash, got %q", passwordHash)
			}
			if !security.CheckPassword(passwordHash, "securepassword") {
				t.Fatal("expected hash to match original password")
			}
			return domain.User{ID: userID, Email: email, PasswordHash: passwordHash}, nil
		},
	}
	uc := New(repo, security.NewJWTService("secret", time.Hour))

	user, err := uc.Register(context.Background(), "  User@Example.COM ", "securepassword")
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if user.ID != userID || user.Email != "user@example.com" {
		t.Fatalf("unexpected user: %+v", user)
	}
}

func TestRegisterMapsDuplicateEmailToConflict(t *testing.T) {
	repo := fakeRepository{
		createUserFn: func(_ context.Context, _, _ string) (domain.User, error) {
			return domain.User{}, errors.New("duplicate key value violates unique constraint")
		},
	}
	uc := New(repo, security.NewJWTService("secret", time.Hour))

	_, err := uc.Register(context.Background(), "user@example.com", "securepassword")
	if !errors.Is(err, response.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}

func TestRegisterRejectsShortPassword(t *testing.T) {
	uc := &Usecases{}

	_, err := uc.Register(context.Background(), "user@example.com", "short")
	if !errors.Is(err, response.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}

func TestLoginNormalizesEmailAndReturnsToken(t *testing.T) {
	passwordHash, err := security.HashPassword("securepassword")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	userID := uuid.New()
	repo := fakeRepository{
		findUserByEmailFn: func(_ context.Context, email string) (domain.User, error) {
			if email != "user@example.com" {
				t.Fatalf("expected normalized email, got %q", email)
			}
			return domain.User{ID: userID, Email: email, PasswordHash: passwordHash}, nil
		},
	}
	jwtService := security.NewJWTService("secret", time.Hour)
	uc := New(repo, jwtService)

	result, err := uc.Login(context.Background(), " User@Example.COM ", "securepassword")
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	if result.Token == "" {
		t.Fatal("expected token")
	}
	claims, err := jwtService.Parse(result.Token)
	if err != nil {
		t.Fatalf("generated token should parse: %v", err)
	}
	if claims.UserID != userID.String() || claims.Email != "user@example.com" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestLoginRejectsWrongPassword(t *testing.T) {
	passwordHash, err := security.HashPassword("securepassword")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	repo := fakeRepository{
		findUserByEmailFn: func(_ context.Context, email string) (domain.User, error) {
			return domain.User{ID: uuid.New(), Email: email, PasswordHash: passwordHash}, nil
		},
	}
	uc := New(repo, security.NewJWTService("secret", time.Hour))

	_, err = uc.Login(context.Background(), "user@example.com", "wrongpassword")
	if !errors.Is(err, response.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestListAccountTypesDelegatesToRepository(t *testing.T) {
	repo := fakeRepository{
		listAccountTypesFn: func(_ context.Context) ([]domain.AccountType, error) {
			return []domain.AccountType{{ID: 1, Name: "bank", Description: "Bank account"}}, nil
		},
	}
	uc := New(repo, security.NewJWTService("secret", time.Hour))

	accountTypes, err := uc.ListAccountTypes(context.Background())
	if err != nil {
		t.Fatalf("ListAccountTypes returned error: %v", err)
	}
	if len(accountTypes) != 1 || accountTypes[0].Name != "bank" {
		t.Fatalf("unexpected account types: %+v", accountTypes)
	}
}

func TestCreateAccountTrimsNameAndDelegates(t *testing.T) {
	userID := uuid.New()
	repo := fakeRepository{
		createAccountFn: func(_ context.Context, gotUserID uuid.UUID, name string, accountTypeID int, balance float64, goldGrams *float64, goldPrice *float64) (domain.Account, error) {
			if gotUserID != userID || name != "BCA Savings" || accountTypeID != 1 || balance != 5000000 {
				t.Fatalf("unexpected args: userID=%s name=%q accountTypeID=%d balance=%f", gotUserID, name, accountTypeID, balance)
			}
			return domain.Account{ID: uuid.New(), UserID: gotUserID, Name: name, AccountTypeID: accountTypeID, Balance: balance}, nil
		},
	}
	uc := New(repo, security.NewJWTService("secret", time.Hour))

	account, err := uc.CreateAccount(context.Background(), userID, "  BCA Savings  ", 1, 5000000, nil)
	if err != nil {
		t.Fatalf("CreateAccount returned error: %v", err)
	}
	if account.Name != "BCA Savings" {
		t.Fatalf("expected trimmed account name, got %q", account.Name)
	}
}

func TestCreateCategoryMapsDuplicateToConflict(t *testing.T) {
	repo := fakeRepository{
		createCategoryFn: func(_ context.Context, _ uuid.UUID, _, _ string) (domain.Category, error) {
			return domain.Category{}, errors.New("duplicate key value violates unique constraint")
		},
	}
	uc := New(repo, security.NewJWTService("secret", time.Hour))

	_, err := uc.CreateCategory(context.Background(), uuid.New(), "Food", "expense")
	if !errors.Is(err, response.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}

func TestCreateTransactionValidExpenseSetsDateAndDelegates(t *testing.T) {
	userID := uuid.New()
	accountID := uuid.New()
	categoryID := uuid.New()
	repo := fakeRepository{
		createTransactionFn: func(_ context.Context, tx domain.Transaction) (domain.Transaction, error) {
			if tx.UserID != userID || tx.AccountID != accountID || tx.CategoryID == nil || *tx.CategoryID != categoryID {
				t.Fatalf("unexpected transaction ownership: %+v", tx)
			}
			if tx.Type != "expense" || tx.Amount != 50000 || tx.Description != "Lunch" {
				t.Fatalf("unexpected transaction details: %+v", tx)
			}
			if tx.Date.IsZero() {
				t.Fatal("expected usecase to set default date")
			}
			tx.ID = uuid.New()
			return tx, nil
		},
	}
	uc := New(repo, security.NewJWTService("secret", time.Hour))

	tx, err := uc.CreateTransaction(context.Background(), domain.Transaction{UserID: userID, AccountID: accountID, CategoryID: &categoryID, Type: "expense", Amount: 50000, Description: "Lunch"})
	if err != nil {
		t.Fatalf("CreateTransaction returned error: %v", err)
	}
	if tx.ID == uuid.Nil {
		t.Fatal("expected created transaction id")
	}
}

func TestListTransactionsRejectsInvalidType(t *testing.T) {
	uc := &Usecases{}

	_, err := uc.ListTransactions(context.Background(), uuid.New(), "", "", "", "", "refund", 50, 0)
	if !errors.Is(err, response.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}

func TestListTransactionsNormalizesPaginationAndDelegates(t *testing.T) {
	userID := uuid.New()
	repo := fakeRepository{
		listTransactionsFn: func(_ context.Context, gotUserID uuid.UUID, start, end, accountID, categoryID, typ string, limit, offset int) ([]domain.Transaction, error) {
			if gotUserID != userID || typ != "expense" || limit != 100 || offset != 10 {
				t.Fatalf("unexpected args: userID=%s type=%q limit=%d offset=%d", gotUserID, typ, limit, offset)
			}
			return []domain.Transaction{{ID: uuid.New(), Type: "expense", Amount: 50000}}, nil
		},
	}
	uc := New(repo, security.NewJWTService("secret", time.Hour))

	transactions, err := uc.ListTransactions(context.Background(), userID, "", "", "", "", "expense", 500, 10)
	if err != nil {
		t.Fatalf("ListTransactions returned error: %v", err)
	}
	if len(transactions) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(transactions))
	}
}

func TestListTransactionsRejectsInvalidPagination(t *testing.T) {
	uc := &Usecases{}

	_, err := uc.ListTransactions(context.Background(), uuid.New(), "", "", "", "", "", -1, 0)
	if !errors.Is(err, response.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest for negative limit, got %v", err)
	}

	_, err = uc.ListTransactions(context.Background(), uuid.New(), "", "", "", "", "", 50, -1)
	if !errors.Is(err, response.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest for negative offset, got %v", err)
	}
}

func TestSpendingByCategoryDelegatesWithExclusiveEndAndReturnsInclusiveEnd(t *testing.T) {
	userID := uuid.New()
	repo := fakeRepository{
		spendingByCategoryFn: func(_ context.Context, gotUserID uuid.UUID, start, end time.Time) (float64, []domain.SpendingCategory, float64, error) {
			if gotUserID != userID {
				t.Fatalf("expected userID %s, got %s", userID, gotUserID)
			}
			expectedStart := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
			expectedExclusiveEnd := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
			if !start.Equal(expectedStart) || !end.Equal(expectedExclusiveEnd) {
				t.Fatalf("unexpected range: start=%s end=%s", start, end)
			}
			return 150000, []domain.SpendingCategory{{Name: "Food", Amount: 150000, Percentage: 100}}, 200000, nil
		},
	}
	uc := New(repo, security.NewJWTService("secret", time.Hour))

	start, end, total, items, totalIncome, err := uc.SpendingByCategory(context.Background(), userID, "2026-06-01", "2026-06-30")
	if err != nil {
		t.Fatalf("SpendingByCategory returned error: %v", err)
	}
	if start.Format("2006-01-02") != "2026-06-01" || end.Format("2006-01-02") != "2026-06-30" {
		t.Fatalf("expected inclusive output range, got %s to %s", start, end)
	}
	if total != 150000 || len(items) != 1 || items[0].Name != "Food" || totalIncome != 200000 {
		t.Fatalf("unexpected report result: total=%f items=%+v totalIncome=%f", total, items, totalIncome)
	}
}
