package payload

type OptionCertReq struct {
	Id int `uri:"id" binding:"required"`
}

// CertInfo certs info
type CertInfo struct {
	Csr        string `json:"csr_file"`
	Key        string `json:"key_file"`
	ClientAuth string `json:"client_auth"`
}
