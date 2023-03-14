package models

import (
	"api_gateway/internal/backend/payload"
)

type EndPoint struct {
	Id      int      `json:"id"`
	Name    string   `gorm:"column:endpoint_name;type:varchar(64);not null;default:''" json:"name"`
	Type    string   `gorm:"column:type;size:64;not null;default:''" json:"type"`
	Routers []Router `gorm:"constraint:OnDelete:SET NULL;" json:"routers"`
}

func (t *EndPoint) TableName() string {
	return "endpoint"
}

func InsertEndPoint(name string, tp string, routers []payload.Router) (*EndPoint, error) {

	endPoint := EndPoint{
		Name: name,
		Type: tp,
	}
	err := db.Create(&endPoint).Error
	if err != nil {
		return nil, err
	}

	return &endPoint, nil
}

func DeleteEndPointById(id int) (*EndPoint, error) {
	endPoint := EndPoint{}
	err := db.Where("id = ?", id).First(&endPoint).Error
	if err != nil {
		return nil, err
	}
	if err := db.Model(&endPoint).Association("Routers").Clear(); err != nil {
		logger.Error().AnErr("DeleteEndPointById error Association tag Clear ", err)
	}
	err = db.Delete(&endPoint).Error
	if err != nil {
		logger.Error().AnErr("DeleteEndPointById error", err)

		return nil, err
	}

	return &endPoint, nil
}

func UpDataEndPoint(id int, name string, tp string, routers []payload.Router) (*EndPoint, error) {
	endPoint := EndPoint{Id: id}
	err := db.Where("id = ?", id).First(&endPoint).Error
	if err != nil {
		return nil, err
	}

	if name != "" {
		endPoint.Name = name
	}
	if tp != "" {
		endPoint.Type = tp
	}
	return &endPoint, nil
}

func GetEndPointById(id int) (*EndPoint, error) {
	endPoint := EndPoint{}
	err := db.Where("id = ?", id).First(&endPoint).Error
	if err != nil {
		return nil, err
	}
	return &endPoint, nil
}

func GetEndPointList(pageSize, page int) ([]*EndPoint, error) {
	switch {
	case pageSize <= 0:
		pageSize = 20
	case page <= 0:
		page = 1
	}
	var endPoint []*EndPoint
	offset := (page - 1) * pageSize
	err := db.Order(defaultSort).Offset(offset).Limit(pageSize).Find(&endPoint).Error
	if err != nil {
		return nil, err
	}
	return endPoint, nil
}

func EndPointUpsert(point payload.EndPoint) (*EndPoint, error) {
	var (
		endPoint *EndPoint
		err      error
	)
	if point.Id != 0 {
		endPoint, err = UpDataEndPoint(point.Id, point.Name, point.Type, point.Routers)
		if err != nil {
			return nil, err
		}
	} else {
		endPoint, err = InsertEndPoint(point.Name, point.Type, point.Routers)
		if err != nil {
			return nil, err
		}
	}

	return endPoint, nil
}
