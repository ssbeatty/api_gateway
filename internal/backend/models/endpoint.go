package models

import (
	"api_gateway/internal/backend/utils"
	"time"
)

type EndPoint struct {
	Id         int      `json:"id"`
	Name       string   `gorm:"column:endpoint_name;type:varchar(64);not null;default:''" json:"name"`
	Type       string   `gorm:"column:type;size:64;not null;default:''" json:"type"`
	Routers    []Router `gorm:"foreignKey:EndPointID" json:"routers"`
	UpdateTime string   `gorm:"column:update_time" description:"update_time" json:"update_time"`
	CreatTime  string   `gorm:"column:creat_time" description:"creat_time" json:"creat_time"`
}

func (t *EndPoint) TableName() string {
	return "endpoint"
}

func InsertEndPoint(name string, tp string, routers []int) (*EndPoint, error) {
	var rs []Router
	for _, tagId := range rs {
		router := Router{}
		err := db.Where("id = ?", tagId).First(&router).Error
		if err != nil {
			logger.Error().AnErr("InsertEndPoint error when First router ", err)
			continue
		}
		rs = append(rs, router)
	}
	endPoint := EndPoint{
		Name:       name,
		Type:       tp,
		Routers:    rs,
		CreatTime:  time.Now().Format(utils.StandardFormat),
		UpdateTime: time.Now().Format(utils.StandardFormat),
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

func UpDataEndPoint(id int, name string, tp string, routers []int) (*EndPoint, error) {
	endPoint := EndPoint{Id: id}
	err := db.Where("id = ?", id).First(&endPoint).Error
	if err != nil {
		return nil, err
	}

	if len(routers) > 0 {
		var rs []Router
		for _, routerId := range routers {
			r := Router{}
			err = db.Where("id = ?", routerId).First(&r).Error
			if err != nil {
				logger.Error().AnErr("upDataEndpoint error when first router", err)
				continue
			}
			rs = append(rs, r)
		}
		if len(rs) != 0 {
			if err := db.Model(&endPoint).Association("Routers").Clear(); err != nil {
				logger.Error().AnErr("upDataEndpoint error when Association router Clear", err)

			}
			endPoint.Routers = rs
		}
	} else {
		if err := db.Model(&endPoint).Association("Routers").Clear(); err != nil {
			logger.Error().AnErr("upDataEndpoint Association router Clear failed", err)
		}
	}

	if name != "" {
		endPoint.Name = name
	}
	if tp != "" {
		endPoint.Type = tp
	}
	endPoint.UpdateTime = time.Now().Format(utils.StandardFormat)
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
