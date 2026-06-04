package http

import (
	"context"
	"encoding/json"
	"errors"
	stdhttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"fintrack-backend/internal/config"
	"fintrack-backend/internal/domain"
	"fintrack-backend/internal/platform/response"
	"fintrack-backend/internal/platform/security"
	"fintrack-backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type healthPingerStub struct {
	err error
}

func (p healthPingerStub) PingContext(_ context.Context) error {
	return p.err
}

func testRouter(uc Usecase, jwtService security.JWTService) *gin.Engine {
	return testRouterWithHealthPinger(uc, jwtService, healthPingerStub{})
}

func testRouterWithHealthPinger(uc Usecase, jwtService security.JWTService, healthPinger HealthPinger) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return Router(config.Config{AppEnv: "test", CORSAllowedOrigins: []string{"http://localhost:3000"}}, New(uc), jwtService, healthPinger)
}

func authHeader(t *testing.T, jwtService security.JWTService, userID uuid.UUID) string {
	t.Helper()
	token, err := jwtService.Generate(userID, "user@example.com")
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	return "Bearer " + token
}

func performRequest(r *gin.Engine, method, path, body, authorization string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	r.ServeHTTP(w, req)
	return w
}

func assertJSONContains(t *testing.T, body, key string, want any) {
	t.Helper()
	var data map[string]any
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		t.Fatalf("response is not valid json: %v body=%s", err, body)
	}
	if got := data[key]; got != want {
		t.Fatalf("expected %s=%v, got %v in body %s", key, want, got, body)
	}
}

func assertValidationField(t *testing.T, body, field string) {
	t.Helper()
	var data struct {
		Error  string            `json:"error"`
		Fields map[string]string `json:"fields"`
	}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		t.Fatalf("response is not valid json: %v body=%s", err, body)
	}
	if data.Error != "validation_error" {
		t.Fatalf("expected validation_error, got %q body=%s", data.Error, body)
	}
	if _, ok := data.Fields[field]; !ok {
		t.Fatalf("expected validation field %q, got fields=%+v body=%s", field, data.Fields, body)
	}
}

func TestRouterHealth(t *testing.T) {
	r := testRouter(fakeUsecase{}, security.NewJWTService("secret", time.Hour))

	w := performRequest(r, stdhttp.MethodGet, "/api/v1/health", "", "")

	if w.Code != stdhttp.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", w.Code, w.Body.String())
	}
	assertJSONContains(t, w.Body.String(), "status", "ok")
	assertJSONContains(t, w.Body.String(), "database", "ok")
}

func TestRouterHealthDatabaseError(t *testing.T) {
	r := testRouterWithHealthPinger(fakeUsecase{}, security.NewJWTService("secret", time.Hour), healthPingerStub{err: errors.New("database unavailable")})

	w := performRequest(r, stdhttp.MethodGet, "/api/v1/health", "", "")

	if w.Code != stdhttp.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d body=%s", w.Code, w.Body.String())
	}
	assertJSONContains(t, w.Body.String(), "status", "degraded")
	assertJSONContains(t, w.Body.String(), "database", "error")
}

func TestRegisterSuccess(t *testing.T) {
	userID := uuid.New()
	r := testRouter(fakeUsecase{
		registerFn: func(_ context.Context, email, password string) (domain.User, error) {
			if email != "user@example.com" || password != "securepassword" {
				t.Fatalf("unexpected register args: email=%q password=%q", email, password)
			}
			return domain.User{ID: userID, Email: email}, nil
		},
	}, security.NewJWTService("secret", time.Hour))

	w := performRequest(r, stdhttp.MethodPost, "/api/v1/auth/register", `{"email":"user@example.com","password":"securepassword"}`, "")

	if w.Code != stdhttp.StatusCreated {
		t.Fatalf("expected status 201, got %d body=%s", w.Code, w.Body.String())
	}
	assertJSONContains(t, w.Body.String(), "email", "user@example.com")
}

