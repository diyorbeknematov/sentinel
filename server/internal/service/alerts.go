package service

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/diyorbek/sentinel/internal/config"
	"github.com/diyorbek/sentinel/internal/models"
	"github.com/diyorbek/sentinel/internal/repository"
	"github.com/diyorbek/sentinel/pkg/apperrors"
	"github.com/google/uuid"
)

type alertService struct {
	repo   *repository.Repository
	cfg    *config.Config
	logger *slog.Logger
}

func NewAlertService(repo *repository.Repository, cfg *config.Config, logger *slog.Logger) *alertService {
	return &alertService{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *alertService) CreateAlert(alert models.CreateAlert) (uuid.UUID, error) {
	id, err := s.repo.Alert.CreateAlert(alert)
	if err != nil {
		return uuid.Nil, apperrors.Internal(err)
	}

	return id, nil
}

func (s *alertService) GetAlertByID(id uuid.UUID) (models.Alert, error) {
	alert, err := s.repo.Alert.GetAlertByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Alert{}, apperrors.NotFound("alert topilmadi", err)
		}
		return models.Alert{}, apperrors.Internal(err)
	}

	return alert, nil
}

func (s *alertService) ListAlerts(filter models.FilterAlert) ([]models.Alert, int, error) {
	alerts, total, err := s.repo.Alert.ListAlerts(filter)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.Alert{}, 0, apperrors.NotFound("alert topilmadi", err)
		}
		return []models.Alert{}, 0, apperrors.Internal(err)
	}

	return alerts, total, nil
}

func (s *alertService) MarkAlertRead(isRead models.MarkAlertRead) (err error) {
	err = s.repo.Alert.MarkAlertRead(isRead)
	if err != nil {
		if apperrors.Is(err, apperrors.EerrNoRowsAffected) {
			return apperrors.BadRequest("aler mavjud emas")
		}
		return apperrors.Internal(err)
	}
	return
}
