package middleware

import (
	"context"
	"monitoring-system/server/pkg/app_error"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		if ctx.Err() == context.DeadlineExceeded {
			c.Error(app_error.NewApiError(http.StatusRequestTimeout, "Request timeout"))
			c.Abort()
		}
	}
}
