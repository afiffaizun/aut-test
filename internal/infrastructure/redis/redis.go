package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	*redis.Client
}

func Connect(ctx context.Context, addr, password string, db int) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Redis{Client: client}, nil
}

func (r *Redis) Close() error {
	return r.Client.Close()
}

type BlacklistService struct {
	redis *Redis
}

func NewBlacklistService(redis *Redis) *BlacklistService {
	return &BlacklistService{redis: redis}
}

func (s *BlacklistService) AddToBlacklist(ctx context.Context, tokenID string, expiry time.Duration) error {
	return s.redis.Set(ctx, "blacklist:"+tokenID, "1", expiry).Err()
}

func (s *BlacklistService) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	result, err := s.redis.Exists(ctx, "blacklist:"+tokenID).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (s *BlacklistService) RemoveFromBlacklist(ctx context.Context, tokenID string) error {
	return s.redis.Del(ctx, "blacklist:"+tokenID).Err()
}