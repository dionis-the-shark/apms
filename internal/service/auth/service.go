package authservice

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

	"github.com/dionis-the-shark/apms-task-tracker/internal/models/user"
	"github.com/google/uuid"
)

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
)

type Repository interface {
	CreateUser(ctx context.Context, u *user.User) error
	GetUserByEmail(ctx context.Context, email string) (user.User, error)
}

type Service struct {
	repo        Repository
	tokenTTL    time.Duration
	tokenSecret []byte
}

func New(repo Repository) *Service {
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

func (s *Service) Register(ctx context.Context, name, email, password, role string) (string, error) {
	if name == "" || email == "" || password == "" || role == "" {
		return "", errors.New("name, email, password, and role are required")
	}

	if _, err := s.repo.GetUserByEmail(ctx, email); err == nil {
		return "", ErrEmailExists
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	hash, err := hashPassword(password)
	if err != nil {
		return "", err
	}

	u := &user.User{
		UserID:       uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		Role:         role,
	}
	if err := s.repo.CreateUser(ctx, u); err != nil {
		return "", err
	}
	return s.issueToken(u.UserID.String(), u.Role)
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	if email == "" || password == "" {
		return "", ErrInvalidCredentials
	}

	u, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	ok, err := verifyPassword(password, u.PasswordHash)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", ErrInvalidCredentials
	}

	return s.issueToken(u.UserID.String(), u.Role)
}

func (s *Service) ValidateToken(ctx context.Context, token string) (string, string, error) {
	_ = ctx

	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return "", "", ErrInvalidToken
	}

	payloadB64 := parts[0]
	sigB64 := parts[1]

	expectedSig := sign(payloadB64, s.tokenSecret)
	if !hmac.Equal([]byte(expectedSig), []byte(sigB64)) {
		return "", "", ErrInvalidToken
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	var payload tokenPayload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return "", "", ErrInvalidToken
	}

	if payload.Exp < time.Now().Unix() {
		return "", "", ErrTokenExpired
	}
	if payload.Sub == "" || payload.Role == "" {
		return "", "", ErrInvalidToken
	}

	return payload.Sub, payload.Role, nil
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
