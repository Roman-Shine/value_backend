package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Roman-Shine/value_backend/internal/config"
	v1 "github.com/Roman-Shine/value_backend/internal/delivery/http/v1"
	"github.com/Roman-Shine/value_backend/internal/service"
	"github.com/Roman-Shine/value_backend/pkg/auth"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		//limiter.Limit(cfg.Limiter.RPS, cfg.Limiter.Burst, cfg.Limiter.TTL),
		//corsMiddleware,
	)

	//docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	//if cfg.Environment != config.EnvLocal {
	//	docs.SwaggerInfo.Host = cfg.HTTP.Host
	//}

	//if cfg.Environment != config.Prod {
	//	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//}

	// Init router
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, h.tokenManager)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
