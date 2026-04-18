package service

import (
	"log/slog"

	"github.com/diyorbek/sentinel/internal/config"
	"github.com/diyorbek/sentinel/internal/models"
	"github.com/diyorbek/sentinel/internal/repository"
	"github.com/google/uuid"
)

type Service struct {
	Agent
	Alert
	AppLog
	Metric
}

func NewService(repos *repository.Repository, cfg *config.Config, logger *slog.Logger) *Service {
	return &Service{
		Agent:  NewAgentService(repos, cfg, logger),
		Alert:  NewAlertService(repos, cfg, logger),
		AppLog: NewLogService(repos, cfg, logger),
		Metric: NewMetricService(repos, cfg, logger),
	}
}

type Agent interface {
	CreateAgent(req models.CreateAgent) (uuid.UUID, error)
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
