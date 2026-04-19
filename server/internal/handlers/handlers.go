package handlers

import (
	"log/slog"

	"github.com/diyorbek/sentinel/internal/config"
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	h.SetupPublicRoutes(router)

	// Routes
	h.SetupRoutes(router)

	return router
}

func (h *Handler) SetupPublicRoutes(router *gin.Engine) {
	router.POST("/register", h.register)
	router.POST("/login", h.login)
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	users := router.Group("/sentinal/api/users")
	{
		users.POST("/", h.CreateUser)
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", h.UpdateUser)
		users.PUT("/:id/role", h.UpdateUserRole)
		users.DELETE("/:id", h.DeleteUser)
	}
	agents := router.Group("/sentinal/api/agents")
	{
		agents.POST("/", h.CreateAgent)
		agents.GET("/:id", h.GetAgentByID)
		agents.GET("/", h.ListAgents)
		agents.PUT("/:id/lastseen", h.UpdateLastSeen)
		agents.DELETE("/:id", h.DeleteAgent)
	}
	logs := router.Group("/sentinal/api")
	{
		logs.GET("/applogs/:id", h.GetAppLogByID)
		logs.GET("/applogs", h.GetListAppLog)
		logs.GET("/nginxlogs/:id", h.GetNginxLogByID)
		logs.GET("/nginxlogs", h.ListNginxLogs)
		logs.GET("/metrics/:id", h.GetMetricsByID)
		logs.GET("/metrics", h.ListMetrics)
	}
	alerts := router.Group("/sentinal/alerts")
	{
		alerts.GET("/:id", h.GetAlertByID)
		alerts.GET("/", h.ListAlerts)
		alerts.PUT("/:id/markread", h.MarkAlertRead)
	}
}
