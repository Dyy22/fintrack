package config

import (
	"testing"
	"time"
)

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("APP_ENV", "")
	t.Setenv("PORT", "")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("JWT_SECRET", "")
	t.Setenv("JWT_EXPIRES_IN", "")
	t.Setenv("CORS_ALLOWED_ORIGINS", "")

	cfg := Load()

	if cfg.AppEnv != "development" || cfg.Port != "8080" {
		t.Fatalf("unexpected default app config: %+v", cfg)
	}
	if cfg.DBHost != "localhost" || cfg.DBPort != "5432" || cfg.DBName != "fintrack" || cfg.DBUser != "fintrack" || cfg.DBPassword != "fintrack" {
		t.Fatalf("unexpected default db config: %+v", cfg)
	}
	if cfg.JWTSecret != "change-me-in-production" || cfg.JWTExpiresIn != 24*time.Hour {
		t.Fatalf("unexpected default jwt config: %+v", cfg)
	}
	if len(cfg.CORSAllowedOrigins) != 2 || cfg.CORSAllowedOrigins[0] != "http://localhost:3000" || cfg.CORSAllowedOrigins[1] != "http://localhost:5173" {
		t.Fatalf("unexpected default cors origins: %+v", cfg.CORSAllowedOrigins)
	}
}

func TestLoadUsesEnvironmentValues(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("PORT", "9090")
	t.Setenv("DB_HOST", "postgres")
	t.Setenv("DB_PORT", "15432")
	t.Setenv("DB_NAME", "custom")
	t.Setenv("DB_USER", "custom_user")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("JWT_SECRET", "jwt-secret")
	t.Setenv("JWT_EXPIRES_IN", "2h")
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://app.example.com, http://localhost:3000 ")

	cfg := Load()

	if !cfg.IsProduction() {
		t.Fatal("expected production mode")
	}
	if cfg.PortInt() != 9090 {
		t.Fatalf("expected port 9090, got %d", cfg.PortInt())
	}
	if cfg.DatabaseURL() != "postgres://custom_user:secret@postgres:15432/custom?sslmode=disable" {
		t.Fatalf("unexpected database URL: %s", cfg.DatabaseURL())
	}
	if cfg.JWTExpiresIn != 2*time.Hour {
		t.Fatalf("expected 2h jwt expiry, got %s", cfg.JWTExpiresIn)
	}
	if len(cfg.CORSAllowedOrigins) != 2 || cfg.CORSAllowedOrigins[0] != "https://app.example.com" || cfg.CORSAllowedOrigins[1] != "http://localhost:3000" {
		t.Fatalf("unexpected cors origins: %+v", cfg.CORSAllowedOrigins)
	}
}

func TestLoadUsesDatabaseURLWhenProvided(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://railway:secret@postgres.railway.internal:5432/railway")
	t.Setenv("DB_HOST", "ignored")
	t.Setenv("DB_PORT", "15432")
	t.Setenv("DB_NAME", "ignored")
	t.Setenv("DB_USER", "ignored")
	t.Setenv("DB_PASSWORD", "ignored")

	cfg := Load()

	if cfg.DatabaseURL() != "postgres://railway:secret@postgres.railway.internal:5432/railway" {
		t.Fatalf("unexpected database URL: %s", cfg.DatabaseURL())
	}
}

func TestValidateRequiresDatabaseURLInProduction(t *testing.T) {
	cfg := Config{AppEnv: "production"}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected production config without DATABASE_URL to be invalid")
	}

	cfg.DatabaseURLValue = "postgres://user:secret@example.com:5432/db?sslmode=require"
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected production config with DATABASE_URL to be valid: %v", err)
	}
}

func TestLoadFallsBackWhenDurationOrPortInvalid(t *testing.T) {
	t.Setenv("JWT_EXPIRES_IN", "invalid")
	t.Setenv("PORT", "invalid")

	cfg := Load()

	if cfg.JWTExpiresIn != 24*time.Hour {
		t.Fatalf("expected fallback jwt expiry, got %s", cfg.JWTExpiresIn)
	}
	if cfg.PortInt() != 8080 {
		t.Fatalf("expected fallback port 8080, got %d", cfg.PortInt())
	}
}
