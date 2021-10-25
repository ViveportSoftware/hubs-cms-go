package handler

import (
	"hubs-cms-go/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// VersionHandler is version checker API
// @Success 200 {string} string "1.0.0"
// @Router /version [get]
func VersionHandler(c *gin.Context) {
	version := config.EnvVariable.Version
	c.String(http.StatusOK, version)
}
