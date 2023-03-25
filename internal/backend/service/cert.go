package service

import (
	"api_gateway/internal/backend/models"
	"api_gateway/internal/backend/payload"
)

// CertsQuery query all certs
func (s *Service) CertsQuery(c *Context) {
	records, err := models.GetAllCerts()
	if err != nil {
		c.ResponseError(err.Error())
		return
	}
	c.ResponseOk(records)
}

// CertsCreate create certs
func (s *Service) CertsCreate(c *Context) {
	var req payload.CertInfo
	err := c.ShouldBind(&req)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		record, err := models.InsertCerts(req)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(record)
	}
}

// CertsUpdate update certs
func (s *Service) CertsUpdate(c *Context) {
	var uri payload.OptionCertReq
	err := c.ShouldBindUri(&uri)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		var req payload.CertInfo
		err := c.ShouldBind(&req)
		if err != nil {
			c.ResponseError(err.Error())
		} else {
			record, err := models.UpdateCerts(uri.Id, req)
			if err != nil {
				c.ResponseError(err.Error())
				return
			}
			c.ResponseOk(record)
		}
	}
}

// CertsDetail  get certs detail
func (s *Service) CertsDetail(c *Context) {
	var uri payload.OptionCertReq
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

// CertsDelete delete certs
func (s *Service) CertsDelete(c *Context) {
	var uri payload.OptionCertReq
	err := c.ShouldBindUri(&uri)
	if err != nil {
		c.ResponseError(err.Error())
	} else {
		err := models.DeleteCertsById(uri.Id)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(nil)
	}
}
