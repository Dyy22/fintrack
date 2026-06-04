package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestErrorMapsKnownErrorsToStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{name: "bad request", err: ErrBadRequest, wantStatus: http.StatusBadRequest, wantCode: "bad_request"},
		{name: "unauthorized", err: ErrUnauthorized, wantStatus: http.StatusUnauthorized, wantCode: "unauthorized"},
		{name: "forbidden", err: ErrForbidden, wantStatus: http.StatusForbidden, wantCode: "forbidden"},
		{name: "not found", err: ErrNotFound, wantStatus: http.StatusNotFound, wantCode: "not_found"},
		{name: "conflict", err: ErrConflict, wantStatus: http.StatusConflict, wantCode: "conflict"},
		{name: "wrapped conflict", err: errors.Join(errors.New("wrapped"), ErrConflict), wantStatus: http.StatusConflict, wantCode: "conflict"},
	}

	gin.SetMode(gin.TestMode)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			Error(c, tt.err)

			if w.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
			if body := w.Body.String(); !strings.Contains(body, `"error":"`+tt.wantCode+`"`) {
				t.Fatalf("expected error code %q in body, got %s", tt.wantCode, body)
			}
		})
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := map[string]string{
		"Email":             "email",
		"Password":          "password",
		"AccountTypeID":     "account_type_id",
		"TransferAccountID": "transfer_account_id",
	}

	for input, want := range tests {
		if got := toSnakeCase(input); got != want {
			t.Fatalf("expected %q -> %q, got %q", input, want, got)
		}
	}
}

func TestErrorMapsUnknownErrorToInternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Error(c, errors.New("boom"))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
	if body := w.Body.String(); !strings.Contains(body, `"error":"internal_server_error"`) {
		t.Fatalf("expected internal_server_error body, got %s", body)
	}
}
