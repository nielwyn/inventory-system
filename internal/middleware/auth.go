package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nielwyn/inventory-system/internal/service"
	"github.com/nielwyn/inventory-system/pkg/logger"
	"github.com/nielwyn/inventory-system/pkg/response"
	"go.uber.org/zap"
)

// Auth middleware validates JWT tokens
func Auth(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, 401, "Authorization header required")
			c.Abort()
			return
		}

		// Extract token (format: "Bearer <token>")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, 401, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			logger.Error("Token validation failed", zap.Error(err))
			response.Error(c, 401, "Invalid or expired token")
			c.Abort()
			return
		}

		// Extract user ID from token
		userID, err := authService.GetUserFromToken(token)
		if err != nil {
			logger.Error("Failed to extract user from token", zap.Error(err))
			response.Error(c, 401, "Invalid token claims")
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("user_id", userID)
		c.Next()
	}
}
