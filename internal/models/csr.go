package models

type Csr struct {
	Id         int    `json:"id"`
	FileName   string `gorm:"column:filename;size:128;not null;default:''" json:"file_name"`
	CsrFile    string `gorm:"type:text" json:"csr_file"`
	KeyFile    string `gorm:"type:text" json:"key_file"`
	ClientAuth string `gorm:"column:auth" json:"client_auth"`
}

func (t *Csr) TableName() string {
	return "csr_file"
}
