package http

import (
	"github.com/gin-gonic/gin"
	"net/http"

)

func JSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

func Error(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}