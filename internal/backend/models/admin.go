package models

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	Id       int    `json:"id"`
	Username string `gorm:"uniqueIndex:a_u_username_unique;column:username;size:128;not null;default:''" json:"username"`
	Password string `gorm:"column:password;size:255;not null;default:''" json:"-"`
	Avatar   string `gorm:"column:avatar" json:"avatar"`
}

func (t *Admin) TableName() string {
	return "admin"
}

func InsertAdmin(name string, ps string) (*Admin, error) {
	pass, err := PasswordHash(ps)
	if err != nil {
		return nil, err
	}
	admin := Admin{
		Username: name,
		Password: pass,
	}
	err = db.Create(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func UpdateAdmin(id int, Username string, Password string) (*Admin, error) {
	admin := Admin{Id: id}
	err := db.Where("id = ?", id).First(&admin).Error
	if err != nil {
		return nil, err
	}
	if Username != "" {
		admin.Username = Username
	}
	if Password != "" {
		admin.Password, err = PasswordHash(Password)
		if err != nil {
			log.Error().AnErr("update admin info failed", err)
			return nil, err
		}
	}
	if err := db.Save(&admin).Error; err != nil {
		log.Error().AnErr("update admin info failed", err)
	}
	return &admin, nil

}

func GetAdminById(id int) (*Admin, error) {
	admin := Admin{}
	err := db.Where("id = ?", id).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func GetAdminByUserName(username string) (*Admin, error) {
	admin := Admin{}
	err := db.Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func DeleteAdminById(id int) (*Admin, error) {
	admin := Admin{}
	err := db.Delete(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// PasswordHash 密码hash
func PasswordHash(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash), err
}
