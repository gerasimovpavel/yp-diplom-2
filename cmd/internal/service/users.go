package service

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"time"
	"yp-diplom-2/cmd/internal/auth"
	"yp-diplom-2/cmd/internal/domain"
	"yp-diplom-2/cmd/internal/hash"
	"yp-diplom-2/cmd/internal/repository"
	"yp-diplom-2/cmd/internal/uuid7"
)

type UserService struct {
	repo            repository.Users
	tokenManager    auth.TokenManager
	hasher          hash.PasswordHasher
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewUserService(repo repository.Users, tokenManager auth.TokenManager, hasher hash.PasswordHasher, accessTTL, refreshTTL time.Duration) *UserService {
	return &UserService{
		repo:         repo,
		tokenManager: tokenManager,
		hasher:       hasher,
	}
}

func (s *UserService) SignUp(ctx context.Context, input SignUpInput) error {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}
	gofakeit.Seed(0)
	verificationCode := gofakeit.LetterN(32)

	uuid, err := uuid7.New()
	if err != nil {
		return err
	}
	user := domain.User{
		UserID:       uuid,
		Email:        input.Email,
		LastName:     input.LastName,
		FirstName:    input.FirstName,
		Password:     passwordHash,
		RegisteredAt: time.Now(),
		LastVisitAt:  time.Now(),
		Verification: domain.Verification{
			Code: verificationCode,
		},
	}

	if err := s.repo.Create(ctx, user); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return err
		}
		return err
	}

	return nil
}

func (s *UserService) SignIn(ctx context.Context, input SignInInput) (Tokens, error) {
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return Tokens{}, err
	}

	user, err := s.repo.GetByCredentials(ctx, input.Email, passwordHash)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return Tokens{}, err
		}
		return Tokens{}, err
	}
	return s.CreateSession(ctx, user.UserID)
}

func (s *UserService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	user, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return s.CreateSession(ctx, user.UserID)
}

func (s *UserService) Verify(ctx context.Context, userID uuid.UUID, hash string) error {
	err := s.repo.Verify(ctx, userID, hash)
	if err != nil {
		if errors.Is(err, domain.ErrVerificationCodeInvalid) {
			return err
		}
		return err
	}

	return nil
}

func (s *UserService) CreateSession(ctx context.Context, userId uuid.UUID) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userId.String(), s.AccessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken = s.tokenManager.NewRefreshToken()

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.RefreshTokenTTL),
	}

	err = s.repo.SetSession(ctx, userId, session)

	return res, err
}
