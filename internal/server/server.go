package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// New builds and configures the Gin engine with all routes registered.
func New() *gin.Engine {
	router := gin.Default()

	router.GET("/health", healthCheck)

	api := router.Group("/api/v1")
	{
		api.GET("/ping", ping)
	}

	return router
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
