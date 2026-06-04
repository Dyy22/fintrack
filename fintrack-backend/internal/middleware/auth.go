package middleware

import (
	"strings"

	"fintrack-backend/internal/platform/response"
	"fintrack-backend/internal/platform/security"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const userIDKey = "user_id"

func Auth(jwtService security.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, response.ErrUnauthorized)
			c.Abort()
			return
		}

		claims, err := jwtService.Parse(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			response.Error(c, response.ErrUnauthorized)
			c.Abort()
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			response.Error(c, response.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set(userIDKey, userID)
		c.Next()
	}
}

func UserID(c *gin.Context) (uuid.UUID, bool) {
	value, ok := c.Get(userIDKey)
	if !ok {
		return uuid.Nil, false
	}
	userID, ok := value.(uuid.UUID)
	return userID, ok
}
