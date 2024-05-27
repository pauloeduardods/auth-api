package middleware

import (
	"monitoring-system/server/domain/auth"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	AuthMiddleware() gin.HandlerFunc
}

type AuthMiddlewareImpl struct {
	auth auth.Auth
}

func NewAuthMiddleware(a auth.Auth) AuthMiddleware {
	return &AuthMiddlewareImpl{
		auth: a,
	}
}

func (a *AuthMiddlewareImpl) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")

		token := auth[7:] // remove Bearer from token

		claims, err := a.auth.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		c.Set("jwtToken", token)
		c.Set("user", claims)

		c.Next()
	}
}
