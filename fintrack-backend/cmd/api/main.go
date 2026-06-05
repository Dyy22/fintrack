package main

import (
	"log/slog"
	"os"

	"fintrack-backend/internal/config"
	httpHandler "fintrack-backend/internal/handler/http"
	"fintrack-backend/internal/platform/database"
	"fintrack-backend/internal/platform/gold"
	"fintrack-backend/internal/platform/security"
	"fintrack-backend/internal/platform/stock"
	"fintrack-backend/internal/repository/postgres"
	"fintrack-backend/internal/usecase"
)

func main() {
	cfg := config.Load()

	logLevel := slog.LevelInfo
	if cfg.IsProduction() {
		logLevel = slog.LevelWarn
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel, AddSource: true})))

	if err := cfg.Validate(); err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}
	slog.Info("starting fintrack api", "env", cfg.AppEnv, "port", cfg.Port, "database_url_configured", cfg.HasDatabaseURL())

	db, err := database.Open(cfg.DatabaseURL())
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("database connected")

	jwtService := security.NewJWTService(cfg.JWTSecret, cfg.JWTExpiresIn)
	repo := postgres.New(db)
	goldProvider := gold.NewProvider(cfg.GoldPriceSourceURL, cfg.GoldPriceFallbackPerGram)
	stockProvider := stock.NewProvider()
	uc := usecase.New(repo, jwtService).
		WithGoldPriceProvider(goldProvider, cfg.GoldPriceRefreshInterval).
		WithStockMarketProvider(stockProvider)
	handler := httpHandler.New(uc)
	router := httpHandler.Router(cfg, handler, jwtService, db)

	slog.Info("server starting", "port", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
