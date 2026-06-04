package http

import (
	"context"
	stdhttp "net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func healthHandler(pinger HealthPinger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if pinger == nil {
			c.JSON(stdhttp.StatusOK, gin.H{
				"status":   "ok",
				"database": "unknown",
			})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := pinger.PingContext(ctx); err != nil {
			c.JSON(stdhttp.StatusServiceUnavailable, gin.H{
				"status":   "degraded",
				"database": "error",
			})
			return
		}

		c.JSON(stdhttp.StatusOK, gin.H{
			"status":   "ok",
			"database": "ok",
		})
	}
}
