package backend

import (
	"api_gateway/internal/backend/config"
	"api_gateway/internal/backend/web"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
	"mime"
	"net/http"
)

func CORS(ctx *gin.Context) {
	method := ctx.Request.Method

	// set response header
	ctx.Header("Access-Control-Allow-Origin", ctx.Request.Header.Get("Origin"))
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Access-Control-Allow-Headers",
		"Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With, X-Files")
	ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")

	if method == "OPTIONS" || method == "HEAD" {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.Next()
}

// exportHeaders export header Content-Disposition for axios
func exportHeaders(ctx *gin.Context) {
	ctx.Header("Access-Control-Expose-Headers", "Content-Disposition")
	ctx.Next()
}

type HandlerFunc func(c *web.Context)

func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &web.Context{
			Context: c,
		}
		h(ctx)
	}
}

func InitRouter(s *web.Service) *web.Service {
	r := gin.New()

	r.Use(gin.Recovery()).Use(CORS).Use(exportHeaders)

	// swagger docs
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	err := mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		log.Error().AnErr("error when add extension type, err: %v", err)
	}

	// restapi
	apiV1 := r.Group("/api/v1")
	{
		apiV1.POST("/oauth/register/admin", Handle(s.RegisterAdmin))
		apiV1.POST("/oauth/register/tenant", Handle(s.RegisterTenant))
		apiV1.POST("/oauth/login/password", Handle(s.GetVersion))
		// version
		apiV1.GET("/version", Handle(s.GetVersion))
	}
	s.Engine = r

	return s
}

func Serve(conf config.WebServer) {
	gin.SetMode(gin.ReleaseMode)
	s := InitRouter(web.NewService(conf))
	log.Info().Msg(fmt.Sprintf("Listening and serving HTTP on %s", s.Addr))
	go func() {
		if err := s.Engine.Run(s.Addr); err != nil {
			log.Error().AnErr(fmt.Sprintf("Listening and serving HTTP on %s", s.Addr), err)

		}
	}()
}
