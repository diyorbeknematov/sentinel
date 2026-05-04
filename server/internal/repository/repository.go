package repository

import (
	"context"
	"time"

	"github.com/diyorbek/sentinel/internal/models"
	"github.com/diyorbek/sentinel/internal/repository/postgres"
	r "github.com/diyorbek/sentinel/internal/repository/redis"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	Account
	Agent
	Alert
	AppLog
	NginxLog
	Metric
	Stats
	RedisStore
}

func NewRepository(db *sqlx.DB, rdb *redis.Client) *Repository {
	return &Repository{
		Account:    postgres.NewAccountRepo(db),
		Agent:      postgres.NewAgentRepo(db),
		Alert:      postgres.NewAlertRepo(db),
		AppLog:     postgres.NewAppLogRepo(db),
		NginxLog:   postgres.NewNGINXLogRepo(db),
		Metric:     postgres.NewMetricRepo(db),
		Stats:      postgres.NewStatsRepo(db),
		RedisStore: r.NewRedisStore(rdb),
	}
}

type Account interface {
	CreateAccount(models.CreateAccount) (uuid.UUID, error)
	GetByEmail(string) (models.Account, error)
	GetByUsername(string) (models.Account, error)
	GetByID(uuid.UUID) (models.Account, error)
	GetByAPIKey(string) (models.Account, error)
	UpdateAccount(models.UpdateAccountDB) error
	DeleteAccount(uuid.UUID) error
}

type Agent interface {
	CreateAgent(models.CreateAgentDB) (uuid.UUID, error)
	GetAgentByID(uuid.UUID) (models.Agent, error)
	ListAgents(models.FilterAgentDB) ([]models.Agent, int, error)
	UpdateLastSeen(models.UpdateLastSeen) error
	DeleteAgent(uuid.UUID) error
}

type Alert interface {
	CreateAlert(alert models.CreateAlert) (uuid.UUID, error)
	GetAlertByID(uuid.UUID) (models.Alert, error)
	ListAlerts(models.FilterAlertDB) ([]models.AlertResponse, int, error)
	MarkAlertRead(uuid.UUID) error
}

type AppLog interface {
	CreateAppLog(models.CreateAppLog) (uuid.UUID, error)
	GetLogByID(uuid.UUID) (models.Log, error)
	ListLogs(models.FilterAppLogDB) ([]models.AppLogResponse, int, error)
}

type Metric interface {
	CreateMetric(models.CreateMetric) (uuid.UUID, error)
	GetMetricsByID(uuid.UUID) (models.Metric, error)
	ListMetrics(models.FilterMetricsDB) ([]models.MetricResponse, int, error)
}

type NginxLog interface {
	CreateNginxLog(log models.CreateNginxLog) (uuid.UUID, error)
	GetNginxLogByID(id uuid.UUID) (models.NginxLog, error)
	ListNginxLogs(filter models.FilterNginxLogDB) ([]models.NginxLogResponse, int, error)
}

type Stats interface {
	GetDashboardStats(models.StatsFilter, time.Duration) (models.DashboardStats, error)
	GetLogVolume(uuid.UUID) ([]models.LogVolumeStats, error)
}

type RedisStore interface {
	SaveResetToken(ctx context.Context, token, userID string) error
	GetResetToken(ctx context.Context, token string) (string, error)
	DeleteResetToken(ctx context.Context, token string) error
}
