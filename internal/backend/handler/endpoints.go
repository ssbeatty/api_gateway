package handler

import (
	"api_gateway/internal/backend/models"
	"api_gateway/internal/backend/payload"
)

// EndpointsList
// @Summary 获取所有路由配置
// @Description 获取所有路由配置
// @Param page_num query int false  "页码数"
// @Param page_size query int false  "分页尺寸" default(20)
// @Tags endpoints
// @Accept x-www-form-urlencoded
// @Success 200 {object} string
// @Router /api/v1/endpoints/list [post]
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
// @Router /api/v1/endpoints/add [post]
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
// @Router /api/v1/endpoints/delete [post]
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
// @Router /api/v1/endpoints/detail [post]
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
