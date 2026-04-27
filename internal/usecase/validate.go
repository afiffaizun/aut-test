package usecase

import (
	"context"
	"errors"

	"auth-service/internal/domain"
	"auth-service/internal/infrastructure/redis"
	"auth-service/internal/pkg/jwt"
)

type ValidateUseCase struct {
	blacklist  *redis.BlacklistService
	jwtManager *jwt.JWT
}

func NewValidateUseCase(
	blacklist *redis.BlacklistService,
	jwtManager *jwt.JWT,
) *ValidateUseCase {
	return &ValidateUseCase{
		blacklist:  blacklist,
		jwtManager: jwtManager,
	}
}

func (uc *ValidateUseCase) Execute(ctx context.Context, accessToken string) (*domain.TokenClaims, error) {
	claims, err := uc.jwtManager.ValidateAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	isBlacklisted, err := uc.blacklist.IsBlacklisted(ctx, claims.TokenID)
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		return nil, err
	}

	if isBlacklisted {
		return nil, domain.ErrTokenRevoked
	}

	return &domain.TokenClaims{
		UserID:  claims.UserID,
		TokenID: claims.TokenID,
	}, nil
}