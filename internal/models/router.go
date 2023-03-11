package models

type Router struct {
	Id          int          `json:"id"`
	EId         int          `gorm:"column:endpoint_id;default:0;not null" json:"eid"`
	Type        string       `gorm:"column:router_type;size:64;not null;default:''" json:"type"`
	TlsEnable   int          `gorm:"column:tls_enable" json:"tls_enable"`
	Priority    int          `gorm:"column:priority" json:"priority"`
	Host        string       `gorm:"column:host" json:"host"`
	UpStream    UpStream     `gorm:"" json:"up_stream"`
	Tls         TlsConfig    `gorm:"" json:"tls"`
	Middlewares []MiddleWare `gorm:"" json:"middlewares"`
}

func (t *Router) TableName() string {
	return "router"
}
