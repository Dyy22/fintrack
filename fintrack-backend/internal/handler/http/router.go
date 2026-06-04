package http

import (
	"context"
	"log/slog"
	"time"

	"fintrack-backend/internal/config"
	"fintrack-backend/internal/middleware"
	"fintrack-backend/internal/platform/security"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HealthPinger interface {
	PingContext(ctx context.Context) error
}

func Router(cfg config.Config, h *Handler, jwtService security.JWTService, healthPinger HealthPinger) *gin.Engine {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(ginLogger(), gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api/v1")
	api.GET("/health", healthHandler(healthPinger))

	auth := api.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)

	protected := api.Group("")
	protected.Use(middleware.Auth(jwtService))
	protected.GET("/account-types", h.ListAccountTypes)
	protected.GET("/accounts", h.ListAccounts)
	protected.POST("/accounts", h.CreateAccount)
	protected.PUT("/accounts/:id", h.UpdateAccount)
	protected.DELETE("/accounts/:id", h.DeleteAccount)

	protected.GET("/categories", h.ListCategories)
	protected.POST("/categories", h.CreateCategory)
	protected.PUT("/categories/:id", h.UpdateCategory)
	protected.DELETE("/categories/:id", h.DeleteCategory)

	protected.GET("/transactions", h.ListTransactions)
	protected.POST("/transactions", h.CreateTransaction)

	protected.GET("/gold/price", h.GoldPrice)
	protected.GET("/gold/prices/history", h.GoldPriceHistory)

	protected.GET("/reports/net-worth", h.NetWorth)
	protected.GET("/reports/spending-by-category", h.SpendingByCategory)

	return r
}

func ginLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method

		attrs := []slog.Attr{
			slog.Int("status", status),
			slog.String("method", method),
			slog.String("path", path),
			slog.Duration("latency", latency),
			slog.Int("size", c.Writer.Size()),
		}
		if query != "" {
			attrs = append(attrs, slog.String("query", query))
		}

		level := slog.LevelInfo
		msg := "request"
		if status >= 500 {
			level = slog.LevelError
			msg = "request failed"
		} else if status >= 400 {
			level = slog.LevelWarn
			msg = "request warning"
		}
		slog.LogAttrs(c.Request.Context(), level, msg, attrs...)
	}
}
