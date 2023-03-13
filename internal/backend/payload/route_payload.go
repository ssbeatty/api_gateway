package payload

type MiddleWareInfo struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Config string `json:"config"`
}

type UpStreamInfo struct {
	Type        string `json:"type"`
	Path        string `json:"path"`
	Weights     string `json:"weights"`
	LoadBalance string `json:"load_balance"`
}

type TlsInfo struct {
	Type       string `json:"type"`
	ClientAuth string `json:"client_auth"`
}
