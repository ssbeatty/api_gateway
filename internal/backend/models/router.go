package models

type Router struct {
	Id          int          `json:"id"`
	EndPointId  int          `gorm:"column:endpoint_id" json:"endpoint_id"`
	Rule        string       `gorm:"column:rule" json:"rule"`
	Type        string       `gorm:"column:router_type;size:64;not null;default:''" json:"type"`
	TlsEnable   int          `gorm:"column:tls_enable" json:"tls_enable"`
	Priority    int          `gorm:"column:priority" json:"priority"`
	Host        string       `gorm:"column:host" json:"host"`
	UpStream    UpStream     `gorm:"constraint:OnDelete:SET NULL;" json:"up_stream"`
	Tls         TlsConfig    `gorm:"constraint:OnDelete:SET NULL;" json:"tls"`
	Middlewares []MiddleWare `gorm:"constraint:OnDelete:SET NULL;" json:"middlewares"`
}

func (t *Router) TableName() string {
	return "router"
}

func UpDataRouter(id int, rule string, tp string, tlsEnable int, priority int, host string) (*Router, error) {
	router := Router{Id: id}
	err := db.Where("id = ?", id).First(&router).Error
	if err != nil {
		return nil, err
	}
	if rule != "" {
		router.Rule = rule
	}
	if tp != "" {
		router.Type = tp
	}

	if tlsEnable != router.TlsEnable {
		router.TlsEnable = tlsEnable
	}
	if priority != router.Priority {
		router.Priority = priority
	}
	if host != "" {
		router.Host = host
	}

	if err := db.Save(&router).Error; err != nil {
		return nil, err
	}
	return &router, nil
}

func GetRouterById(id int) (*Router, error) {
	router := Router{}
	err := db.Where("id = ?", id).First(&router).Error
	if err != nil {
		return nil, err
	}
	return &router, nil
}

func DeleteRouterById(id int) (*Router, error) {
	router := Router{}
	err := db.Where("id = ?", id).First(&router).Error
	if err != nil {
		return nil, err
	}
	if err := db.Model(&router).Association("Middlewares").Clear(); err != nil {
		return nil, err
	}
	if err := db.Model(&router).Association("UpStream").Clear(); err != nil {
		return nil, err
	}
	if err := db.Model(&router).Association("Tls").Clear(); err != nil {
		return nil, err
	}
	err = db.Delete(&router).Error
	if err != nil {
		return nil, err
	}

	return &router, nil
}

func InsertRouter(eid int, role string, tp string, tlsEnable int, priority int, host string, UpStreamID int, tlsID int, middlewares []int) (*Router, error) {
	var middleWares []MiddleWare
	for _, rid := range middlewares {
		m := MiddleWare{}
		err := db.Where("id = ?", rid).First(&m).Error
		if err != nil {
			logger.Error().AnErr("InsertRouter error when First middleware", err)
			continue
		}
		middleWares = append(middleWares, m)
	}
	r := Router{
		EndPointId: eid,
		Rule:       role,
		Type:       tp,
		TlsEnable:  tlsEnable,
		Priority:   priority,
		Host:       host,
	}

	return &r, nil
}
