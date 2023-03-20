package models

// Tenant 租户
type Tenant struct {
	Id              int    `json:"id"`
	Username        string `gorm:"uniqueIndex:a_u_username_unique;column:username;size:128;not null" json:"username"` // 用户名
	Password        string `gorm:"column:password;size:255;not null" json:"password"`                                 // 密码
	RequestQuantity string `gorm:"column:create_at" json:"request_quantity"`
	Token           string `gorm:"column:token" json:"token"`
}

func (t *Tenant) TableName() string {
	return "tenants"
}

func InsertTenant(name string, ps string) (*Tenant, error) {
	pass, err := PasswordHash(ps)
	if err != nil {
		return nil, err
	}
	tenant := Tenant{
		Username: name,
		Password: pass,
	}
	err = db.Create(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func GetTenantById(id int) (*Tenant, error) {
	tenant := Tenant{}
	err := db.Where("id = ?", id).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func GetAllTenant() ([]*Tenant, error) {
	var tenants []*Tenant
	err := db.Find(&tenants).Error
	if err != nil {
		return nil, err
	}

	return tenants, nil
}

func DeleteTenantById(id int) (*Tenant, error) {
	tenant := Tenant{}
	err := db.Delete(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}
