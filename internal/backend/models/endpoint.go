package models

import (
	"api_gateway/internal/backend/payload"
	"gorm.io/gorm/clause"
)

type Endpoint struct {
	Id         int      `json:"id"`
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
		_, err = InsertOrUpdateRouter(endpoint.Id, info)
		if err != nil {
			logger.Error().Err(err).Msg("Error when insert router")
			continue
		}
	}

	return &endpoint, nil
}

func QueryEndpoints() ([]Endpoint, error) {
	var records []Endpoint
	err := db.Preload(clause.Associations).Find(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

func DeleteEndPointById(id int) (*Endpoint, error) {
	endPoint := Endpoint{}
	err := db.Where("id = ?", id).First(&endPoint).Error
	if err != nil {
		return nil, err
	}

	err = db.Delete(&endPoint).Error
	if err != nil {
		logger.Error().AnErr("DeleteEndPointById error", err)

		return nil, err
	}

	return &endPoint, nil
}

func GetEndPointById(id int) (*Endpoint, error) {
	endPoint := Endpoint{}
	err := db.Where("id = ?", id).First(&endPoint).Error
	if err != nil {
		return nil, err
	}
	return &endPoint, nil
}
