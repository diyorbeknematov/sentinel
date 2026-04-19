package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/diyorbek/sentinel/internal/config"
	"github.com/diyorbek/sentinel/internal/models"
	"github.com/diyorbek/sentinel/internal/repository"
	"github.com/diyorbek/sentinel/pkg/apperrors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type authService struct {
	repo *repository.Repository
	cfg  *config.Config
}

type jwtCustomClaim struct {
	jwt.StandardClaims
	UserId uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	Type   string    `json:"type"`
}

func NewAuthService(repo *repository.Repository, cfg *config.Config) *authService {
	return &authService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *authService) CreateToken(user models.User, tokenType string, expiresAt time.Time) (*models.Token, error) {
	claims := jwtCustomClaim{
		UserId: user.Id,
		Role:   user.Role,
		Type:   tokenType,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expiresAt.Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	token, err := jwtToken.SignedString([]byte(s.cfg.ACCESS_TOKEN))
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	return &models.Token{
		Token:     token,
		Type:      tokenType,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *authService) GenerateTokens(user models.User) (*models.Token, *models.Token, error) {
	accessExpiresAt := time.Now().Add(time.Hour)
	refreshExpiresAt := time.Now().Add(24 * time.Hour)

	accessToken, err := s.CreateToken(user, "access_token", accessExpiresAt)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := s.CreateToken(user, "refresh_token", refreshExpiresAt)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) ParseToken(token string) (*jwtCustomClaim, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwtCustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.cfg.ACCESS_TOKEN), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*jwtCustomClaim)
	if !ok {
		return nil, errors.New("token claims are not of type *jwtCustomClaim")
	}

	return claims, nil
}

func (s *authService) Register(req models.Register) (*models.Token, *models.Token, error) {
	// Check if the username already exists
	_, err := s.repo.User.GetUserByUserName(req.UserName)
	if err == nil {
		return nil, nil, apperrors.BadRequest("The username is already existed")
	} else if err != sql.ErrNoRows {
		return nil, nil, apperrors.Internal(err)
	}

	// hashing the password
	hashedPass, err := HassPassword(req.Password)
	if err != nil {
		return nil, nil, err
	}
	req.Password = hashedPass

	userId, err := s.repo.User.CreateUser(models.CreateUser{
		UserName: req.UserName,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		return nil, nil, err
	}

	return s.GenerateTokens(models.User{
		Id:       userId,
		Role:     req.Role,
		UserName: req.UserName,
		Password: req.Password,
	})
}

func (s *authService) Login(req models.Login) (*models.Token, *models.Token, error) {
	user, err := s.repo.User.GetUserByUserName(req.UserName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, apperrors.BadRequest("bunday foydanaluvchi nomi yo'q")
		}
		return nil, nil, apperrors.Internal(err)
	}

	if VerifyPassword(req.Password, user.Password) {
		return nil, nil, apperrors.BadRequest("username yoki password xato")
	}

	return s.GenerateTokens(user)
}
