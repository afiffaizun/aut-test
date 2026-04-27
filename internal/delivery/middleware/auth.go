package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

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

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c.Writer, domain.ErrInvalidToken)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(c.Writer, domain.ErrInvalidToken)
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := m.jwtManager.ValidateAccessToken(token)
		if err != nil {
			response.Unauthorized(c.Writer, err)
			c.Abort()
			return
		}

		isBlacklisted, err := m.blacklist.IsBlacklisted(c.Request.Context(), claims.TokenID)
		if err != nil {
			response.InternalServerError(c.Writer, err)
			c.Abort()
			return
		}

		if isBlacklisted {
			response.Unauthorized(c.Writer, domain.ErrTokenRevoked)
			c.Abort()
			return
		}

		c.Header("X-User-ID", claims.UserID)
		c.Header("X-Token-ID", claims.TokenID)

		c.Next()
	}
}