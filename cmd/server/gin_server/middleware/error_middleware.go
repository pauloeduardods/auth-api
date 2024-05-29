package middleware

import (
	"monitoring-system/server/pkg/app_error"
	"monitoring-system/server/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			switch e := err.Err.(type) {
			case *app_error.ApiError:
				c.AbortWithStatusJSON(e.StatusCode, e)
				c.Abort()
			default:
				log.Error("Error occurred %v", e)
				c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{"message": e.Error()})
				c.Abort()
			}
		}
	}
}
