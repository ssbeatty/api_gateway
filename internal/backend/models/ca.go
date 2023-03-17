package models

type Ca struct {
	Id       int    `json:"id"`
	TlsId    int    `gorm:"column:tls_id" json:"tls_id"`
	FileName string `gorm:"column:filename" json:"file_name"`
	KeyFile  string `gorm:"column:key_file;type:text" json:"key_file"`
}

func (t *Ca) TableName() string {
	return "cas"
}

func InsertCa(tid int, fileName string, keyFile string) (*Ca, error) {
	ca := Ca{
		TlsId:    tid,
		FileName: fileName,
		KeyFile:  keyFile,
	}
	err := db.Create(&ca).Error
	if err != nil {
		return nil, err
	}
	return &ca, nil
}

func DeleteCaById(id int) (*Ca, error) {
	ca := Ca{}
	err := db.Where("id = ?", id).First(&ca).Error
	if err != nil {
		return nil, err
	}
	return &ca, nil
}

func UpDataCa(id int, fileName string, keyFile string) (*Ca, error) {
	ca := Ca{Id: id}
	err := db.Where("id = ?", id).First(&ca).Error
	if err != nil {
		return nil, err
	}
	if fileName != "" {
		ca.FileName = fileName
	}
	if keyFile != "" {
		ca.KeyFile = keyFile
	}

	if err := db.Save(&ca).Error; err != nil {
		return nil, err
	}
	return &ca, nil
}

func GetCaById(id int) (*Ca, error) {
	ca := Ca{}
	err := db.Where("id = ?", id).First(&ca).Error
	if err != nil {
		return nil, err
	}
	return &ca, nil
}
