package service

import (
	"log/slog"
	"time"

	"github.com/diyorbek/sentinel/internal/config"
	"github.com/diyorbek/sentinel/internal/mailer"
	"github.com/diyorbek/sentinel/internal/models"
	"github.com/diyorbek/sentinel/internal/repository"
	"github.com/google/uuid"
)

type Service struct {
	Account
	Authorization
	Agent
	Alert
	AppLog
	NginxLog
	Metric
	Stats
}

func NewService(repos *repository.Repository, cfg *config.Config, logger *slog.Logger, mail *mailer.Mailer) *Service {
	return &Service{
		Account:       NewAccountService(repos, cfg),
		Authorization: NewAuthService(repos, cfg, mail),
		Agent:         NewAgentService(repos, cfg),
		Alert:         NewAlertService(repos, cfg),
		AppLog:        NewLogService(repos, cfg),
		NginxLog:      NewNginxLogService(repos, cfg),
		Metric:        NewMetricService(repos, cfg),
		Stats:         NewStatsService(repos, cfg),
	}
}

type Account interface {
	CreateAccount(models.CreateAccount) (uuid.UUID, error)
	GetAccountByID(string) (models.Account, error)
	GetAccountByAPIKey(string) (uuid.UUID, error)
	UpdateAccount(models.UpdateAccount) error
	DeleteAccount(string) error
}

type Authorization interface {
	CreateToken(models.Account, string, time.Time) (*models.Token, error)
	GenerateTokens(models.Account) (*models.Token, *models.Token, error)
	ParseToken(string) (*jwtCustomClaim, error)
	Login(models.Login) (*models.Token, *models.Token, error)
	Register(models.Register) (*models.Token, *models.Token, error)
	ForgotPassword(string) error
	ResetPassword(token string, newPassword string) error
}

type Agent interface {
	CreateAgent(models.CreateAgent) (models.CreateAgentResponse, error)
	GetAgentByID(string) (models.Agent, error)
	ListAgents(models.FilterAgent) ([]models.Agent, int, error)
	UpdateLastSeen(uuid.UUID) error
	DeleteAgent(string) error
}

type Alert interface {
	CreateAlert(alert models.CreateAlert) (uuid.UUID, error)
	GetAlertByID(string) (models.Alert, error)
	ListAlerts(filter models.FilterAlert) ([]models.AlertResponse, int, error)
	MarkAlertRead(string) error
}

type AppLog interface {
	CreateAppLog(log models.CreateAppLog) (uuid.UUID, error)
	GetLogByID(string) (models.Log, error)
	ListLogs(filter models.FilterAppLog) ([]models.AppLogResponse, int, error)
}

type Metric interface {
	CreateMetric(metric models.CreateMetric) (uuid.UUID, error)
	GetMetricsByID(string) (models.Metric, error)
	ListMetrics(filter models.FilterMetrics) ([]models.MetricResponse, int, error)
}

type NginxLog interface {
	CreateNginxLog(log models.CreateNginxLog) (uuid.UUID, error)
	GetNginxLogByID(string) (models.NginxLog, error)
	ListNginxLogs(filter models.FilterNginxLog) ([]models.NginxLogResponse, int, error)
}

type Stats interface {
	GetDashboardStats(filter models.StatsFilter) (models.DashboardStats, error)
	GetLogVolumeStats(string) ([]models.LogVolumeStats, error)
}
