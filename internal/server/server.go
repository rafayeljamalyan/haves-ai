package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"haves/internal/db"
)

// New builds and configures the Gin engine with all routes registered.
func New(database *db.DB) *gin.Engine {
	router := gin.Default()

	h := &handler{db: database}

	router.GET("/health", h.healthCheck)
	router.GET("/ready", h.ready)

	api := router.Group("/api/v1")
	{
		api.GET("/ping", ping)
	}

	return router
}

// handler holds dependencies shared across route handlers.
type handler struct {
	db *db.DB
}

// codeDBUnavailable marks a failed readiness check (database unreachable).
const codeDBUnavailable = 1001

// healthCheck is a liveness probe: the process is up.
func (h *handler) healthCheck(c *gin.Context) {
	respond(c, http.StatusOK, gin.H{"status": "ok"})
}

// ready is a readiness probe: dependencies (the database) are reachable.
func (h *handler) ready(c *gin.Context) {
	if err := h.db.Ping(c.Request.Context()); err != nil {
		respondError(c, http.StatusServiceUnavailable, codeDBUnavailable, err.Error())
		return
	}
	respond(c, http.StatusOK, gin.H{"status": "ready"})
}

func ping(c *gin.Context) {
	respond(c, http.StatusOK, gin.H{"message": "pong"})
}
