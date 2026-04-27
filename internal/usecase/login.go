package usecase

import (
	"context"
	"errors"
	"time"

	"auth-service/internal/domain"
	"auth-service/internal/pkg/hashing"
	"auth-service/internal/pkg/jwt"
)

type LoginUseCase struct {
	userRepo    domain.UserRepository
	tokenRepo   domain.TokenRepository
	hasher      *hashing.Hasher
	jwtManager  *jwt.JWT
}

func NewLoginUseCase(
	userRepo domain.UserRepository,
	tokenRepo domain.TokenRepository,
	hasher *hashing.Hasher,
	jwtManager *jwt.JWT,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		hasher:     hasher,
		jwtManager: jwtManager,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, email, password string) (*domain.User, *domain.TokenPair, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, nil, domain.ErrInvalidCredentials
		}
		return nil, nil, err
	}

	if err := uc.hasher.VerifyPassword(user.Password, password); err != nil {
		return nil, nil, domain.ErrInvalidCredentials
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