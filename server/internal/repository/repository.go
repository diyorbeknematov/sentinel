package database

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	
}

func NewRepository(db *sqlx.DB, logger *slog.Logger) *Repository {
	return &Repository{}
}
