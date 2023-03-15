package handler

import (
	"api_gateway/internal/backend/models"
	"api_gateway/internal/backend/payload"
	"api_gateway/internal/version"
	"net/http"
)

func (s *Service) GetVersion(c *Context) {
	c.String(http.StatusOK, version.Version)
}

func (s *Service) RegisterAdmin(c *Context) {
	var req payload.AdminRegisterReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		user, err := models.InsertAdmin(req.UserName, req.Password)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(user)
	}
}
