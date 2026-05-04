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

type accountService struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewAccountService(repo *repository.Repository, cfg *config.Config) *accountService {
	return &accountService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *accountService) CreateAccount(req models.CreateAccount) (uuid.UUID, error) {
	req.Id = uuid.New()

	apiKey, err := generateAPIKey()
	if err != nil {
		return uuid.Nil, apperrors.Internal(err)
	}
	req.APIKey = apiKey

	id, err := s.repo.Account.CreateAccount(req)
	if err != nil {
		return uuid.Nil, apperrors.Internal(err)
	}

	return id, nil
}

func (s *accountService) GetAccountByID(id string) (models.Account, error) {
	UId, err := uuid.Parse(id)
	if err != nil {
		return models.Account{}, apperrors.BadRequest("noto'g'ri UUID format")
	}

	account, err := s.repo.Account.GetByID(UId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Account{}, apperrors.NotFound(
				"user not found",
				err,
			)
		}
		return models.Account{}, apperrors.Internal(err)
	}

	return account, nil
}

func (s *accountService) GetAccountByAPIKey(apiKey string) (uuid.UUID, error) {
	account, err := s.repo.Account.GetByAPIKey(apiKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, apperrors.NotFound(
				"user not found",
				err,
			)
		}
		return uuid.Nil, apperrors.Internal(err)
	}

	return account.Id, err
}

func (s *accountService) UpdateAccount(req models.UpdateAccount) error {
	UId, err := uuid.Parse(req.Id)
	if err != nil {
		return apperrors.BadRequest("noto'g'ri UUID format")
	}

	err = s.repo.Account.UpdateAccount(models.UpdateAccountDB{
		Id:       UId,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNoRowsAffected) {
			return apperrors.BadRequest("bunday ma'lumotga ega user yo'q")
		}
		return apperrors.Internal(err)
	}

	return nil
}

func (s *accountService) DeleteAccount(id string) error {
	UId, err := uuid.Parse(id)
	if err != nil {
		return apperrors.BadRequest("noto'g'ri UUID format")
	}

	err = s.repo.Account.DeleteAccount(UId)
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNoRowsAffected) {
			return apperrors.BadRequest("bunday ma'lumotga ega user yo'q")
		}
		return apperrors.Internal(err)
	}

	return nil
}
