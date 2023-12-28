package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
	"yp-diplom-2/cmd/internal/config"
)

type TokenManager interface {
	NewJWT(userID string, ttl time.Duration) (string, error)
	Parse(accessToken string) (string, error)
	NewRefreshToken() string
}

type Manager struct {
	cfg *config.Config
}

func NewManager(cfg *config.Config) (*Manager, error) {

	if cfg.Auth.JWT.SigningKey == "" {
		return nil, errors.New("empty signing key")
	}
	return &Manager{
		cfg: cfg,
	}, nil
}
func (m *Manager) NewJWT(userID string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(ttl)},
		Subject:   userID,
	})
	return token.SignedString([]byte(m.cfg.Auth.JWT.SigningKey))
}

func (m *Manager) Parse(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}
		return m.cfg.Auth.JWT.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("error get user claims from tokens")
	}
	return claims["sub"].(string), nil
}

func (m *Manager) NewRefreshToken() string {
	return uuid.New().String()
}
