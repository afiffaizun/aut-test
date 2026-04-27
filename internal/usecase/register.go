package usecase

import (
	"context"
	"errors"
	"time"

	"auth-service/internal/domain"
	"auth-service/internal/pkg/hashing"
	"auth-service/internal/pkg/jwt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrWeakPassword      = errors.New("password does not meet security requirements")
)

type RegisterUseCase struct {
	userRepo     domain.UserRepository
	tokenRepo    domain.TokenRepository
	hasher       *hashing.Hasher
	jwtManager   *jwt.JWT
}

func NewRegisterUseCase(
	userRepo domain.UserRepository,
	tokenRepo domain.TokenRepository,
	hasher *hashing.Hasher,
	jwtManager *jwt.JWT,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		hasher:     hasher,
		jwtManager: jwtManager,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, email, password string) (*domain.User, *domain.TokenPair, error) {
	existingUser, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, nil, err
	}
	if existingUser != nil {
		return nil, nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := uc.hasher.HashPassword(password)
	if err != nil {
		return nil, nil, err
	}

	user := &domain.User{
		Email:    email,
		Password: hashedPassword,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, nil, err
	}

	accessToken, err := uc.jwtManager.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := uc.jwtManager.GenerateRefreshToken()
	if err != nil {
		return nil, nil, err
	}

	_, err = uc.tokenRepo.StoreRefreshToken(ctx, &domain.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshToken,
		ExpiresAt: time.Now().Add(uc.jwtManager.GetRefreshTokenExpiry()),
	})
	if err != nil {
		return nil, nil, err
	}

	return user, &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}