package postgres

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type agentRepo struct {
	db *sqlx.DB
	logger *slog.Logger
}

func NewAgentRepo(db *sqlx.DB, logger *slog.Logger) *agentRepo {
	return &agentRepo{
		db: db,
		logger: logger,
	}
}

func (r *agentRepo) CreateAgent()