package service

import (
	"log/slog"

	"github.com/diyorbek/sentinel/internal/config"
	"github.com/diyorbek/sentinel/internal/models"
	"github.com/diyorbek/sentinel/internal/repository"
	"github.com/diyorbek/sentinel/pkg/apperrors"
	"github.com/google/uuid"
)

type logService struct {
	repo   *repository.Repository
	cfg    *config.Config
	logger *slog.Logger
}

func NewLogService(repo *repository.Repository, cfg *config.Config, logger *slog.Logger) *logService {
	return &logService{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *logService) CreateAppLog(log models.CreateAppLog) (uuid.UUID, error) {
	// 1. DB ga yozish
	id, err := s.repo.AppLog.CreateAppLog(log)
	if err != nil {
		return uuid.Nil, apperrors.Internal(err)
	}

	return id, nil
}

func (s *logService) GetLogByID(id uuid.UUID) (models.Log, error) {
	log, err := s.repo.AppLog.GetLogByID(id)
	if err != nil {
		if apperrors.Is(err, apperrors.EerrNoRowsAffected) {
			return models.Log{}, apperrors.NotFound("log topilmadi", err)
		}
		return models.Log{}, apperrors.Internal(err)
	}

	return log, nil
}

func (s *logService) ListLogs(filter models.FilterAppLog) ([]models.Log, int, error) {
	logs, total, err := s.repo.AppLog.ListLogs(filter)
	if err != nil {
		if apperrors.Is(err, apperrors.EerrNoRowsAffected) {
			return []models.Log{}, 0, apperrors.NotFound("log topilmadi", err)
		}
		return []models.Log{}, 0, apperrors.Internal(err)
	}

	return logs, total, nil
}
