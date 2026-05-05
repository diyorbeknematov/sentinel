package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/diyorbek/sentinel/internal/config"
	"github.com/diyorbek/sentinel/internal/mailer"
	"github.com/diyorbek/sentinel/internal/models"
	"github.com/diyorbek/sentinel/internal/repository"
	"github.com/diyorbek/sentinel/pkg/apperrors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type authService struct {
	repo   *repository.Repository
	mailer *mailer.Mailer
	cfg    *config.Config
}

type jwtCustomClaim struct {
	jwt.StandardClaims
	AccountID uuid.UUID `json:"user_id"`
	Type      string    `json:"type"`
}

func NewAuthService(repo *repository.Repository, cfg *config.Config, mail *mailer.Mailer) *authService {
	return &authService{
		repo:   repo,
		mailer: mail,
		cfg:    cfg,
	}
}

func (s *authService) CreateToken(account models.Account, tokenType string, expiresAt time.Time) (*models.Token, error) {
	claims := jwtCustomClaim{
		AccountID: account.Id,
		Type:      tokenType,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expiresAt.Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString([]byte(s.cfg.AccessToken))
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	return &models.Token{
		Token:     token,
		Type:      tokenType,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *authService) GenerateTokens(account models.Account) (*models.Token, *models.Token, error) {
	accessExpiresAt := time.Now().Add(time.Hour)
	refreshExpiresAt := time.Now().Add(24 * time.Hour)

	accessToken, err := s.CreateToken(account, "access_token", accessExpiresAt)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := s.CreateToken(account, "refresh_token", refreshExpiresAt)
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

		return []byte(s.cfg.AccessToken), nil
	})
	if err != nil {
		return nil, err
	}

	if !jwtToken.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := jwtToken.Claims.(*jwtCustomClaim)
	if !ok {
		return nil, errors.New("token claims are not of type *jwtCustomClaim")
	}

	return claims, nil
}

func (s *authService) Register(req models.Register) (*models.Token, *models.Token, error) {
	// hashing the password
	hashedPass, err := HashPassword(req.Password)
	if err != nil {
		return nil, nil, err
	}
	req.Password = hashedPass

	apiKey, err := generateAPIKey()
	if err != nil {
		return nil, nil, err
	}

	userId, err := s.repo.Account.CreateAccount(models.CreateAccount{
		Id:       uuid.New(),
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		APIKey:   apiKey,
	})
	if err != nil {
		if isDuplicateError(err) {
			return nil, nil, apperrors.BadRequest("email or username already exists")
		}
		return nil, nil, apperrors.Internal(err)
	}

	return s.GenerateTokens(models.Account{
		Id:       userId,
		Email:    req.Email,
		Username: req.Username,
	})
}

func (s *authService) Login(req models.Login) (*models.Token, *models.Token, error) {
	user, err := s.repo.Account.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, apperrors.BadRequest("bunday foydanaluvchi nomi yo'q")
		}
		return nil, nil, apperrors.Internal(err)
	}

	if !VerifyPassword(req.Password, user.Password) {
		return nil, nil, apperrors.BadRequest("username yoki password xato")
	}

	return s.GenerateTokens(user)
}

func (s *authService) ForgotPassword(email string) error {
	account, err := s.repo.Account.GetByEmail(email)
	if err != nil {
		return apperrors.Internal(err)
	}

	token := generateToken()

	// DB ga saqlash (expiry bilan)
	s.repo.RedisStore.SaveResetToken(context.Background(), token, account.Id.String())

	link := fmt.Sprintf(
		"%s/reset-password?token=%s",
		s.cfg.FrontendURL,
		token,
	)

	body := fmt.Sprintf(`
Hello %s,

Reset password link:
%s

This link expires in 15 minutes.
`, account.Username, link)

	return s.mailer.Send(account.Email, "Reset Password", body)
}

func (s *authService) ResetPassword(token string, newPassword string) error {
	ctx := context.Background()

	// 1. tokenni tekshirish
	idStr, err := s.repo.RedisStore.GetResetToken(ctx, token)
	if err != nil {
		return apperrors.BadRequest("invalid or expired token")
	}

	userID, err := uuid.Parse(idStr)
	if err != nil {
		apperrors.BadRequest("noto'g'ri UUID format")
	}

	// 2. userni olish
	account, err := s.repo.Account.GetByID(userID)
	if err != nil {
		return apperrors.Internal(err)
	}

	// 3. passwordni hash qilish
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return apperrors.Internal(err)
	}

	// 4. update qilish
	account.Password = hashedPassword
	if err := s.repo.Account.UpdateAccount(models.UpdateAccountDB{
		Id: userID,
		Username: account.Username,
		Password: account.Password,
	}); err != nil {
		return apperrors.Internal(err)
	}

	// 5. tokenni o‘chirish (security uchun MUHIM)
	if err := s.repo.DeleteResetToken(ctx, token); err != nil {
		return apperrors.Internal(err)
	}

	return nil
}
