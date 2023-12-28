package service

import (
	"time"
	"yp-diplom-2/cmd/internal/auth"
	"yp-diplom-2/cmd/internal/hash"
	"yp-diplom-2/cmd/internal/repository"
)

type Dependencies struct {
	Repos           *repository.Repositories
	TokenManager    auth.TokenManager
	Hasher          hash.PasswordHasher
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type Services struct {
	Users Users
}

func NewServices(deps Dependencies) *Services {
	return &Services{
		Users: NewUserService(deps.Repos.Users, deps.TokenManager, deps.Hasher, deps.AccessTokenTTL, deps.RefreshTokenTTL),
	}
}
