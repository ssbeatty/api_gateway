package models

import (
	"api_gateway/internal/backend/payload"
)

type CA struct {
	Id  int    `gorm:"primaryKey" json:"id"`
	Csr string `gorm:"type:text" json:"certs_file"`
	Key string `gorm:"type:text" json:"key_file"`
}

func (t *CA) TableName() string {
	return "ca_certs"
}

func InsertCACerts(CAInfo payload.CAInfo) (*CA, error) {
	ca := CA{
		Csr: CAInfo.Csr,
		Key: CAInfo.Key,
	}
	err := db.Create(&ca).Error
	if err != nil {
		return nil, err
	}
	return &ca, nil
}

func UpdateCACerts(id int, info payload.CAInfo) (*CA, error) {
	ca := &CA{Id: id}

	if info.Csr != "" {
		ca.Csr = info.Csr
	}
	if info.Key != "" {
		ca.Key = info.Key
	}

	if err := db.Updates(&ca).Error; err != nil {
		return nil, err
	}

	return ca, nil
}

func DeleteCACertsById(id int) error {
	ca := CA{}
	err := db.Where("id = ?", id).First(&ca).Error
	if err != nil {
		return err
	}
	err = db.Delete(&ca).Error
	if err != nil {
		return err
	}
	return nil
}

func GetCACertsById(id int) (*CA, error) {
	ca := CA{}
	err := db.Where("id = ?", id).First(&ca).Error
	if err != nil {
		return nil, err
	}
	return &ca, nil
}

func GetAllCACerts() ([]CA, error) {
	var caCerts []CA
	err := db.Find(&caCerts).Error
	if err != nil {
		return nil, err
	}

	return caCerts, nil
}
