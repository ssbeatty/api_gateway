package payload

type AdminLoginPasswordReq struct {
	UserName string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"123456"`
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