func TestRegisterValidationAndConflict(t *testing.T) {
	r := testRouter(fakeUsecase{
		registerFn: func(_ context.Context, _, _ string) (domain.User, error) {
			return domain.User{}, response.ErrConflict
		},
	}, security.NewJWTService("secret", time.Hour))

	bad := performRequest(r, stdhttp.MethodPost, "/api/v1/auth/register", `{"email":"invalid","password":"short"}`, "")
	if bad.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected bad request status 400, got %d body=%s", bad.Code, bad.Body.String())
	}
	assertValidationField(t, bad.Body.String(), "email")
	assertValidationField(t, bad.Body.String(), "password")

	conflict := performRequest(r, stdhttp.MethodPost, "/api/v1/auth/register", `{"email":"user@example.com","password":"securepassword"}`, "")
	if conflict.Code != stdhttp.StatusConflict {
		t.Fatalf("expected conflict status 409, got %d body=%s", conflict.Code, conflict.Body.String())
	}
}

func TestLoginSuccessAndUnauthorized(t *testing.T) {
	jwtService := security.NewJWTService("secret", time.Hour)
	r := testRouter(fakeUsecase{
		loginFn: func(_ context.Context, email, password string) (usecase.LoginResult, error) {
			if password == "wrongpassword" {
				return usecase.LoginResult{}, response.ErrUnauthorized
			}
			if email != "user@example.com" || password != "securepassword" {
				t.Fatalf("unexpected login args: email=%q password=%q", email, password)
			}
			return usecase.LoginResult{Token: "token", User: domain.User{ID: uuid.New(), Email: email}}, nil
		},
	}, jwtService)

	ok := performRequest(r, stdhttp.MethodPost, "/api/v1/auth/login", `{"email":"user@example.com","password":"securepassword"}`, "")
	if ok.Code != stdhttp.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", ok.Code, ok.Body.String())
	}
	assertJSONContains(t, ok.Body.String(), "token", "token")

	unauthorized := performRequest(r, stdhttp.MethodPost, "/api/v1/auth/login", `{"email":"user@example.com","password":"wrongpassword"}`, "")
	if unauthorized.Code != stdhttp.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d body=%s", unauthorized.Code, unauthorized.Body.String())
	}
}

