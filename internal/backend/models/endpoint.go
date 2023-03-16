package models

import (
	"api_gateway/internal/backend/payload"
	"gorm.io/gorm/clause"
)

type Endpoint struct {
	Id         int      `gorm:"primaryKey" json:"id"`
	Name       string   `gorm:"column:endpoint_name;not null;index:type_port,unique" json:"name"`
	Type       string   `gorm:"column:type;not null" json:"type"`
	ListenPort int      `gorm:"column:listen_port;index:type_port,unique" json:"listen_port"`
	Routers    []Router `gorm:"constraint:OnDelete:CASCADE;" json:"routers"`
}

func (t *Endpoint) TableName() string {
	return "endpoints"
}

func InsertEndpoint(info payload.PostEndPointReq) (*Endpoint, error) {
	endpoint := Endpoint{
		Name:       info.Name,
		Type:       string(info.Type),
		ListenPort: info.ListenPort,
	}

	err := db.Create(&endpoint).Error
	if err != nil {
		return nil, err
	}

	for _, info := range info.Routers {
		_, err = InsertRouter(endpoint.Id, info)
		if err != nil {
			logger.Error().Err(err).Msg("Error when insert router")
			continue
		}
	}

	return &endpoint, nil
}

func GetAllEndpoints() ([]Endpoint, error) {
	var records []Endpoint
	err := db.Preload(clause.Associations).Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

func UpdateEndpoint(id int, info payload.PostEndPointReq) (*Endpoint, error) {
	endpoint := Endpoint{Id: id}

	if info.Name != "" {
		endpoint.Name = info.Name
	}
	if info.Type != "" {
		endpoint.Type = string(info.Type)
	}
	if info.ListenPort != 0 {
		endpoint.ListenPort = info.ListenPort
	}

	if err := db.Updates(&endpoint).Error; err != nil {
		return nil, err
	}

	err := ClearEndpointRouters(id)
	if err != nil {
		return nil, err
	}

	for _, i := range info.Routers {
		_, err = InsertRouter(endpoint.Id, i)
		if err != nil {
			logger.Error().Err(err).Msg("Error when insert router")
			continue
		}
	}
	return GetEndPointById(id)
}

func DeleteEndPointById(id int) error {
	endPoint := Endpoint{}
	err := db.Where("id = ?", id).First(&endPoint).Error
	if err != nil {
		return err
	}

	err = db.Delete(&endPoint).Error
	if err != nil {
		return err
	}

	return nil
}

func GetEndPointById(id int) (*Endpoint, error) {
	endPoint := Endpoint{}
	err := db.Where("id = ?", id).Preload(clause.Associations).First(&endPoint).Error
	if err != nil {
		return nil, err
	}
	return &endPoint, nil
}
