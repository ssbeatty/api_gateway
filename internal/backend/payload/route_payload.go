package payload

type EndPoint struct {
	Id      int      `form:"id"`
	Name    string   `form:"name"`
	Type    string   `form:"type"`
	Routers []Router `form:"routers"`
}

type Router struct {
	Id          int              `form:"id"`
	EndPointId  int              `form:"endpoint_id"`
	Rule        string           `form:"rule"`
	Type        string           `form:"type"`
	TlsEnable   int              `form:"tls_enable"`
	Priority    int              `form:"priority"`
	Host        string           `form:"host"`
	UpStream    UpStreamInfo     `form:"up_stream"`
	Tls         TlsInfo          `form:"tls"`
	Middlewares []MiddleWareInfo `form:"middlewares"`
}

type MiddleWareInfo struct {
	Id     int    `form:"id"`
	Name   string `form:"name"`
	Type   string `form:"type"`
	Config string `form:"config"`
}

type UpStreamInfo struct {
	Type        string `form:"type"`
	Path        string `form:"path"`
	Weights     string `form:"weights"`
	LoadBalance string `form:"load_balance"`
}

type TlsInfo struct {
	Type       string `form:"type"`
	ClientAuth string `form:"client_auth"`
}

type Pages struct {
	PageNum  int `form:"page_num"`
	PageSize int `form:"page_size"`
}
