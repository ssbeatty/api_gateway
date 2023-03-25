package models

import "api_gateway/internal/backend/payload"

type Cert struct {
	Id         int    `json:"id"`
	Csr        string `gorm:"type:text" json:"csr_file"`
	Key        string `gorm:"type:text" json:"key_file"`
	ClientAuth string `gorm:"column:auth" json:"client_auth"`
}

func (t *Cert) TableName() string {
	return "certs"
}

// GetAllCerts get all certs
func GetAllCerts() ([]Cert, error) {
	var records []Cert
	err := db.Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// GetCertById get cert by id
func GetCertById(id int) (*Cert, error) {
	record := Cert{}
	err := db.Where("id = ?", id).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// InsertCerts insert cert
func InsertCerts(cert payload.CertInfo) (*Cert, error) {
	certs := Cert{
		Csr:        cert.Csr,
		Key:        cert.Key,
		ClientAuth: cert.ClientAuth,
	}
	err := db.Create(&certs).Error
	if err != nil {
		return nil, err
	}
	return &certs, nil
}

// UpdateCerts update cert
func UpdateCerts(id int, info payload.CertInfo) (*Cert, error) {
	cert := &Cert{Id: id}

	if info.Csr != "" {
		cert.Csr = info.Csr
	}
	if info.Key != "" {
		cert.Key = info.Key
	}
	if info.ClientAuth != "" {
		cert.ClientAuth = info.ClientAuth
	}

	if err := db.Updates(&cert).Error; err != nil {
		return nil, err
	}

	return cert, nil
}

// DeleteCertsById delete cert by id
func DeleteCertsById(id int) error {
	cert := Cert{}
	err := db.Where("id = ?", id).First(&cert).Error
	if err != nil {
		return err
	}
	err = db.Delete(&cert).Error
	if err != nil {
		return err
	}
	return nil
}
