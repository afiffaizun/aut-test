package usecase

import (
	"context"
	"time"

	"auth-service/internal/domain"
	"auth-service/internal/infrastructure/redis"
	"auth-service/internal/pkg/jwt"
)

type LogoutUseCase struct {
	tokenRepo  domain.TokenRepository
	blacklist  *redis.BlacklistService
	jwtManager *jwt.JWT
}

func NewLogoutUseCase(
	tokenRepo domain.TokenRepository,
	blacklist *redis.BlacklistService,
	jwtManager *jwt.JWT,
) *LogoutUseCase {
	return &LogoutUseCase{
		tokenRepo:  tokenRepo,
		blacklist:  blacklist,
		jwtManager: jwtManager,
	}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, refreshToken string, accessTokenID string, accessTokenExpiry time.Duration) error {
	if refreshToken != "" {
		if err := uc.tokenRepo.RevokeRefreshToken(ctx, refreshToken); err != nil {
			return err
		}
	}

	if accessTokenID != "" && accessTokenExpiry > 0 {
		if err := uc.blacklist.AddToBlacklist(ctx, accessTokenID, accessTokenExpiry); err != nil {
			return err
		}
	}

	return nil
}