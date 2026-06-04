package main

import (
	"log"

	"fintrack-backend/internal/config"
	httpHandler "fintrack-backend/internal/handler/http"
	"fintrack-backend/internal/platform/database"
	"fintrack-backend/internal/platform/gold"
	"fintrack-backend/internal/platform/security"
	"fintrack-backend/internal/repository/postgres"
	"fintrack-backend/internal/usecase"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	jwtService := security.NewJWTService(cfg.JWTSecret, cfg.JWTExpiresIn)
	repo := postgres.New(db)
	goldProvider := gold.NewProvider(cfg.GoldPriceSourceURL, cfg.GoldPriceFallbackPerGram)
	uc := usecase.New(repo, jwtService).WithGoldPriceProvider(goldProvider, cfg.GoldPriceRefreshInterval)
	handler := httpHandler.New(uc)
	router := httpHandler.Router(cfg, handler, jwtService)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
