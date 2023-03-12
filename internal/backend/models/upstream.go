package models

type UpStream struct {
	Id          int    `json:"id"`
	RId         int    `gorm:"column:router_id" json:"r_id"`
	Type        string `gorm:"column:router_type;size:64;default:'';not null" json:"type"`
	Path        string `gorm:"column:path" json:"path"`
	Weights     string `gorm:"column:weights" json:"weights"`
	LoadBalance string `gorm:"column:load_balance" json:"load_balance"`
}

func (t *UpStream) TableName() string {
	return "upstream"
}

func InsertUpStream(rid int, path string, tp string, weights string, loadBalance string) (*UpStream, error) {
	upStream := UpStream{
		RId:         rid,
		Path:        path,
		Type:        tp,
		Weights:     weights,
		LoadBalance: loadBalance,
	}
	err := db.Create(&upStream).Error
	if err != nil {
		return nil, err
	}
	return &upStream, nil
}

func DeleteUpStreamById(id int) (*UpStream, error) {
	upStream := UpStream{}
	err := db.Where("id = ?", id).First(&upStream).Error
	if err != nil {
		return nil, err
	}
	return &upStream, nil
}

func UpDataUpStream(id int, path string, tp string, weights string, loadBalance string) (*UpStream, error) {
	upStream := UpStream{Id: id}
	err := db.Where("id = ?", id).First(&upStream).Error
	if err != nil {
		return nil, err
	}
	if path != "" {
		upStream.Path = path
	}
	if tp != "" {
		upStream.Type = tp
	}

	if weights != "" {
		upStream.Weights = weights
	}

	if loadBalance != "" {
		upStream.LoadBalance = loadBalance
	}

	if err := db.Save(&upStream).Error; err != nil {
		return nil, err
	}
	return &upStream, nil
}

func GetUpStreamById(id int) (*UpStream, error) {
	upStream := UpStream{}
	err := db.Where("id = ?", id).First(&upStream).Error
	if err != nil {
		return nil, err
	}
	return &upStream, nil
}
