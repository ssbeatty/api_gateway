package backend

import (
	"api_gateway/internal/backend/config"
	"api_gateway/internal/backend/docs"
	"api_gateway/internal/backend/handler"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"mime"
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

	// swagger docs
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	err := mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		log.Error().AnErr("error when add extension type, err: %v", err)
	}

	r.POST("/auth/admin/register", Handle(s.RegisterAdmin))
	r.POST("/auth/admin/login", Handle(s.AdminLoginPassword))
	// version
	r.GET("/version", Handle(s.GetVersion))

	// restapi
	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/endpoints", Handle(s.EndpointsList))
		apiV1.POST("/endpoints", Handle(s.EndpointsAdd))
		apiV1.DELETE("/endpoints", Handle(s.EndpointsDelete))
		apiV1.GET("/endpoints/:id", Handle(s.EndpointsGetDetail))
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
