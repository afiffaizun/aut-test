package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"auth-service/internal/config"
	"auth-service/internal/delivery/handler"
	"auth-service/internal/delivery/middleware"
	"auth-service/internal/domain"
	"auth-service/internal/infrastructure/database"
	"auth-service/internal/infrastructure/redis"
	"auth-service/internal/pkg/hashing"
	"auth-service/internal/pkg/jwt"
	"auth-service/internal/usecase"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(context.Background(), cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()
	log.Println("Database connected successfully!")

	rdb, err := redis.Connect(context.Background(), cfg.Redis.Addr(), cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	defer rdb.Close()
	log.Println("Redis connected successfully!")

	hasher := hashing.NewHasher(hashing.Params{
		Memory:      cfg.Argon2.Memory,
		Iterations:  cfg.Argon2.Iterations,
		Parallelism: cfg.Argon2.Parallelism,
		SaltLength:  cfg.Argon2.SaltLength,
		KeyLength:   cfg.Argon2.KeyLength,
	})

	jwtManager := jwt.NewJWT(cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry, cfg.JWT.RefreshTokenExpiry)

	blacklist := redis.NewBlacklistService(rdb)

	registerUC := usecase.NewRegisterUseCase(nil, nil, hasher, jwtManager)
	loginUC := usecase.NewLoginUseCase(nil, nil, hasher, jwtManager)
	refreshUC := usecase.NewRefreshUseCase(nil, blacklist, jwtManager)
	logoutUC := usecase.NewLogoutUseCase(nil, blacklist, jwtManager)
	validateUC := usecase.NewValidateUseCase(blacklist, jwtManager)

	authHandler := handler.NewAuthHandler(registerUC, loginUC, refreshUC, logoutUC, validateUC)

	authMiddleware := middleware.NewAuthMiddleware(blacklist, jwtManager)
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit.Requests, cfg.RateLimit.Window)

	r := gin.Default()

	r.Use(middleware.SecureHeaders())
	r.Use(middleware.Logger())

	api := r.Group("/api/v1/auth")
	{
		api.POST("/register", rateLimiter.Limit(), authHandler.Register)
		api.POST("/login", rateLimiter.Limit(), authHandler.Login)
		api.POST("/refresh", authHandler.Refresh)
		api.POST("/logout", authMiddleware.Authenticate(), authHandler.Logout)
		api.GET("/me", authMiddleware.Authenticate(), authHandler.Me)
		api.GET("/validate", authMiddleware.Authenticate(), authHandler.Validate)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Printf("Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := r.Run(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

type UserRepository struct{}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, domain.ErrUserNotFound
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return nil, domain.ErrUserNotFound
}

type TokenRepository struct{}

func (r *TokenRepository) StoreRefreshToken(ctx context.Context, token *domain.RefreshToken) (string, error) {
	return token.TokenHash, nil
}

func (r *TokenRepository) GetRefreshToken(ctx context.Context, tokenID string) (*domain.RefreshToken, error) {
	return nil, domain.ErrUserNotFound
}

func (r *TokenRepository) RevokeRefreshToken(ctx context.Context, tokenID string) error {
	return nil
}

func (r *TokenRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	return nil
}