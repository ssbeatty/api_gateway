package payload

type AdminLoginPasswordReq struct {
	UserName string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type AdminLoginPasswordResp struct {
	UserName      string `json:"username"`
	Toke          string `json:"token"`
	TokenExpireAt int32  `json:"token_expire_at"`
}

type AdminRegisterReq struct {
	UserName string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type OauthLogoutReq struct {
	UserId int `form:"user_id"`
}
