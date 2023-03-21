package service

import (
	"api_gateway/internal/backend/models"
	"api_gateway/internal/backend/payload"
)

// CAsQuery query all CA
func (s *Service) CAsQuery(c *Context) {
	records, err := models.GetAllCAs()
	if err != nil {
		c.ResponseError(err.Error())
		return
	}
	c.ResponseOk(records)
}

// CAsCreate create CA
func (s *Service) CAsCreate(c *Context) {
	var req payload.CAInfo
	err := c.ShouldBind(&req)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		record, err := models.InsertCA(req)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(record)
	}
}

// CAsDetail  get CA detail
func (s *Service) CAsDetail(c *Context) {
	var uri payload.OptionCAReq
	err := c.ShouldBindUri(&uri)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		record, err := models.GetCAById(uri.Id)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(record)
	}
}

// CAsDelete delete CA
func (s *Service) CAsDelete(c *Context) {
	var uri payload.OptionCAReq
	err := c.ShouldBindUri(&uri)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		_, err := models.DeleteCAById(uri.Id)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(nil)
	}
}
