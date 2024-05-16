package middleware

import (
	"monitoring-system/server/pkg/app_error"
	"monitoring-system/server/pkg/jwtToken"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	AuthMiddleware() gin.HandlerFunc
}

type AuthMiddlewareImpl struct {
	JwtToken *jwtToken.JwtToken
}

func NewAuthMiddleware(j *jwtToken.JwtToken) AuthMiddleware {
	return &AuthMiddlewareImpl{
		JwtToken: j,
	}
}

func (a *AuthMiddlewareImpl) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			c.Error(app_error.NewApiError(http.StatusUnauthorized, "Authorization token missing"))
			c.Abort()
			return
		}

		splitToken := strings.Split(token, " ")

		//TODO: Check this
		a.JwtToken.CacheJWK()

		jwtToken, err := a.JwtToken.ParseJWT(splitToken[1])

		if err != nil {
			c.Error(app_error.NewApiError(http.StatusUnauthorized, err.Error()))
			c.Abort()
			return
		}

		c.Set("jwtToken", jwtToken)
		c.Set("user", jwtToken.Claims)

		c.Next()
	}
}
