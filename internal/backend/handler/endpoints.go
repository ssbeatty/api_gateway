package handler

import (
	"api_gateway/internal/backend/models"
	"api_gateway/internal/backend/payload"
)

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
// @Summary 新增路由配置
// @Description 新增路由配置
// @Param data body payload.PostEndPointReq true "endpoint info"
// @Tags endpoints
// @Accept json
// @Success 200 {object} string
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router /api/v1/endpoints [post]
func (s *Service) EndpointsAdd(c *Context) {
	var form payload.PostEndPointReq
	err := c.ShouldBind(&form)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		//Data, err := models.EndPointUpsert(form)
		//if err != nil {
		//	c.ResponseError(err.Error())
		//	return
		//}
		c.ResponseOk(nil)
	}
}
