package models

type Cert struct {
	Id         int    `json:"id"`
	CsrFile    string `gorm:"type:text" json:"csr_file"`
	KeyFile    string `gorm:"type:text" json:"key_file"`
	ClientAuth string `gorm:"column:auth" json:"client_auth"`
}

func (t *Cert) TableName() string {
	return "certs"
}

func GetCertById(id int) (*Cert, error) {
	record := Cert{}
	err := db.Where("id = ?", id).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}
