package service

import (
	"api_gateway/internal/backend/models"
	"api_gateway/internal/backend/payload"
)

// EndpointsQuery query all endpoints
func (s *Service) EndpointsQuery(c *Context) {
	records, err := models.GetAllEndpoints()
	if err != nil {
		c.ResponseError(err.Error())
		return
	}
	c.ResponseOk(records)
}

// EndpointsCreate create endpoint
func (s *Service) EndpointsCreate(c *Context) {
	var req payload.PostEndPointReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		record, err := models.InsertEndpoint(req)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		s.ReloadAllEndpoint()
		c.ResponseOk(record)
	}
}

// EndpointsUpdate update endpoint
func (s *Service) EndpointsUpdate(c *Context) {
	var uri payload.OptionEndpointReq
	err := c.ShouldBindUri(&uri)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		var req payload.PostEndPointReq
		err = c.ShouldBindJSON(&req)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		resp, err := models.UpdateEndpoint(uri.Id, req)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		s.ReloadAllEndpoint()
		c.ResponseOk(resp)
	}
}

// EndpointsDetail get endpoint detail
func (s *Service) EndpointsDetail(c *Context) {
	var uri payload.OptionEndpointReq
	err := c.ShouldBindUri(&uri)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		record, err := models.GetEndPointById(uri.Id)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(record)
	}
}

// EndpointsDelete delete endpoint
func (s *Service) EndpointsDelete(c *Context) {
	var uri payload.OptionEndpointReq
	err := c.ShouldBindUri(&uri)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		err = models.DeleteEndPointById(uri.Id)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		s.ReloadAllEndpoint()
		c.ResponseOk(nil)
	}
}
