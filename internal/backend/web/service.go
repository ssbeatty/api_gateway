package web

import (
	"api_gateway/internal/backend/config"
	"api_gateway/internal/backend/payload"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

const (
	HttpStatusOk        = "200"
	HttpStatusError     = "400"
	HttpResponseSuccess = "success"
)

type Service struct {
	Engine *gin.Engine
	Addr   string
	conf   config.WebServer
}

func NewService(conf config.WebServer) *Service {
	service := &Service{
		Addr: fmt.Sprintf("%s:%d", conf.Addr, conf.Port),
		conf: conf,
	}

	return service
}

type Context struct {
	*gin.Context
}

func (c *Context) ResponseError(msg string) {
	d := payload.GenerateErrorResponse(HttpStatusError, msg)
	c.JSON(http.StatusOK, d)
}

func (c *Context) ResponseOk(data interface{}) {
	d := payload.GenerateDataResponse(HttpStatusOk, HttpResponseSuccess, data)
	c.JSON(http.StatusOK, d)
}

func (c *Context) readFormFile(header *multipart.FileHeader) ([]byte, error) {
	ff, err := header.Open()
	if err != nil {
		return nil, err
	}
	fileBytes, err := ioutil.ReadAll(ff)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}
