package models

import (
	"api_gateway/internal/backend/payload"
	"encoding/json"
	"gorm.io/gorm"
)

type Router struct {
	Id          int      `gorm:"primaryKey" json:"id"`
	Rule        string   `gorm:"column:rule" json:"rule"`
	Type        string   `gorm:"column:router_type;not null" json:"type"`
	TlsEnable   bool     `gorm:"column:tls_enable" json:"tls_enable"`
	Priority    int      `gorm:"column:priority" json:"priority"`
	Host        string   `gorm:"column:host" json:"host"`
	UpStream    string   `gorm:"column:upstream" json:"upstream"`
	Middlewares string   `gorm:"column:middlewares" json:"middlewares"`
	EndpointId  int      `gorm:"column:endpoint_id" json:"endpoint_id"`
	CertID      int      `gorm:"column:cert_id" json:"cert_id"`
	Cert        Cert     `gorm:"constraint:OnDelete:SET NULL;" json:"-"`
	CaID        int      `gorm:"column:ca_id" json:"ca_id"`
	CA          CA       `gorm:"constraint:OnDelete:SET NULL;" json:"-"`
	Endpoint    Endpoint `json:"-"`
}

func (t *Router) TableName() string {
	return "routers"
}

func InsertRouter(endpointID int, info payload.RouterInfo) (*Router, error) {
	session := db.Session(&gorm.Session{})

	upstreamJson, err := json.Marshal(info.UpStream)
	if err != nil {
		return nil, err
	}
	middlewareJson, err := json.Marshal(info.Middlewares)
	if err != nil {
		return nil, err
	}
	r := Router{
		Id:          info.Id,
		EndpointId:  endpointID,
		Rule:        info.Rule,
		Type:        string(info.Type),
		TlsEnable:   info.TlsEnable,
		Priority:    info.Priority,
		Host:        info.Host,
		UpStream:    string(upstreamJson),
		Middlewares: string(middlewareJson),
	}

	if info.CertId != 0 {
		if cert, errF := GetCertById(info.CertId); err == nil {
			r.CertID = cert.Id
		} else {
			return nil, errF
		}
	} else {
		session = session.Omit("CertID")
	}

	if err = session.Save(&r).Error; err != nil {
		return nil, err
	}

	return &r, nil
}

func ClearEndpointRouters(endpointID int) error {
	return db.Where("endpoint_id = ?", endpointID).Delete(&Router{}).Error
}
