package service

import (
	"log/slog"
	"time"

	"github.com/diyorbek/sentinel/internal/config"
	"github.com/diyorbek/sentinel/internal/models"
	"github.com/diyorbek/sentinel/internal/repository"
	"github.com/google/uuid"
)

type Service struct {
	User
	Authorization
	Agent
	Alert
	AppLog
	NginxLog
	Metric
}

func NewService(repos *repository.Repository, cfg *config.Config, logger *slog.Logger) *Service {
	return &Service{
		User:          NewUserService(repos, cfg),
		Authorization: NewAuthService(repos, cfg),
		Agent:         NewAgentService(repos, cfg, logger),
		Alert:         NewAlertService(repos, cfg, logger),
		AppLog:        NewLogService(repos, cfg, logger),
		Metric:        NewMetricService(repos, cfg, logger),
	}
}

type User interface {
	CreateUser(req models.CreateUser) (uuid.UUID, error)
	GetUser(id uuid.UUID) (models.User, error)
	UpdateUser(req models.UpdateUser) error
	UpdateUserRole(req models.UpdateRole) error
	DeleteUser(id uuid.UUID) error
}

type Authorization interface {
	CreateToken(models.User, string, time.Time) (*models.Token, error)
	GenerateTokens(models.User) (*models.Token, *models.Token, error)
	ParseToken(string) (*jwtCustomClaim, error)
	Login(models.Login) (*models.Token, *models.Token, error)
	Register(models.Register) (*models.Token, *models.Token, error)
}

type Agent interface {
	CreateAgent(req models.CreateAgent) (uuid.UUID, string, error)
	GetAgentByID(id uuid.UUID) (models.Agent, error)
	GetAgentByAPIKey(apiKey string) (models.Agent, error)
	ListAgents(filter models.FilterAgent) ([]models.Agent, int, error)
	UpdateLastSeen(lastSeen models.UpdateLastSeen) error
	DeleteAgent(id uuid.UUID) error
}

type Alert interface {
	CreateAlert(alert models.CreateAlert) (uuid.UUID, error)
	GetAlertByID(id uuid.UUID) (models.Alert, error)
	ListAlerts(filter models.FilterAlert) ([]models.Alert, int, error)
	MarkAlertRead(isRead models.MarkAlertRead) error
}

type AppLog interface {
	CreateAppLog(log models.CreateAppLog) (uuid.UUID, error)
	GetLogByID(id uuid.UUID) (models.Log, error)
	ListLogs(filter models.FilterAppLog) ([]models.Log, int, error)
}

type Metric interface {
	CreateMetric(metric models.CreateMetric) (uuid.UUID, error)
	GetMetricsByID(id uuid.UUID) (models.Metric, error)
	ListMetrics(filter models.FilterMetrics) ([]models.Metric, int, error)
}

type NginxLog interface {
	CreateNginxLog(log models.CreateNginxLog) (uuid.UUID, error)
	GetNginxLogByID(id uuid.UUID) (models.NginxLog, error)
	ListNginxLogs(filter models.FilterNginxLog) ([]models.NginxLog, int, error)
}
