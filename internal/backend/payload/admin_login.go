package payload

type OauthLoginPasswordReq struct {
	Id       int    `form:"id"`
	UserName string `form:"user_name" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type OauthLoginPasswordResp struct {
	code    int         `form:"code"`
	message string      `form:"message"`
	data    interface{} `form:"data"`
}

type OauthLogoutReq struct {
	UserId int `form:"user_id"`
}

type OauthSuccessData struct {
	Id            int    `form:"id"`
	UserName      string `form:"user_name"`
	Toke          string `form:"token"`
	TokenExpireAt string `form:"token_expire_at"`
}

type RegisterUser struct {
	UserName string `form:"user_name" binding:"required"`
	Password string `form:"password" binding:"required"`
}
