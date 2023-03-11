package models

type UpStream struct {
	Id          int    `json:"id"`
	Type        string `gorm:"column:router_type;size:64;default:'';not null" json:"type"`
	Path        string `gorm:"column:path" json:"path"`
	Weights     string `gorm:"column:weights" json:"weights"`
	LoadBalance string `gorm:"column:load_balance" json:"load_balance"`
}

func (t *UpStream) TableName() string {
	return "upstream"
}
