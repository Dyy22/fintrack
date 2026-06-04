package security

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTServiceGenerateAndParse(t *testing.T) {
	service := NewJWTService("test-secret", time.Hour)
	userID := uuid.New()
	email := "user@example.com"

	token, err := service.Generate(userID, email)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := service.Parse(token)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if claims.UserID != userID.String() {
		t.Fatalf("expected user id %s, got %s", userID.String(), claims.UserID)
	}
	if claims.Email != email {
		t.Fatalf("expected email %s, got %s", email, claims.Email)
	}
	if claims.ExpiresAt == nil || !claims.ExpiresAt.After(time.Now()) {
		t.Fatal("expected token expiry to be in the future")
	}
}

func TestJWTServiceParseRejectsWrongSecret(t *testing.T) {
	service := NewJWTService("test-secret", time.Hour)
	otherService := NewJWTService("other-secret", time.Hour)

	token, err := service.Generate(uuid.New(), "user@example.com")
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if _, err := otherService.Parse(token); err == nil {
		t.Fatal("expected token parse to fail with wrong secret")
	}
}
