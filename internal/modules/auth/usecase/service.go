package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	authports "github.com/dionis-the-shark/apms-task-tracker/internal/modules/auth/ports"
	usermodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/user"
	"github.com/google/uuid"
)

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidInput       = errors.New("invalid input")
)

type Service struct {
	repo        authports.Repository
	tokenTTL    time.Duration
	tokenSecret []byte
}

type API interface {
	Register(ctx context.Context, input RegisterInput) (TokenResponse, error)
	Login(ctx context.Context, input LoginInput) (TokenResponse, error)
	ValidateToken(ctx context.Context, token string) (ValidateResponse, error)
}

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type ValidateResponse struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

func New(repo authports.Repository) *Service {
	secret := os.Getenv("AUTH_TOKEN_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}
	ttl := 24 * time.Hour
	if envTTL := strings.TrimSpace(os.Getenv("AUTH_TOKEN_TTL")); envTTL != "" {
		if parsed, err := time.ParseDuration(envTTL); err == nil {
			ttl = parsed
		}
	}
	return &Service{
		repo:        repo,
		tokenTTL:    ttl,
		tokenSecret: []byte(secret),
	}
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (TokenResponse, error) {
	if input.Name == "" || input.Email == "" || input.Password == "" || input.Role == "" {
		return TokenResponse{}, fmt.Errorf("%w: name, email, password, and role are required", ErrInvalidInput)
	}

	if _, err := s.repo.GetUserByEmail(ctx, input.Email); err == nil {
		return TokenResponse{}, ErrEmailExists
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return TokenResponse{}, err
	}

	hash, err := hashPassword(input.Password)
	if err != nil {
		return TokenResponse{}, err
	}

	u := &usermodule.User{
		UserID:       uuid.New(),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hash,
		Role:         input.Role,
	}
	if err := s.repo.CreateUser(ctx, u); err != nil {
		return TokenResponse{}, err
	}
	token, err := s.issueToken(u.UserID.String(), u.Role)
	if err != nil {
		return TokenResponse{}, err
	}
	return TokenResponse{Token: token}, nil
}

func (s *Service) Login(ctx context.Context, input LoginInput) (TokenResponse, error) {
	if input.Email == "" || input.Password == "" {
		return TokenResponse{}, fmt.Errorf("%w: email and password are required", ErrInvalidInput)
	}

	u, err := s.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TokenResponse{}, ErrInvalidCredentials
		}
		return TokenResponse{}, err
	}

	ok, err := verifyPassword(input.Password, u.PasswordHash)
	if err != nil {
		return TokenResponse{}, err
	}
	if !ok {
		return TokenResponse{}, ErrInvalidCredentials
	}

	token, err := s.issueToken(u.UserID.String(), u.Role)
	if err != nil {
		return TokenResponse{}, err
	}
	return TokenResponse{Token: token}, nil
}

func (s *Service) ValidateToken(ctx context.Context, token string) (ValidateResponse, error) {
	_ = ctx

	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return ValidateResponse{}, ErrInvalidToken
	}

	payloadB64 := parts[0]
	sigB64 := parts[1]

	expectedSig := sign(payloadB64, s.tokenSecret)
	if !hmac.Equal([]byte(expectedSig), []byte(sigB64)) {
		return ValidateResponse{}, ErrInvalidToken
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return ValidateResponse{}, ErrInvalidToken
	}

	var payload tokenPayload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return ValidateResponse{}, ErrInvalidToken
	}

	if payload.Exp < time.Now().Unix() {
		return ValidateResponse{}, ErrTokenExpired
	}
	if payload.Sub == "" || payload.Role == "" {
		return ValidateResponse{}, ErrInvalidToken
	}

	return ValidateResponse{UserID: payload.Sub, Role: payload.Role}, nil
}

type tokenPayload struct {
	Sub  string `json:"sub"`
	Role string `json:"role"`
	Exp  int64  `json:"exp"`
}

func (s *Service) issueToken(userID, role string) (string, error) {
	payload := tokenPayload{
		Sub:  userID,
		Role: role,
		Exp:  time.Now().Add(s.tokenTTL).Unix(),
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadJSON)
	sigB64 := sign(payloadB64, s.tokenSecret)
	return fmt.Sprintf("%s.%s", payloadB64, sigB64), nil
}

func sign(payloadB64 string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(payloadB64))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	sum := sha256.Sum256(append(salt, []byte(password)...))
	return fmt.Sprintf("%s$%s",
		base64.RawURLEncoding.EncodeToString(salt),
		base64.RawURLEncoding.EncodeToString(sum[:]),
	), nil
}

func verifyPassword(password, stored string) (bool, error) {
	parts := strings.Split(stored, "$")
	if len(parts) != 2 {
		return false, errors.New("invalid password hash format")
	}
	salt, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}
	expected, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}
	sum := sha256.Sum256(append(salt, []byte(password)...))
	if len(expected) != len(sum) {
		return false, nil
	}
	return subtle.ConstantTimeCompare(expected, sum[:]) == 1, nil
}
