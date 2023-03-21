package payload

const (
	RespTypeMsg   = "msg"
	RespTypeError = "error"
	RespTypeData  = "data"

	SuccessMessage = "success"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data,omitempty"`
	Type string      `json:"type,omitempty"` // data msg error
}

func GenerateDataResponse(code int, msg string, data interface{}) Response {
	return Response{code, msg, data, RespTypeData}
}

func GenerateMsgResponse(code int, msg string) Response {
	return Response{code, msg, nil, RespTypeMsg}
}

func GenerateErrorResponse(code int, msg string) Response {
	return Response{code, msg, nil, RespTypeError}
}

type Page struct {
	PageNum  int `form:"page_num"`
	PageSize int `form:"page_size"`
}

type PageData struct {
	Total   int64       `json:"total"`
	PageNum int         `json:"page_num"`
	Data    interface{} `json:"data"`
}
