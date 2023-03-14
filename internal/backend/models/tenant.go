package models

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm/clause"
)

// Tenant 租户
type Tenant struct {
	Id              int    `json:"id"`
	Username        string `gorm:"uniqueIndex:a_u_username_unique;column:username;size:128;not null;default:''" json:"username"` // 用户名
	Password        string `gorm:"column:password;size:255;not null;default:''" json:"password"`                                 // 密码
	RequestQuantity string `gorm:"column:create_at" json:"request_quantity"`
	Token           string `gorm:"column:token" json:"token"`
}

func (t *Tenant) TableName() string {
	return "tenant"
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

func UpdateTenant(id int, Username string, Password string, token, RequestQuantity string) (*Tenant, error) {
	tenant := Tenant{Id: id}
	err := db.Where("id = ?", id).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	if Username != "" {
		tenant.Username = Username
	}
	if Password != "" {
		tenant.Password, err = PasswordHash(Password)
		if err != nil {
			log.Error().AnErr("update admin info failed", err)
			return nil, err
		}
	}

	if token != "" {
		tenant.Token = token
	}
	if RequestQuantity != "" {
		tenant.RequestQuantity = RequestQuantity
	}

	if err := db.Save(&tenant).Error; err != nil {
		log.Error().AnErr("update admin info failed", err)
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
	err := db.Preload(clause.Associations).Find(&tenants).Error
	if err != nil {
		return nil, err
	}

	return tenants, nil
}

func DeleteTenantById(id int) (*Admin, error) {
	admin := Admin{}
	err := db.Delete(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}
