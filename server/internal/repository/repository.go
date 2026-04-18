package repository

import (
	"github.com/diyorbek/sentinel/internal/models"
	"github.com/diyorbek/sentinel/internal/repository/postgres"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Agent
	Alert
	AppLog
	NginxLog
	Metric
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Agent:    postgres.NewAgentRepo(db),
		Alert:    postgres.NewAlertRepo(db),
		AppLog:   postgres.NewAppLogRepo(db),
		NginxLog: postgres.NewNGINXLogRepo(db),
		Metric:   postgres.NewMetricRepo(db),
	}
}

type Agent interface {
	CreateAgent(models.CreateAgent) (uuid.UUID, error)
	GetAgentByID(uuid.UUID) (models.Agent, error)
	GetAgentByAPIKey(string) (models.Agent, error)
	ListAgents(models.FilterAgent) ([]models.Agent, int, error)
	UpdateLastSeen(models.UpdateLastSeen) error
	DeleteAgent(uuid.UUID) error
}

type Alert interface {
	CreateAlert(alert models.CreateAlert) (uuid.UUID, error)
	GetAlertByID(uuid.UUID) (models.Alert, error)
	ListAlerts(models.FilterAlert) ([]models.Alert, int, error)
	MarkAlertRead(models.MarkAlertRead) error
}

type AppLog interface {
	CreateAppLog(models.CreateAppLog) (uuid.UUID, error)
	GetLogByID(uuid.UUID) (models.Log, error)
	ListLogs(models.FilterAppLog) ([]models.Log, int, error)
}

type Metric interface {
	CreateMetric(models.CreateMetric) (uuid.UUID, error)
	GetMetricsByID(uuid.UUID) (models.Metric, error)
	ListMetrics(models.FilterMetrics) ([]models.Metric, int, error)
}

type NginxLog interface{}
