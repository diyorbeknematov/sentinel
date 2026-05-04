package handlers

import (
	"log/slog"

	"github.com/diyorbek/sentinel/internal/config"
	"github.com/diyorbek/sentinel/internal/handlers/middleware"
	"github.com/diyorbek/sentinel/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service *service.Service
	logger  *slog.Logger
}

func NewHandler(service *service.Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) InitRoutes(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	h.SetupPublicRoutes(router)

	// Routes
	h.SetupRoutes(router)

	return router
}

func (h *Handler) SetupPublicRoutes(router *gin.Engine) {
	router.POST("/sentinel/register", h.register)
	router.POST("/sentinel/login", h.login)
	router.POST("/sentinel/forgot-password", h.forgotPassword)
	router.POST("/sentinel/reset-password", h.resetPassword)
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	router.POST("/sentinel/agents", middleware.APIKeyMiddleware(h.service), h.CreateAgent)
	router.POST("/sentinel/heartbeat", middleware.APIKeyMiddleware(h.service), h.Heartbeat)

	router.Use(middleware.AuthMiddleware(h.service))
	router.GET("sentinel/me", h.GetMe)

	users := router.Group("/sentinel/accounts")
	{
		users.POST("/", h.CreateAccount)
		users.GET("/:id", h.GetAccount)
		users.PUT("/:id", h.UpdateAccount)
		users.DELETE("/:id", h.DeleteAccount)
	}
	agents := router.Group("/sentinel/agents")
	{
		agents.GET("/", h.ListAgents)
		agents.GET("/:id", h.GetAgentByID)
		agents.DELETE("/:id", h.DeleteAgent)
	}
	logs := router.Group("/sentinel")
	{
		logs.GET("/applogs/:id", h.GetAppLogByID)
		logs.GET("/applogs", h.GetListAppLog)
		logs.GET("/nginxlogs/:id", h.GetNginxLogByID)
		logs.GET("/nginxlogs", h.ListNginxLogs)
		logs.GET("/metrics/:id", h.GetMetricsByID)
		logs.GET("/metrics", h.ListMetrics)
		logs.GET("/stats", h.GetDashboardStats)
		logs.GET("/logvolume", h.GetLogVolumeStats)
	}
	alerts := router.Group("/sentinel/alerts")
	{
		alerts.GET("/:id", h.GetAlertByID)
		alerts.GET("/", h.ListAlerts)
		alerts.PUT("/:id/markread", h.MarkAlertRead)
	}
}
