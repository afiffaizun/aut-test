package middleware

import (
	"net/http"
	"strings"

	"auth-service/internal/domain"
	"auth-service/internal/infrastructure/redis"
	"auth-service/internal/pkg/jwt"
	"auth-service/internal/pkg/response"
)

type AuthMiddleware struct {
	blacklist  *redis.BlacklistService
	jwtManager *jwt.JWT
}

func NewAuthMiddleware(blacklist *redis.BlacklistService, jwtManager *jwt.JWT) *AuthMiddleware {
	return &AuthMiddleware{
		blacklist:  blacklist,
		jwtManager: jwtManager,
	}
}

func (m *AuthMiddleware) Authenticate() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Unauthorized(w, domain.ErrInvalidToken)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				response.Unauthorized(w, domain.ErrInvalidToken)
				return
			}

			token := parts[1]
			claims, err := m.jwtManager.ValidateAccessToken(token)
			if err != nil {
				response.Unauthorized(w, err)
				return
			}

			isBlacklisted, err := m.blacklist.IsBlacklisted(r.Context(), claims.TokenID)
			if err != nil {
				response.InternalServerError(w, err)
				return
			}

			if isBlacklisted {
				response.Unauthorized(w, domain.ErrTokenRevoked)
				return
			}

			r.Header.Set("X-User-ID", claims.UserID)
			r.Header.Set("X-Token-ID", claims.TokenID)

			next.ServeHTTP(w, r)
		})
	}
}