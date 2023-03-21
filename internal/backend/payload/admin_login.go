package payload

import "time"

type AdminLoginPasswordReq struct {
	UserName string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type AdminLoginPasswordResp struct {
	Expire time.Time `json:"expire"`
	Token  string    `json:"token"`
}

type AdminRegisterReq struct {
	UserName string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type TenantLoginPasswordReq struct {
	UserName string `json:"username" binding:"required" example:"tenant1"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type TenantRegisterReq struct {
	UserName string `json:"username" binding:"required" example:"tenant1"`
	Password string `json:"password" binding:"required" example:"123456"`
}
