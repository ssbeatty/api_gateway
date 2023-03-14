package payload

type EndPoint struct {
	Id         int      `json:"id"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Routers    []Router `json:"routers"`
	UpdateTime string   `json:"update_time"`
	CreatTime  string   `json:"creat_time"`
}

type Router struct {
	Id          int              `json:"id"`
	EndPointId  int              `json:"endpoint_id"`
	Rule        string           `json:"rule"`
	Type        string           `json:"type"`
	TlsEnable   int              `json:"tls_enable"`
	Priority    int              `json:"priority"`
	Host        string           `json:"host"`
	UpStream    UpStreamInfo     `json:"up_stream"`
	Tls         TlsInfo          `json:"tls"`
	Middlewares []MiddleWareInfo `json:"middlewares"`
}

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
