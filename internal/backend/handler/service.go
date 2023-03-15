package handler

import (
	"api_gateway/internal/backend/config"
	"api_gateway/internal/backend/models"
	"api_gateway/internal/backend/payload"
	"api_gateway/internal/backend/utils"
	"api_gateway/pkg/logs"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	HttpStatusOk        = "200"
	HttpStatusError     = "400"
	HttpResponseSuccess = "success"
	ServiceBackend      = "backend"
)

// TokenGenerator generates a token for the specified account.
type TokenGenerator interface {
	GenerateToken(accountID string, expire time.Duration) (string, error)
}

type Service struct {
	Engine      *gin.Engine
	Addr        string
	conf        config.WebServer
	jwtGen      TokenGenerator
	TokenExpire time.Duration
}

func NewService(conf config.WebServer) *Service {
	logger := log.With().Str(logs.ServiceName, ServiceBackend).Logger()

	pkFile, err := os.Open(conf.Jwt.JwtSecretPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot open private key")
	}

	pkBytes, err := io.ReadAll(pkFile)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot read private key")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot parse private key")
	}

	jwtGen := utils.NewJWTTokenGen("/auth/admin/login", privateKey)
	service := &Service{
		Addr:        fmt.Sprintf("%s:%d", conf.BindAddr, conf.BindPort),
		conf:        conf,
		jwtGen:      jwtGen,
		TokenExpire: conf.Jwt.JwtExp,
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

func (s *Service) OauthLoginPassword(req payload.AdminLoginPasswordReq) (*payload.AdminLoginPasswordResp, error) {
	var resp = &payload.AdminLoginPasswordResp{}
	admin, err := models.GetAdminByUserName(req.UserName)
	if err != nil {
		return resp, errors.New("用户未注册")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		return resp, fmt.Errorf("密码不正确")
	}
	return s.SetUserJwtToken(admin.Username)
}

// SetUserJwtToken set admin/tenant token
func (s *Service) SetUserJwtToken(userName string) (*payload.AdminLoginPasswordResp, error) {
	var resp = &payload.AdminLoginPasswordResp{}
	token, err := s.jwtGen.GenerateToken(userName, s.TokenExpire)
	if err != nil {
		return resp, err
	}
	resp.UserName = userName
	resp.Toke = token
	resp.TokenExpireAt = int32(s.TokenExpire.Seconds())
	return resp, nil
}
