package middleware

import (
	"auth-api/src/pkg/app_error"
	"auth-api/src/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors[0]
			switch e := err.Err.(type) {
			case *app_error.ApiError:
				c.JSON(e.StatusCode, e)
			default:
				log.Error("Error occurred %v", e)
				c.JSON(http.StatusInternalServerError, map[string]string{"message": e.Error()})
			}
			c.Abort()
		}
	}
}
