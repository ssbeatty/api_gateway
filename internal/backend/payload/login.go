package payload

import "time"

type AdminLoginPasswordReq struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginPasswordResp struct {
	Expire time.Time `json:"expire"`
	Token  string    `json:"token"`
}

type AdminRegisterReq struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
