package web

import (
	"api_gateway/internal/version"
	"net/http"
)

// @BasePath /api/v1

// GetVersion
// @Summary 获取当前版本
// @Description 获取当前版本
// @Tags system
// @Accept x-www-form-urlencoded
// @Success 200 {object} string
// @Router /version [get]
func (s *Service) GetVersion(c *Context) {
	c.String(http.StatusOK, version.Version)
}
