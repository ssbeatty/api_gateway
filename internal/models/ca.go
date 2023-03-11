package models

type Ca struct {
	Id       int    `json:"id"`
	FileName string `gorm:"column:filename;size:128;not null;default:''" json:"file_name"`
	KeyFile  string `gorm:"column:key_file;type:text" json:"key_file"`
}

func (t *Ca) TableName() string {
	return "ca_file"
}
