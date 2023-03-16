package handler

import (
	"api_gateway/internal/backend/models"
	"api_gateway/internal/backend/payload"
)

// EndpointsQuery query all endpoints
func (s *Service) EndpointsQuery(c *Context) {
	records, err := models.QueryEndpoints()
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

		c.ResponseOk(record)
	}
}
