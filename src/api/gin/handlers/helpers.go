package handlers

import (
	"context"
	"io"
	"monitoring-system/server/src/pkg/app_error"
	"net/http"

	"github.com/gin-gonic/gin"
)

func bindJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		if err == io.EOF {
			return app_error.NewApiError(400, "Invalid request")
		}
		return err
	}
	return nil
}

func bindQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return app_error.NewApiError(400, "Invalid request")
	}
	return nil
}

func processRequest[T any, U any](c *gin.Context, input T, executeFunc func(context.Context, T) (U, error)) {
	if err := bindJSON(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := executeFunc(c.Request.Context(), input)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, output)
}

func processRequestNoOutput[T any](c *gin.Context, input T, executeFunc func(context.Context, T) error) {
	if err := bindJSON(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := executeFunc(c.Request.Context(), input); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

func processRequestQuery[T any, U any](c *gin.Context, input T, executeFunc func(context.Context, T) (U, error)) {
	if err := bindQuery(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := executeFunc(c.Request.Context(), input)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, output)
}
