package handler

import (
	"api_gateway/internal/backend/config"
	"api_gateway/internal/backend/payload"
	"api_gateway/internal/gateway/dynamic"
	"api_gateway/pkg/logs"
	"api_gateway/pkg/safe"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	HttpStatusOk        = 200
	HttpStatusError     = 400
	HttpResponseSuccess = "success"
	ServiceBackend      = "backend"
)

type Provider interface {
	Provide(configurationChan chan<- dynamic.Message, pool *safe.Pool) error
}

type Service struct {
	Engine      *gin.Engine
	Addr        string
	conf        config.WebServer
	TokenExpire time.Duration
	JwtKeyBytes []byte
	provider    Provider
}

func NewService(conf config.WebServer, provider Provider) *Service {
	logger := log.With().Str(logs.ServiceName, ServiceBackend).Logger()

	pkFile, err := os.Open(conf.Jwt.JwtSecretPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot open private key")
	}

	pkBytes, err := io.ReadAll(pkFile)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot read private key")
	}

	service := &Service{
		Addr:        fmt.Sprintf("%s:%d", conf.BindAddr, conf.BindPort),
		conf:        conf,
		TokenExpire: conf.Jwt.JwtExp,
		JwtKeyBytes: pkBytes,
		provider:    provider,
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
