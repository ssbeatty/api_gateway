package payload

type CAInfo struct {
	Csr string `json:"csr"`
	Key string `json:"key"`
}

type OptionCAReq struct {
	Id int `uri:"id" binding:"required"`
}
