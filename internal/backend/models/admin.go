package models

import (
	"api_gateway/internal/backend/payload"
	"api_gateway/internal/backend/utils"
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

func OauthLoginPassword(req payload.OauthLoginPasswordReq) (payload.OauthSuccessData, error) {
	var OauthSuccessData payload.OauthSuccessData
	admin, err := GetAdminById(req.Id)
	if err != nil {
		return OauthSuccessData, err
	}
	hash, err := PasswordHash(req.Password)
	if err != nil {
		return OauthSuccessData, err
	}
	if hash != admin.Password {
		return OauthSuccessData, err
	}
	return SetLoginJwtToken(req.Id, admin.Username)
}

func SetLoginJwtToken(userId int, userName string) (payload.OauthSuccessData, error) {
	var OauthSuccessData payload.OauthSuccessData
	token, exp, err := utils.GenerateToken(userId, userName)
	if err != nil {
		return OauthSuccessData, err
	}
	OauthSuccessData.Id = userId
	OauthSuccessData.UserName = userName
	OauthSuccessData.Toke = token
	OauthSuccessData.TokenExpireAt = utils.TimeStandardFormat(time.Unix(exp, 0), false)
	return OauthSuccessData, nil
}
