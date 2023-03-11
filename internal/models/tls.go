package models

type TlsConfig struct {
	Id         int    `json:"id"`
	RId        int    `gorm:"column:router_id" json:"r_id"`
	Type       string `gorm:"column:type;type:varchar(32);not null;default:''" json:"type"`
	CsrFile    Csr    `gorm:"" json:"csr_file"`
	KeyFile    string `gorm:"type:text" json:"tls_enable"`
	ClientAuth string `gorm:"client_auth" json:"client_auth"`
	CaFiles    Ca     `gorm:"" json:"ca_files"`
}

func (t *TlsConfig) TableName() string {
	return "tls_config"
}
