package models

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Admin struct {
	Id         int       `json:"id"`
	Username   string    `gorm:"uniqueIndex:a_u_username_unique;column:username;size:128;not null;default:''" json:"username"` // 用户名
	Password   string    `gorm:"column:password;size:255;not null;default:''" json:"password"`                                 // 密码
	HeadImg    string    `gorm:"column:head_img" json:"head_img"`                                                              //头像
	UpdateTime time.Time `gorm:"column:update_time" description:"update_time" json:"update_time"`
	CreatTime  time.Time `gorm:"column:creat_time" description:"creat_time" json:"creat_time"`
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
		Username:   name,
		Password:   pass,
		UpdateTime: time.Now(),
		CreatTime:  time.Now(),
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
	admin.UpdateTime = time.Now()
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

// ComparePasswords 比对用户密码是否正确
func ComparePasswords(dbPassword, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)); err != nil {
		return false
	}
	return true
}
