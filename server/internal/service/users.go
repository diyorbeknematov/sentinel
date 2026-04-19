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

type userService struct {
	repo *repository.Repository
	cfg  *config.Config
}

func NewUserService(repo *repository.Repository, cfg *config.Config) *userService {
	return &userService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *userService) CreateUser(req models.CreateUser) (uuid.UUID, error) {
	id, err := s.repo.User.CreateUser(req)
	if err != nil {
		return uuid.Nil, apperrors.Internal(err)
	}

	return id, nil
}

func (s *userService) GetUser(id uuid.UUID) (models.User, error) {
	user, err := s.repo.User.GetUserByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, apperrors.NotFound(
				"user not found",
				err,
			)
		}
		return models.User{}, apperrors.Internal(err)
	}

	return user, nil
}

func (s *userService) UpdateUser(req models.UpdateUser) error {
	err := s.repo.User.UpdateUser(req)
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNoRowsAffected) {
			return apperrors.BadRequest("bunday ma'lumotga ega user yo'q")
		}
		return apperrors.Internal(err)
	}

	return nil
}

func (s *userService) UpdateUserRole(req models.UpdateRole) error {
	err := s.repo.User.UpdateUserRole(req)
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNoRowsAffected) {
			return apperrors.BadRequest("bunday ma'lumotga ega user yo'q")
		}
		return apperrors.Internal(err)
	}

	return nil
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	err := s.repo.User.DeleteUser(id)
	if err != nil {
		if apperrors.Is(err, apperrors.ErrNoRowsAffected) {
			return apperrors.BadRequest("bunday ma'lumotga ega user yo'q")
		}
		return apperrors.Internal(err)
	}

	return nil
}
