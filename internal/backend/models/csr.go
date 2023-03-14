package models

type Csr struct {
	Id         int    `json:"id"`
	TlsId      int    `gorm:"column:tls_id" json:"tls_id"`
	FileName   string `gorm:"column:filename;size:128;not null;default:''" json:"file_name"`
	CsrFile    string `gorm:"type:text" json:"csr_file"`
	KeyFile    string `gorm:"type:text" json:"key_file"`
	ClientAuth string `gorm:"column:auth" json:"client_auth"`
}

func (t *Csr) TableName() string {
	return "csr_file"
}

func InsertCsr(tid int, fileName string, keyFile string, csrFile string, clientAuth string) (*Csr, error) {
	csr := Csr{
		TlsId:      tid,
		FileName:   fileName,
		KeyFile:    keyFile,
		CsrFile:    csrFile,
		ClientAuth: clientAuth,
	}
	err := db.Create(&csr).Error
	if err != nil {
		return nil, err
	}
	return &csr, nil
}

func DeleteCsrById(id int) (*Csr, error) {
	ca := Csr{}
	err := db.Where("id = ?", id).First(&ca).Error
	if err != nil {
		return nil, err
	}
	return &ca, nil
}

func UpDataCsr(id int, fileName string, keyFile string, csrFile string, clientAuth string) (*Csr, error) {
	crs := Csr{Id: id}
	err := db.Where("id = ?", id).First(&crs).Error
	if err != nil {
		return nil, err
	}
	if fileName != "" {
		crs.FileName = fileName
	}
	if keyFile != "" {
		crs.KeyFile = keyFile
	}
	if csrFile != "" {
		crs.CsrFile = csrFile
	}
	if clientAuth != "" {
		crs.ClientAuth = clientAuth
	}

	if err := db.Save(&crs).Error; err != nil {
		return nil, err
	}
	return &crs, nil
}

func GetCsrById(id int) (*Csr, error) {
	ca := Csr{}
	err := db.Where("id = ?", id).First(&ca).Error
	if err != nil {
		return nil, err
	}
	return &ca, nil
}
