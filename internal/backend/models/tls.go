package models

type TlsConfig struct {
	Id         int    `json:"id"`
	RId        int    `gorm:"column:router_id" json:"r_id"`
	Type       string `gorm:"column:type;type:varchar(32);not null;default:''" json:"type"`
	CsrFile    Csr    `gorm:"constraint:OnDelete:SET NULL;" json:"csr_file"`
	KeyFile    string `gorm:"type:text" json:"tls_enable"`
	ClientAuth string `gorm:"client_auth" json:"client_auth"`
	CaFiles    []Ca   `gorm:"constraint:OnDelete:SET NULL;" json:"ca_files"`
}

func (t *TlsConfig) TableName() string {
	return "tls_config"
}

func InsertTlsConfig(rid int, tp string, keyFile string, clientAuth string) (*TlsConfig, error) {
	tlsConfig := TlsConfig{
		RId:        rid,
		Type:       tp,
		KeyFile:    keyFile,
		ClientAuth: clientAuth,
	}
	err := db.Create(&tlsConfig).Error
	if err != nil {
		return nil, err
	}
	return &tlsConfig, nil
}

func DeleteTlsConfigById(id int) (*TlsConfig, error) {
	tlsConfig := TlsConfig{}
	err := db.Where("id = ?", id).First(&tlsConfig).Error
	if err != nil {
		return nil, err
	}
	if err := db.Model(&tlsConfig).Association("CsrFile").Clear(); err != nil {
		logger.Error().AnErr("DeleteTlsConfigById error Association tag CsrFile clear", err)
	}
	if err := db.Model(&tlsConfig).Association("CaFiles").Clear(); err != nil {
		logger.Error().AnErr("DeleteTlsConfigById error Association CaFiles Clear", err)
	}
	err = db.Delete(&tlsConfig).Error
	if err != nil {
		return nil, err
	}

	return &tlsConfig, nil
}

func UpDataTlsConfig(id int, tp string, keyFile string, clientAuth string) (*TlsConfig, error) {
	tlsConfig := TlsConfig{Id: id}
	err := db.Where("id = ?", id).First(&tlsConfig).Error
	if err != nil {
		return nil, err
	}
	if tp != "" {
		tlsConfig.Type = tp
	}
	if keyFile != "" {
		tlsConfig.KeyFile = keyFile
	}
	if tp != "" {
		tlsConfig.ClientAuth = clientAuth
	}

	err = db.Save(&tlsConfig).Error
	if err != nil {
		return nil, err
	}
	return &tlsConfig, nil
}

func GetTlsConfigById(id int) (*TlsConfig, error) {
	tlsConfig := TlsConfig{}
	err := db.Where("id = ?", id).First(&tlsConfig).Error
	if err != nil {
		return nil, err
	}
	return &tlsConfig, nil
}