func TestProtectedRouteRequiresBearerToken(t *testing.T) {
	r := testRouter(fakeUsecase{}, security.NewJWTService("secret", time.Hour))

	w := performRequest(r, stdhttp.MethodGet, "/api/v1/account-types", "", "")

	if w.Code != stdhttp.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestListAccountTypesSuccess(t *testing.T) {
	jwtService := security.NewJWTService("secret", time.Hour)
	r := testRouter(fakeUsecase{
		listAccountTypesFn: func(_ context.Context) ([]domain.AccountType, error) {
			return []domain.AccountType{{ID: 1, Name: "bank", Description: "Bank account"}}, nil
		},
	}, jwtService)

	w := performRequest(r, stdhttp.MethodGet, "/api/v1/account-types", "", authHeader(t, jwtService, uuid.New()))

	if w.Code != stdhttp.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"account_types"`) || !strings.Contains(w.Body.String(), `"bank"`) {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

func TestCreateAccountValidationAndSuccess(t *testing.T) {
	jwtService := security.NewJWTService("secret", time.Hour)
	userID := uuid.New()
	r := testRouter(fakeUsecase{
		createAccountFn: func(_ context.Context, gotUserID uuid.UUID, name string, accountTypeID int, balance float64, goldGrams *float64) (domain.Account, error) {
			if gotUserID != userID || name != "BCA" || accountTypeID != 1 || balance != 1000000 {
				t.Fatalf("unexpected create account args: userID=%s name=%q type=%d balance=%f", gotUserID, name, accountTypeID, balance)
			}
			return domain.Account{ID: uuid.New(), UserID: gotUserID, Name: name, AccountTypeID: accountTypeID, Type: "bank", Balance: balance, Currency: "IDR", IsActive: true}, nil
		},
	}, jwtService)
	auth := authHeader(t, jwtService, userID)

	bad := performRequest(r, stdhttp.MethodPost, "/api/v1/accounts", `{"name":"BCA"}`, auth)
	if bad.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected bad request status 400, got %d body=%s", bad.Code, bad.Body.String())
	}
	assertValidationField(t, bad.Body.String(), "account_type_id")

	ok := performRequest(r, stdhttp.MethodPost, "/api/v1/accounts", `{"name":"BCA","account_type_id":1,"balance":1000000}`, auth)
	if ok.Code != stdhttp.StatusCreated {
		t.Fatalf("expected status 201, got %d body=%s", ok.Code, ok.Body.String())
	}
	assertJSONContains(t, ok.Body.String(), "name", "BCA")
}

func TestListAndCreateCategories(t *testing.T) {
	jwtService := security.NewJWTService("secret", time.Hour)
	userID := uuid.New()
	r := testRouter(fakeUsecase{
		listCategoriesFn: func(_ context.Context, gotUserID uuid.UUID, typ string) ([]domain.Category, error) {
			if gotUserID != userID || typ != "expense" {
				t.Fatalf("unexpected list categories args: userID=%s type=%q", gotUserID, typ)
			}
			return []domain.Category{{ID: uuid.New(), Name: "Food", Type: "expense", IsDefault: true}}, nil
		},
		createCategoryFn: func(_ context.Context, gotUserID uuid.UUID, name, typ string) (domain.Category, error) {
			if gotUserID != userID || name != "Coffee" || typ != "expense" {
				t.Fatalf("unexpected create category args: userID=%s name=%q type=%q", gotUserID, name, typ)
			}
			return domain.Category{ID: uuid.New(), Name: name, Type: typ}, nil
		},
	}, jwtService)
	auth := authHeader(t, jwtService, userID)

	list := performRequest(r, stdhttp.MethodGet, "/api/v1/categories?type=expense", "", auth)
	if list.Code != stdhttp.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", list.Code, list.Body.String())
	}

	bad := performRequest(r, stdhttp.MethodPost, "/api/v1/categories", `{"name":"Coffee","type":"transfer"}`, auth)
	if bad.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected bad request status 400, got %d body=%s", bad.Code, bad.Body.String())
	}
	assertValidationField(t, bad.Body.String(), "type")

	ok := performRequest(r, stdhttp.MethodPost, "/api/v1/categories", `{"name":"Coffee","type":"expense"}`, auth)
	if ok.Code != stdhttp.StatusCreated {
		t.Fatalf("expected status 201, got %d body=%s", ok.Code, ok.Body.String())
	}
	assertJSONContains(t, ok.Body.String(), "name", "Coffee")
}

func TestListTransactionsPaginationValidationAndSuccess(t *testing.T) {
	jwtService := security.NewJWTService("secret", time.Hour)
	userID := uuid.New()
	r := testRouter(fakeUsecase{
		listTransactionsFn: func(_ context.Context, gotUserID uuid.UUID, start, end, accountID, categoryID, typ string, limit, offset int) ([]domain.Transaction, error) {
			if gotUserID != userID || typ != "expense" || limit != 100 || offset != 10 {
				t.Fatalf("unexpected list transactions args: userID=%s type=%q limit=%d offset=%d", gotUserID, typ, limit, offset)
			}
			return []domain.Transaction{{ID: uuid.New(), Type: "expense", Amount: 50000}}, nil
		},
	}, jwtService)
	auth := authHeader(t, jwtService, userID)

	bad := performRequest(r, stdhttp.MethodGet, "/api/v1/transactions?limit=abc", "", auth)
	if bad.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected bad request status 400, got %d body=%s", bad.Code, bad.Body.String())
	}

	ok := performRequest(r, stdhttp.MethodGet, "/api/v1/transactions?type=expense&limit=500&offset=10", "", auth)
	if ok.Code != stdhttp.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", ok.Code, ok.Body.String())
	}
	assertJSONContains(t, ok.Body.String(), "limit", float64(100))
	assertJSONContains(t, ok.Body.String(), "offset", float64(10))
}

func TestCreateTransactionValidationAndSuccess(t *testing.T) {
	jwtService := security.NewJWTService("secret", time.Hour)
	userID := uuid.New()
	accountID := uuid.New()
	categoryID := uuid.New()
	r := testRouter(fakeUsecase{
		createTransactionFn: func(_ context.Context, tx domain.Transaction) (domain.Transaction, error) {
			if tx.UserID != userID || tx.AccountID != accountID || tx.CategoryID == nil || *tx.CategoryID != categoryID {
				t.Fatalf("unexpected transaction ids: %+v", tx)
			}
			if tx.Type != "expense" || tx.Amount != 50000 || tx.Description != "Lunch" {
				t.Fatalf("unexpected transaction details: %+v", tx)
			}
			tx.ID = uuid.New()
			return tx, nil
		},
	}, jwtService)
	auth := authHeader(t, jwtService, userID)

	bad := performRequest(r, stdhttp.MethodPost, "/api/v1/transactions", `{"account_id":"not-a-uuid","type":"expense","amount":50000}`, auth)
	if bad.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected bad request status 400, got %d body=%s", bad.Code, bad.Body.String())
	}

	body := `{"account_id":"` + accountID.String() + `","category_id":"` + categoryID.String() + `","type":"expense","amount":50000,"description":"Lunch","date":"2026-06-03T12:00:00Z"}`
	ok := performRequest(r, stdhttp.MethodPost, "/api/v1/transactions", body, auth)
	if ok.Code != stdhttp.StatusCreated {
		t.Fatalf("expected status 201, got %d body=%s", ok.Code, ok.Body.String())
	}
	assertJSONContains(t, ok.Body.String(), "type", "expense")
}

func TestReportsSuccess(t *testing.T) {
	jwtService := security.NewJWTService("secret", time.Hour)
	userID := uuid.New()
	r := testRouter(fakeUsecase{
		netWorthFn: func(_ context.Context, gotUserID uuid.UUID) (float64, []domain.Account, error) {
			if gotUserID != userID {
				t.Fatalf("unexpected user id: %s", gotUserID)
			}
			return 1000000, []domain.Account{{ID: uuid.New(), Name: "BCA", Balance: 1000000}}, nil
		},
		spendingByCategoryFn: func(_ context.Context, gotUserID uuid.UUID, startDate, endDate string) (time.Time, time.Time, float64, []domain.SpendingCategory, float64, error) {
			if gotUserID != userID || startDate != "2026-06-01" || endDate != "2026-06-30" {
				t.Fatalf("unexpected spending args: userID=%s start=%q end=%q", gotUserID, startDate, endDate)
			}
			return time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC), 50000, []domain.SpendingCategory{{Name: "Food", Amount: 50000, Percentage: 100}}, 100000, nil
		},
	}, jwtService)
	auth := authHeader(t, jwtService, userID)

	netWorth := performRequest(r, stdhttp.MethodGet, "/api/v1/reports/net-worth", "", auth)
	if netWorth.Code != stdhttp.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", netWorth.Code, netWorth.Body.String())
	}
	assertJSONContains(t, netWorth.Body.String(), "net_worth", float64(1000000))

	spending := performRequest(r, stdhttp.MethodGet, "/api/v1/reports/spending-by-category?start_date=2026-06-01&end_date=2026-06-30", "", auth)
	if spending.Code != stdhttp.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", spending.Code, spending.Body.String())
	}
	assertJSONContains(t, spending.Body.String(), "total_spending", float64(50000))
}
