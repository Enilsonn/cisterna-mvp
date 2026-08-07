package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"cisterna-mvp/auth-service/internal/domain"
	"cisterna-mvp/auth-service/internal/messaging"

	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByID(ctx context.Context, id string) (domain.User, error)
	CreateUser(ctx context.Context, user domain.User) error
	SaveRefreshToken(ctx context.Context, token domain.RefreshToken) error
	RevokeRefreshToken(ctx context.Context, token string) error
	GetRefreshToken(ctx context.Context, token string) (domain.RefreshToken, error)
}

func (s *AuthService) RegisterUser(ctx context.Context, email, password string, role domain.Role) (domain.User, error) {
	if email == "" || password == "" {
		return domain.User{}, errors.New("email and password are required")
	}

	if _, err := s.repo.GetUserByEmail(ctx, email); err == nil {
		return domain.User{}, errors.New("user already exists")
	}

	switch role {
	case domain.RoleAdmin, domain.RolePipeiro, domain.RoleCidadao:
	default:
		return domain.User{}, errors.New("invalid role")
	}

	user := domain.User{
		ID:           fmt.Sprintf("user-%d", time.Now().UnixNano()),
		Email:        email,
		PasswordHash: HashPassword(password),
		Role:         role,
		IsActive:     true,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return domain.User{}, err
	}

	if s.publisher != nil {
		_ = s.publisher.Publish(ctx, messaging.NewAuthEvent("user_registered", user.Email, string(user.Role), user.ID))
	}

	return user, nil
}

type AuthService struct {
	repo      Repository
	publisher *messaging.Publisher
}

func NewAuthService(repo Repository) *AuthService {
	publisher := messaging.NewPublisher([]string{"localhost:9092"}, "auth-events")
	return &AuthService{repo: repo, publisher: publisher}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (domain.TokenPair, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.TokenPair{}, err
	}
	if !user.IsActive {
		return domain.TokenPair{}, errors.New("user inactive")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return domain.TokenPair{}, errors.New("invalid credentials")
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return domain.TokenPair{}, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return domain.TokenPair{}, err
	}

	if s.publisher != nil {
		_ = s.publisher.Publish(ctx, messaging.NewAuthEvent("user_logged_in", user.Email, string(user.Role), user.ID))
	}

	return domain.TokenPair{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		AccessExpiry:  time.Now().Add(15 * time.Minute).Unix(),
		RefreshExpiry: time.Now().Add(24 * time.Hour).Unix(),
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	stored, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err == nil && s.publisher != nil {
		user, errUser := s.repo.GetUserByID(ctx, stored.UserID)
		if errUser == nil {
			_ = s.publisher.Publish(ctx, messaging.NewAuthEvent("user_logged_out", user.Email, string(user.Role), user.ID))
		}
	}
	return s.repo.RevokeRefreshToken(ctx, refreshToken)
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (domain.TokenPair, error) {
	stored, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return domain.TokenPair{}, err
	}
	if stored.Revoked {
		return domain.TokenPair{}, errors.New("refresh token revoked")
	}

	user, err := s.repo.GetUserByID(ctx, stored.UserID)
	if err != nil {
		return domain.TokenPair{}, err
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return domain.TokenPair{}, err
	}

	newRefreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return domain.TokenPair{}, err
	}

	if s.publisher != nil {
		_ = s.publisher.Publish(ctx, messaging.NewAuthEvent("token_refreshed", user.Email, string(user.Role), user.ID))
	}

	return domain.TokenPair{
		AccessToken:   accessToken,
		RefreshToken:  newRefreshToken,
		AccessExpiry:  time.Now().Add(15 * time.Minute).Unix(),
		RefreshExpiry: time.Now().Add(24 * time.Hour).Unix(),
	}, nil
}

func (s *AuthService) ValidateAccessToken(token string) (domain.JWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return domain.JWTClaims{}, errors.New("invalid token")
	}

	payload, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return domain.JWTClaims{}, err
	}

	var claims struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return domain.JWTClaims{}, err
	}

	return domain.JWTClaims{UserID: claims.UserID, Role: domain.Role(claims.Role)}, nil
}

func (s *AuthService) generateAccessToken(user domain.User) (string, error) {
	claims := map[string]string{
		"user_id": user.ID,
		"role":    string(user.Role),
	}
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	header := []byte(`{"alg":"none","typ":"JWT"}`)
	token := base64.RawStdEncoding.EncodeToString(header) + "." + base64.RawStdEncoding.EncodeToString(payload)
	return token + ".signature", nil
}

func (s *AuthService) generateRefreshToken(user domain.User) (string, error) {
	raw := make([]byte, 24)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	token := hex.EncodeToString(raw)
	stored := domain.RefreshToken{
		ID:       token,
		UserID:   user.ID,
		UserRole: user.Role,
		Token:    token,
	}
	if err := s.repo.SaveRefreshToken(context.Background(), stored); err != nil {
		return "", err
	}
	return token, nil
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func ComparePassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

func (s *AuthService) CreateUser(ctx context.Context, user domain.User) error {
	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) ValidateRefreshToken(ctx context.Context, refreshToken string) (domain.RefreshToken, error) {
	return s.repo.GetRefreshToken(ctx, refreshToken)
}

func (s *AuthService) RemoveRefreshToken(ctx context.Context, refreshToken string) error {
	return s.repo.RevokeRefreshToken(ctx, refreshToken)
}

func (s *AuthService) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *AuthService) ValidateAccessTokenWithContext(ctx context.Context, token string) (domain.JWTClaims, error) {
	return s.ValidateAccessToken(token)
}

func (s *AuthService) LoginWithContext(ctx context.Context, email, password string) (domain.TokenPair, error) {
	return s.Login(ctx, email, password)
}

func (s *AuthService) RefreshTokensWithContext(ctx context.Context, refreshToken string) (domain.TokenPair, error) {
	return s.RefreshTokens(ctx, refreshToken)
}

func (s *AuthService) LogoutWithContext(ctx context.Context, refreshToken string) error {
	return s.Logout(ctx, refreshToken)
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (domain.JWTClaims, error) {
	return s.ValidateAccessToken(token)
}

func (s *AuthService) ParseClaims(token string) (domain.JWTClaims, error) {
	return s.ValidateAccessToken(token)
}

func (s *AuthService) CheckPassword(ctx context.Context, email, password string) (bool, error) {
	return ComparePasswordHash(email, password)
}

func ComparePasswordHash(email, password string) (bool, error) {
	return false, fmt.Errorf("not implemented")
}
