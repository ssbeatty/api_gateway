package models

import (
	"api_gateway/internal/backend/payload"
	"gorm.io/gorm/clause"
)

type CA struct {
	Id         int    `gorm:"primaryKey" json:"id"`
	CertsFile  string `gorm:"type:text" json:"certs_file"`
	KeyFile    string `gorm:"type:text" json:"key_file"`
	ClientAuth string `gorm:"column:auth" json:"client_auth"`
}

func (t *CA) TableName() string {
	return "cas"
}

func InsertCA(CAInfo payload.CAInfo) (*CA, error) {
	ca := CA{
		CertsFile:  CAInfo.CertsFile,
		KeyFile:    CAInfo.KeyFile,
		ClientAuth: CAInfo.ClientAuth,
	}
	err := db.Create(&ca).Error
	if err != nil {
		return nil, err
	}
	return &ca, nil
}

func DeleteCAById(id int) (*CA, error) {
	ca := CA{}
	err := db.Where("id = ?", id).First(&ca).Error
	if err != nil {
		return nil, err
	}
	return &ca, nil
}

func GetCAById(id int) (*CA, error) {
	ca := CA{}
	err := db.Where("id = ?", id).First(&ca).Error
	if err != nil {
		return nil, err
	}
	return &ca, nil
}

func GetAllCAs() ([]CA, error) {
	var cas []CA
	err := db.Preload(clause.Associations).Find(&cas).Error
	if err != nil {
		return nil, err
	}

	return cas, nil
}
