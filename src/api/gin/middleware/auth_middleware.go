package middleware

import (
	"auth-api/src/internal/modules/user-manager/domain/auth"
	"auth-api/src/pkg/app_error"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	AuthMiddleware(groupNames ...auth.UserGroup) gin.HandlerFunc
}

type AuthMiddlewareImpl struct {
	auth auth.AuthService
}

func NewAuthMiddleware(a auth.AuthService) AuthMiddleware {
	return &AuthMiddlewareImpl{
		auth: a,
	}
}

func (a *AuthMiddlewareImpl) AuthMiddleware(groupNames ...auth.UserGroup) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(app_error.NewApiError(401, "Unauthorized"))
			c.Abort()
			return
		}

		token := authHeader[7:] // remove Bearer from token

		claims, err := a.auth.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.Error(app_error.NewApiError(401, "Unauthorized"))
			c.Abort()
			return
		}

		userGroups := make(map[string]struct{}, len(claims.UserGroups))
		for _, group := range claims.UserGroups {
			userGroups[group] = struct{}{}
		}

		authorized := false
		for _, groupName := range groupNames {
			if _, exists := userGroups[string(groupName)]; exists {
				authorized = true
				break
			}
		}

		if !authorized {
			c.Error(app_error.NewApiError(401, "Unauthorized"))
			c.Abort()
			return
		}

		c.Set("jwtToken", token)
		c.Set("claims", claims)

		c.Next()
	}
}
