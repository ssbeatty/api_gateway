package backend

import (
	"api_gateway/internal/backend/config"
	"api_gateway/internal/backend/handler"
	"api_gateway/internal/backend/models"
	"api_gateway/internal/backend/payload"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"mime"
	"time"
)

const (
	identityKey = "id"
)

// exportHeaders export header Content-Disposition for axios
func exportHeaders(ctx *gin.Context) {
	ctx.Header("Access-Control-Expose-Headers", "Content-Disposition")
	ctx.Next()
}

type HandlerFunc func(c *handler.Context)

func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &handler.Context{
			Context: c,
		}
		h(ctx)
	}
}

func InitRouter(s *handler.Service) *handler.Service {
	r := gin.New()

	corsConf := cors.DefaultConfig()
	corsConf.AllowAllOrigins = true
	r.Use(gin.Recovery()).Use(cors.New(corsConf)).Use(exportHeaders)

	err := mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		log.Error().AnErr("error when add extension type, err: %v", err)
	}

	// version
	r.GET("/version", Handle(s.GetVersion))

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "gateway backend",
		Key:         s.JwtKeyBytes,
		Timeout:     s.TokenExpire,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*payload.AdminLoginPasswordReq); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &payload.AdminLoginPasswordReq{
				UserName: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals payload.AdminLoginPasswordReq
			if err = c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.UserName
			password := loginVals.Password

			admin, errA := models.GetAdminByUserName(userID)
			if errA != nil {
				return nil, jwt.ErrForbidden
			}
			if _ = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
				return nil, jwt.ErrFailedAuthentication
			} else {
				return &loginVals, nil
			}

		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*payload.AdminLoginPasswordReq); ok {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, payload.Response{
				Code: code,
				Msg:  message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal().Err(err).Msg("JWT Error")
	}

	r.POST("/auth/admin/register", Handle(s.RegisterAdmin))
	// login
	r.POST("/auth/admin/login", authMiddleware.LoginHandler)
	r.POST("/auth/admin/logout", authMiddleware.LogoutHandler)
	r.POST("/auth/admin/refresh_token", authMiddleware.RefreshHandler)
	// restapi
	apiV1 := r.Group("/api/v1")
	apiV1.Use(authMiddleware.MiddlewareFunc())
	{
		apiV1.GET("/endpoints", Handle(s.EndpointsQuery))
		apiV1.POST("/endpoints", Handle(s.EndpointsCreate))
	}

	s.Engine = r

	return s
}

func Serve(conf config.WebServer) {
	gin.SetMode(gin.ReleaseMode)
	s := InitRouter(handler.NewService(conf))
	log.Info().Msg(fmt.Sprintf("Listening and serving HTTP on %s", s.Addr))
	go func() {
		if err := s.Engine.Run(s.Addr); err != nil {
			log.Error().AnErr(fmt.Sprintf("Listening and serving HTTP on %s", s.Addr), err)
		}
	}()
}
