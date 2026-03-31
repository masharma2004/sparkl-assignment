package services

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"sparklassignment/backend/internal/config"
	"sparklassignment/backend/internal/dto"
	"sparklassignment/backend/internal/models"
	"sparklassignment/backend/internal/utils"
)

type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

type SessionBundle struct {
	AccessToken  string
	RefreshToken string
	Response     *dto.LoginResponse
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		db:  db,
		cfg: cfg,
	}
}

func toUserResponse(user *models.User) dto.UserResponse {
	return dto.UserResponse{
		UserID:   user.UserID,
		Username: user.Username,
		Role:     user.Role,
		FullName: user.FullName,
		Email:    user.Email,
	}
}

func (s *AuthService) issueSession(user *models.User) (*SessionBundle, error) {
	accessTokenTTL := time.Duration(s.cfg.AccessTokenTTLMinutes) * time.Minute
	refreshTokenTTL := time.Duration(s.cfg.RefreshTokenTTLDays) * 24 * time.Hour
	accessToken, err := utils.GenerateAccessToken(
		user.UserID,
		user.Role,
		s.cfg.JWTSecret,
		s.cfg.JWTIssuer,
		accessTokenTTL,
	)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateOpaqueToken(32)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	refreshSession := models.RefreshSession{
		UserID:    user.UserID,
		TokenHash: utils.HashToken(refreshToken),
		ExpiresAt: time.Now().UTC().Add(refreshTokenTTL),
	}

	if err := s.db.Create(&refreshSession).Error; err != nil {
		return nil, fmt.Errorf("create refresh session: %w", err)
	}

	return &SessionBundle{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Response: &dto.LoginResponse{
			User: toUserResponse(user),
		},
	}, nil
}

func (s *AuthService) LoginByRole(username, password, role string) (*SessionBundle, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, fmt.Errorf("find user by username: %w", err)
	}

	if user.Role != role {
		return nil, ErrInvalidCredentials
	}

	if err := utils.CheckPasswordHash(password, user.HashedPassword); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.issueSession(&user)
}

func (s *AuthService) SignupStudent(req dto.StudentSignupRequest) (*SessionBundle, error) {
	username := strings.TrimSpace(req.Username)
	fullName := strings.TrimSpace(req.FullName)
	email := strings.ToLower(strings.TrimSpace(req.Email))

	switch {
	case username == "":
		return nil, ValidationError{Message: "username is required"}
	case len(username) < 3:
		return nil, ValidationError{Message: "username must be at least 3 characters"}
	case len(username) > 50:
		return nil, ValidationError{Message: "username must be at most 50 characters"}
	case fullName == "":
		return nil, ValidationError{Message: "full_name is required"}
	case len(fullName) > 100:
		return nil, ValidationError{Message: "full_name must be at most 100 characters"}
	case email == "":
		return nil, ValidationError{Message: "email is required"}
	case len(email) > 50:
		return nil, ValidationError{Message: "email must be at most 50 characters"}
	case len(req.Password) < 8:
		return nil, ValidationError{Message: "password must be at least 8 characters"}
	case len(req.Password) > 72:
		return nil, ValidationError{Message: "password must be at most 72 characters"}
	}

	parsedEmail, err := mail.ParseAddress(email)
	if err != nil || !strings.EqualFold(parsedEmail.Address, email) {
		return nil, ValidationError{Message: "email must be a valid address"}
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := models.User{
		Username:       username,
		HashedPassword: hashedPassword,
		Role:           "student",
		FullName:       fullName,
		Email:          email,
	}

	if err := s.db.Create(&user).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ConflictError{Message: "student account already exists"}
		}

		return nil, fmt.Errorf("create student user: %w", err)
	}

	return s.issueSession(&user)
}

func (s *AuthService) RefreshSession(refreshToken string) (*SessionBundle, error) {
	trimmedToken := strings.TrimSpace(refreshToken)
	if trimmedToken == "" {
		return nil, ErrInvalidRefreshSession
	}

	now := time.Now().UTC()
	var bundle *SessionBundle

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var refreshSession models.RefreshSession
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("token_hash = ?", utils.HashToken(trimmedToken)).
			First(&refreshSession).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInvalidRefreshSession
			}

			return fmt.Errorf("find refresh session: %w", err)
		}

		if refreshSession.RevokedAt != nil || now.After(refreshSession.ExpiresAt) {
			return ErrInvalidRefreshSession
		}

		var user models.User
		if err := tx.First(&user, refreshSession.UserID).Error; err != nil {
			return fmt.Errorf("find refresh user: %w", err)
		}

		newRefreshToken, err := utils.GenerateOpaqueToken(32)
		if err != nil {
			return fmt.Errorf("generate refresh token: %w", err)
		}

		newRefreshHash := utils.HashToken(newRefreshToken)
		revokedAt := now
		refreshSession.RevokedAt = &revokedAt
		refreshSession.ReplacedByHash = newRefreshHash
		if err := tx.Save(&refreshSession).Error; err != nil {
			return fmt.Errorf("revoke refresh session: %w", err)
		}

		replacementSession := models.RefreshSession{
			UserID:    user.UserID,
			TokenHash: newRefreshHash,
			ExpiresAt: now.Add(time.Duration(s.cfg.RefreshTokenTTLDays) * 24 * time.Hour),
		}
		if err := tx.Create(&replacementSession).Error; err != nil {
			return fmt.Errorf("create rotated refresh session: %w", err)
		}

		accessToken, err := utils.GenerateAccessToken(
			user.UserID,
			user.Role,
			s.cfg.JWTSecret,
			s.cfg.JWTIssuer,
			time.Duration(s.cfg.AccessTokenTTLMinutes)*time.Minute,
		)
		if err != nil {
			return fmt.Errorf("generate rotated access token: %w", err)
		}

		bundle = &SessionBundle{
			AccessToken:  accessToken,
			RefreshToken: newRefreshToken,
			Response: &dto.LoginResponse{
				User: toUserResponse(&user),
			},
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return bundle, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	trimmedToken := strings.TrimSpace(refreshToken)
	if trimmedToken == "" {
		return nil
	}

	now := time.Now().UTC()
	return s.db.Model(&models.RefreshSession{}).
		Where("token_hash = ? AND revoked_at IS NULL", utils.HashToken(trimmedToken)).
		Updates(map[string]any{
			"revoked_at": revokedAtValue(now),
			"updated_at": now,
		}).Error
}

func (s *AuthService) GetMe(userID uint) (*dto.MeResponse, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &dto.MeResponse{
		User: toUserResponse(&user),
	}, nil
}

func revokedAtValue(value time.Time) *time.Time {
	return &value
}
