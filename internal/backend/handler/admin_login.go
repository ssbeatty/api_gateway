package handler

import (
	"api_gateway/internal/backend/models"
	"api_gateway/internal/backend/payload"
	"api_gateway/internal/version"
	"net/http"
)

// GetVersion @Summary 获取当前版本
// @Description 获取当前版本
// @Tags system
// @Accept x-www-form-urlencoded
// @Success 200 {object} string
// @Router /version [get]
func (s *Service) GetVersion(c *Context) {
	c.String(http.StatusOK, version.Version)
}

// RegisterAdmin
// @Summary 注册管理员用户
// @Description 注册管理员用户
// @Param data body payload.AdminRegisterReq true "auth info"
// @Tags auth
// @Accept json
// @Success 200 {object} payload.Response
// @Router /auth/admin/register [post]
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

// AdminLoginPassword 账号密码登录
// @ID AdminLoginPassword
// @Summary admin账号密码登录
// @Description admin账号密码登录
// @Param data body payload.AdminLoginPasswordReq true "auth info"
// @Tags auth
// @Accept json
// @Success 200 {object} payload.Response{data=payload.AdminLoginPasswordResp}
// @Router /auth/admin/login [post]
func (s *Service) AdminLoginPassword(c *Context) {
	var req payload.AdminLoginPasswordReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		successData, err := s.OauthLoginPassword(req)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(successData)
	}
}
