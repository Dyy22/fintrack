package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fintrack-backend/internal/platform/security"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestAuthAcceptsValidBearerTokenAndSetsUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtService := security.NewJWTService("secret", time.Hour)
	userID := uuid.New()
	token, err := jwtService.Generate(userID, "user@example.com")
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	r := gin.New()
	r.Use(Auth(jwtService))
	r.GET("/protected", func(c *gin.Context) {
		gotUserID, ok := UserID(c)
		if !ok {
			t.Fatal("expected user id in context")
		}
		if gotUserID != userID {
			t.Fatalf("expected userID %s, got %s", userID, gotUserID)
		}
		c.Status(http.StatusNoContent)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthRejectsMissingMalformedOrInvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtService := security.NewJWTService("secret", time.Hour)

	tests := []struct {
		name   string
		header string
	}{
		{name: "missing header", header: ""},
		{name: "malformed header", header: "Token abc"},
		{name: "invalid token", header: "Bearer invalid-token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(Auth(jwtService))
			r.GET("/protected", func(c *gin.Context) {
				t.Fatal("handler should not be called")
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			r.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Fatalf("expected status 401, got %d body=%s", w.Code, w.Body.String())
			}
		})
	}
}
