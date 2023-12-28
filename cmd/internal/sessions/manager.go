package sessions

import (
	"time"
	"yp-diplom-2/cmd/internal/domain"
)

type Session struct {
	User         domain.User
	AccessToken  string
	RefreshToken struct {
		Key string
		TTL time.Time
	}
}

type SessionManager struct {
	Sessions []*Session
}

func New() *SessionManager {
	return &SessionManager{Sessions: []*Session{}}
}
func (sm *SessionManager) Add() error {

	return nil
}
