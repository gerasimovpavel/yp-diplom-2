package service

import (
	"context"
	"github.com/google/uuid"
)

type SignUpInput struct {
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type SignInInput struct {
	Email    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Users interface {
	SignUp(ctx context.Context, input SignUpInput) error
	SignIn(ctx context.Context, input SignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	Verify(ctx context.Context, userID uuid.UUID, hash string) error
}
