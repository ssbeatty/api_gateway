package models

type MiddleWare struct {
	Id     int    `json:"id"`
	Name   int    `gorm:"column:name;size:64;default:'';not null" json:"name"`
	Type   string `gorm:"column:router_type;size:32;default:'';not null" json:"type"`
	Config string `gorm:"column:path;type:text" json:"path"`
	Route  Router `gorm:"" json:"-"`
}

func (t *MiddleWare) TableName() string {
	return "middle_ware"
}
