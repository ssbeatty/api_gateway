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
// @Param user_name query int false  "用户名"
// @Param password query int false  "密码"
// @Tags register
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
			c.ResponseError(err.Error())
			return
		}
		successData, err := models.SetUserJwtToken(user.Id, user.Username)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(successData)
	}
}

// RegisterTenant
// @Summary 注册租户
// @Description 注册租户
// @Param user_name query int false  "用户名"
// @Param password query int false  "密码"
// @Tags register
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
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(user)
	}
}

// OauthLoginPassword 账号密码登录
// @ID OauthLoginPassword
// @Summary 账号密码登录
// @Description 账号密码登录
// @Param user_name query int false  "用户名"
// @Param password query int false  "密码"
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
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(successData)
	}
}

// EndpointsList
// @Summary 获取所有路由配置
// @Description 获取所有路由配置
// @Param page_num query int false  "页码数"
// @Param page_size query int false  "分页尺寸" default(20)
// @Tags endpoints
// @Accept x-www-form-urlencoded
// @Success 200 {object} string
// @Router /endpoints/list [post]
func (s *Service) EndpointsList(c *Context) {
	var form payload.Pages
	err := c.ShouldBind(&form)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		Data, err := models.GetEndPointList(form.PageNum, form.PageSize)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(Data)
	}
}

// EndpointsAdd
// @Summary 新增/修改路由配置
// @Description 新增/修改路由配置（未携带id信息为新增）
// @Tags endpoints
// @Accept x-www-form-urlencoded
// @Success 200 {object} string
// @Router /endpoints/add [post]
func (s *Service) EndpointsAdd(c *Context) {
	var form payload.EndPoint
	err := c.ShouldBind(&form)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		Data, err := models.EndPointUpsert(form)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(Data)
	}
}

// EndpointsDelete
// @Summary 删除路由配置
// @Description 删除路由配置
// @Param id query int false  "网关配置的id"
// @Tags endpoints
// @Accept x-www-form-urlencoded
// @Success 200 {object} string
// @Router /endpoints/delete [post]
func (s *Service) EndpointsDelete(c *Context) {
	var form payload.EndPoint
	err := c.ShouldBind(&form)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		if err := c.ShouldBindJSON(&form); err != nil {
			c.ResponseError(err.Error())
			return
		}
		Data, err := models.DeleteEndPointById(form.Id)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(Data)
	}
}

// EndpointsGetDetail
// @Summary 获取路由配置详情
// @Description 获取路由配置详情
// @Param id query int false  "网关配置的id"
// @Tags endpoints
// @Accept x-www-form-urlencoded
// @Success 200 {object} string
// @Router /endpoints/detail [post]
func (s *Service) EndpointsGetDetail(c *Context) {
	var form payload.EndPoint
	err := c.ShouldBind(&form)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		if err := c.ShouldBindJSON(&form); err != nil {
			c.ResponseError(err.Error())
			return
		}
		Data, err := models.GetEndPointById(form.Id)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(Data)
	}
}
