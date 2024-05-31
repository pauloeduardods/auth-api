package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		done := make(chan struct{})

		go func() {
			c.Request = c.Request.WithContext(ctx)
			c.Next()
			done <- struct{}{}
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				c.JSON(http.StatusRequestTimeout, gin.H{"message": "Request timeout"})
				c.Abort()
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
				c.Abort()
				return
			}
		}
	}
}
