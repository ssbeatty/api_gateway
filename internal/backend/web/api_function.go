package web

import (
	"api_gateway/internal/backend/models"
	"api_gateway/internal/backend/payload"
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

// RegisterAdmin
// @Summary 注册管理员用户
// @Description 注册管理员用户
// @Tags register_admin
// @Accept x-www-form-urlencoded
// @Success 200 {object} string
// @Router /register/admin [post]
func (s *Service) RegisterAdmin(c *Context) {
	var form payload.RegisterUser
	err := c.ShouldBind(&form)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		if err := c.ShouldBindJSON(&form); err != nil {
			c.ResponseError(err.Error())
			return
		}
		user, err := models.InsertAdmin(form.UserName, form.Password)
		if err != nil {
			return
		}
		successData, err := models.SetLoginJwtToken(user.Id, user.Username)
		if err != nil {
			return
		}
		c.ResponseOk(successData)
	}
}

// RegisterTenant
// @Summary 注册租户
// @Description 注册租户
// @Tags register_tenant
// @Accept x-www-form-urlencoded
// @Success 200 {object} string
// @Router /register/tenant [post]
func (s *Service) RegisterTenant(c *Context) {
	var form payload.RegisterUser
	err := c.ShouldBind(&form)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		if err := c.ShouldBindJSON(&form); err != nil {
			c.ResponseError(err.Error())
			return
		}
		user, err := models.InsertAdmin(form.UserName, form.Password)
		if err != nil {
			return
		}
		c.ResponseOk(user)
	}
}

// OauthLoginPassword 账号密码登录
// @ID OauthLoginPassword
// @Summary 账号密码登录
// @Description 账号密码登录
// @Tags Oauth
// @Accept x-www-form-urlencoded
// @Success 200 {object} string
// @Router /oauth/login/password [post]
func (s *Service) OauthLoginPassword(c *Context) {
	var form payload.OauthLoginPasswordReq
	err := c.ShouldBind(&form)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		if err := c.ShouldBindJSON(&form); err != nil {
			c.ResponseError(err.Error())
			return
		}
		successData, err := models.OauthLoginPassword(form)
		if err != nil {
			return
		}
		c.ResponseOk(successData)
	}
}
