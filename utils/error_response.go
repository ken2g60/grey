package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"message": message,
		"status":  http.StatusBadRequest,
	})
}
