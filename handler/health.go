package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler is health checker API
// @Success 200 {string} string "ok"
// @Router /health [get]
func HealthHandler(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
