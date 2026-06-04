package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppEnv                   string
	Port                     string
	DBHost                   string
	DBPort                   string
	DBName                   string
	DBUser                   string
	DBPassword               string
	JWTSecret                string
	JWTExpiresIn             time.Duration
	CORSAllowedOrigins       []string
	GoldPriceSourceURL       string
	GoldPriceFallbackPerGram float64
	GoldPriceRefreshInterval time.Duration
}

func Load() Config {
	jwtExpiresIn, err := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "24h"))
	if err != nil {
		jwtExpiresIn = 24 * time.Hour
	}

	goldPriceRefreshInterval, err := time.ParseDuration(getEnv("GOLD_PRICE_REFRESH_INTERVAL", "1h"))
	if err != nil {
		goldPriceRefreshInterval = time.Hour
	}

	return Config{
		AppEnv:                   getEnv("APP_ENV", "development"),
		Port:                     getEnv("PORT", "8080"),
		DBHost:                   getEnv("DB_HOST", "localhost"),
		DBPort:                   getEnv("DB_PORT", "5432"),
		DBName:                   getEnv("DB_NAME", "fintrack"),
		DBUser:                   getEnv("DB_USER", "fintrack"),
		DBPassword:               getEnv("DB_PASSWORD", "fintrack"),
		JWTSecret:                getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpiresIn:             jwtExpiresIn,
		CORSAllowedOrigins:       splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")),
		GoldPriceSourceURL:       getEnv("GOLD_PRICE_SOURCE_URL", "https://logam-mulia-api.iamutaki.workers.dev/api/prices/logammulia"),
		GoldPriceFallbackPerGram: getEnvFloat("GOLD_PRICE_FALLBACK_PER_GRAM", 0),
		GoldPriceRefreshInterval: goldPriceRefreshInterval,
	}
}

func (c Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func (c Config) IsProduction() bool {
	return strings.EqualFold(c.AppEnv, "production")
}

func (c Config) PortInt() int {
	port, err := strconv.Atoi(c.Port)
	if err != nil {
		return 8080
	}
	return port
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getEnvFloat(key string, fallback float64) float64 {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}
	return parsed
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
