package service

import (
	"api_gateway/internal/backend/models"
	"api_gateway/internal/backend/payload"
)

// CACertsQuery query all CA
func (s *Service) CACertsQuery(c *Context) {
	records, err := models.GetAllCACerts()
	if err != nil {
		c.ResponseError(err.Error())
		return
	}
	c.ResponseOk(records)
}

// CACertsCreate create CA
func (s *Service) CACertsCreate(c *Context) {
	var req payload.CAInfo
	err := c.ShouldBind(&req)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		record, err := models.InsertCACerts(req)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(record)
	}
}

// CACertsUpdate update CA
func (s *Service) CACertsUpdate(c *Context) {
	var uri payload.OptionCAReq
	err := c.ShouldBindUri(&uri)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		var req payload.CAInfo
		err := c.ShouldBind(&req)
		if err != nil {
			c.ResponseError(err.Error())
		} else {
			record, err := models.UpdateCACerts(uri.Id, req)
			if err != nil {
				c.ResponseError(err.Error())
				return
			}
			c.ResponseOk(record)
		}
	}
}

// CACertsDetail  get CA detail
func (s *Service) CACertsDetail(c *Context) {
	var uri payload.OptionCAReq
	err := c.ShouldBindUri(&uri)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		record, err := models.GetCACertsById(uri.Id)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(record)
	}
}

// CACertsDelete delete CA
func (s *Service) CACertsDelete(c *Context) {
	var uri payload.OptionCAReq
	err := c.ShouldBindUri(&uri)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		err := models.DeleteCACertsById(uri.Id)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(nil)
	}
}
