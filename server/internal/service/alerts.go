package service

import (
	"database/sql"
	"errors"

	"github.com/diyorbek/sentinel/internal/config"
	"github.com/diyorbek/sentinel/internal/models"
	"github.com/diyorbek/sentinel/internal/repository"
	"github.com/diyorbek/sentinel/pkg/apperrors"
	"github.com/google/uuid"
)

type alertService struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewAlertService(repo *repository.Repository, cfg *config.Config) *alertService {
	return &alertService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *alertService) CreateAlert(alert models.CreateAlert) (uuid.UUID, error) {
	alert.Id = uuid.New()

	id, err := s.repo.Alert.CreateAlert(alert)
	if err != nil {
		return uuid.Nil, apperrors.Internal(err)
	}

	return id, nil
}

func (s *alertService) GetAlertByID(id string) (models.Alert, error) {
	UId, err := uuid.Parse(id)
	if err != nil {
		return models.Alert{}, apperrors.BadRequest("noto'g'ri UUID format")
	}
	alert, err := s.repo.Alert.GetAlertByID(UId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Alert{}, apperrors.NotFound("alert topilmadi", err)
		}
		return models.Alert{}, apperrors.Internal(err)
	}

	return alert, nil
}

func (s *alertService) ListAlerts(filter models.FilterAlert) ([]models.AlertResponse, int, error) {
	var agentUId *uuid.UUID = nil
	if filter.AgentId != nil {
		id, err := uuid.Parse(*filter.AgentId)
		if err != nil {
			return []models.AlertResponse{}, 0, apperrors.BadRequest("noto'g'ri UUID format")
		}
		agentUId = &id
	}

	alerts, total, err := s.repo.Alert.ListAlerts(models.FilterAlertDB{
		AgentId:  agentUId,
		Severity: filter.Severity,
		IsRead:   filter.IsRead,
		From:     filter.From,
		To:       filter.To,
		Limit:    filter.Limit,
		Offset:   filter.Offset,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.AlertResponse{}, 0, apperrors.NotFound("alert topilmadi", err)
		}
		return []models.AlertResponse{}, 0, apperrors.Internal(err)
	}

	return alerts, total, nil
}

func (s *alertService) MarkAlertRead(id string) (err error) {
	UId, err := uuid.Parse(id)
	if err != nil {
		return apperrors.BadRequest("noto'g'ri UUID format")
	}

	err = s.repo.Alert.MarkAlertRead(UId)
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNoRowsAffected) {
			return apperrors.BadRequest("alert mavjud emas")
		}
		return apperrors.Internal(err)
	}
	return
}
