package usecase

import (
	"context"
	"errors"
	"time"

	"auth-service/internal/domain"
	"auth-service/internal/infrastructure/redis"
	"auth-service/internal/pkg/jwt"
)

type RefreshUseCase struct {
	tokenRepo    domain.TokenRepository
	blacklist    *redis.BlacklistService
	jwtManager   *jwt.JWT
}

func NewRefreshUseCase(
	tokenRepo domain.TokenRepository,
	blacklist *redis.BlacklistService,
	jwtManager *jwt.JWT,
) *RefreshUseCase {
	return &RefreshUseCase{
		tokenRepo:  tokenRepo,
		blacklist:  blacklist,
		jwtManager: jwtManager,
	}
}

func (uc *RefreshUseCase) Execute(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	rt, err := uc.tokenRepo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidToken
		}
		return nil, err
	}

	if rt.IsRevoked {
		return nil, domain.ErrTokenRevoked
	}

	if rt.ExpiresAt.Before(time.Now()) {
		return nil, domain.ErrTokenExpired
	}

	if err := uc.tokenRepo.RevokeRefreshToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	accessToken, err := uc.jwtManager.GenerateAccessToken(rt.UserID)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := uc.jwtManager.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	_, err = uc.tokenRepo.StoreRefreshToken(ctx, &domain.RefreshToken{
		UserID:    rt.UserID,
		TokenHash: newRefreshToken,
		ExpiresAt: time.Now().Add(uc.jwtManager.GetRefreshTokenExpiry()),
	})
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}