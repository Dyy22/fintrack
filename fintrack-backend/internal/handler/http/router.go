package http

import (
	"time"

	"fintrack-backend/internal/config"
	"fintrack-backend/internal/middleware"
	"fintrack-backend/internal/platform/security"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router(cfg config.Config, h *Handler, jwtService security.JWTService) *gin.Engine {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api/v1")
	api.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

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
