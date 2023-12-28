package repository

import (
	"context"
	"github.com/google/uuid"
	"yp-diplom-2/cmd/internal/domain"
)

type Users interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	Verify(ctx context.Context, userID uuid.UUID, code string) error
	SetSession(ctx context.Context, userID uuid.UUID, session domain.Session) error
}
