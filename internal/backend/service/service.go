package service

import (
	"api_gateway/internal/backend/config"
	"api_gateway/internal/backend/models"
	"api_gateway/internal/backend/payload"
	gatewayConfig "api_gateway/internal/gateway/config"
	"api_gateway/internal/gateway/dynamic"
	"api_gateway/pkg/logs"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
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
	BackendService      = "backend"
)

type Provider interface {
	ReloadConfig(msg dynamic.Message)
	Name() string
}

type Service struct {
	Engine      *gin.Engine
	Addr        string
	conf        config.WebServer
	TokenExpire time.Duration
	JwtKeyBytes []byte
	provider    Provider
	logger      zerolog.Logger
}

func (s *Service) ReloadAllEndpoint() {
	var cfgs []gatewayConfig.Endpoint
	endpoints, err := models.GetAllEndpoints()
	if err != nil {
		s.logger.Error().Err(err).Msg("error when reload endpoints")
	}

	for _, endpoint := range endpoints {
		var routers []gatewayConfig.Router

		cfg := gatewayConfig.Endpoint{
			Name:       endpoint.Name,
			ListenPort: endpoint.ListenPort,
			Type:       gatewayConfig.EndpointType(endpoint.Type),
		}
		for _, r := range endpoint.Routers {
			var (
				upstream    = payload.UpstreamInfo{}
				middlewares []gatewayConfig.Middleware
			)

			err = json.Unmarshal([]byte(r.UpStream), &upstream)
			if err != nil {
				s.logger.Error().Err(err).Msg("Error when load upstream config")
				continue
			}

			err = json.Unmarshal([]byte(r.Middlewares), &middlewares)
			if err != nil {
				s.logger.Error().Err(err).Msg("Error when load upstream config")
				continue
			}

			route := gatewayConfig.Router{
				Rule:       r.Rule,
				Host:       r.Host,
				TlsEnabled: r.TlsEnable,
				Type:       gatewayConfig.RuleType(r.Type),
				Priority:   r.Priority,
				Upstream: gatewayConfig.Upstream{
					Type:                upstream.Type,
					Paths:               upstream.Path,
					Weights:             upstream.Weights,
					LoadBalancerType:    upstream.LoadBalance,
					MaxIdleConnsPerHost: upstream.MaxIdleConnsPerHost,
				},
				Middlewares: middlewares,
				TLSConfig: gatewayConfig.TLS{
					Type:       "bytes",
					CsrFile:    r.Cert.Csr,
					KeyFile:    r.Cert.Key,
					CaFiles:    nil,
					ClientAuth: r.Cert.ClientAuth,
				},
			}

			routers = append(routers, route)
		}

		cfg.Routers = routers
		cfgs = append(cfgs, cfg)
	}
	msg := dynamic.Message{
		ProviderName:  s.provider.Name(),
		Configuration: cfgs,
	}

	s.provider.ReloadConfig(msg)
}

func NewService(conf config.WebServer, provider Provider) *Service {
	logger := log.With().Str(logs.ServiceName, BackendService).Logger()

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
		logger:      logger,
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
