package models

import "time"

type EndPoint struct {
	Id         int       `json:"id"`
	Name       string    `gorm:"column:endpoint_name;type:varchar(64);not null;default:''" json:"name"`
	Type       string    `gorm:"column:type;size:64;not null;default:''" json:"type"`
	Routers    []Router  `gorm:"" json:"routers"`
	UpdateTime time.Time `gorm:"column:update_time" description:"update_time" json:"update_time"`
	CreatTime  time.Time `gorm:"column:creat_time" description:"creat_time" json:"creat_time"`
}

func (t *EndPoint) TableName() string {
	return "endpoint"
}
