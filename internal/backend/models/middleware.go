package models

type MiddleWare struct {
	Id     int    `json:"id"`
	RId    int    `gorm:"column:router_id" json:"r_id"`
	Name   string `gorm:"column:name;size:64;default:'';not null" json:"name"`
	Type   string `gorm:"column:router_type;size:32;default:'';not null" json:"type"`
	Config string `gorm:"column:path;type:text" json:"config"`
}

func (t *MiddleWare) TableName() string {
	return "middle_ware"
}

func InsertMiddleWare(rid int, name string, tp string, config string) (*MiddleWare, error) {
	middleWare := MiddleWare{
		RId:    rid,
		Name:   name,
		Type:   tp,
		Config: config,
	}
	err := db.Create(&middleWare).Error
	if err != nil {
		return nil, err
	}
	return &middleWare, nil
}

func DeleteMiddleWareById(id int) (*MiddleWare, error) {
	middleWare := MiddleWare{}
	err := db.Where("id = ?", id).First(&middleWare).Error
	if err != nil {
		return nil, err
	}
	return &middleWare, nil
}

func UpDataMiddleWare(id int, name string, tp string, config string) (*MiddleWare, error) {
	middleWare := MiddleWare{Id: id}
	err := db.Where("id = ?", id).First(&middleWare).Error
	if err != nil {
		return nil, err
	}
	if name != "" {
		middleWare.Name = name
	}
	if tp != "" {
		middleWare.Type = tp
	}

	if config != "" {
		middleWare.Config = config
	}

	if err := db.Save(&middleWare).Error; err != nil {
		return nil, err
	}
	return &middleWare, nil
}

func GetMiddleWareById(id int) (*MiddleWare, error) {
	middleWare := MiddleWare{}
	err := db.Where("id = ?", id).First(&middleWare).Error
	if err != nil {
		return nil, err
	}
	return &middleWare, nil
}
