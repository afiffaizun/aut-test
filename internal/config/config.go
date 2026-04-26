package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	JWT       JWTConfig
	Argon2    Argon2Config
	RateLimit RateLimitConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret             string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

type Argon2Config struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

type RateLimitConfig struct {
	Enabled  bool
	Requests int
	Window   time.Duration
}

func Load(configPath string) (*Config, error) {
	if err := godotenv.Load(configPath); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	getOrDefault := func(key, defaultValue string) string {
		if val := os.Getenv(key); val != "" {
			return val
		}
		return defaultValue
	}

	getOrDefaultDuration := func(key, defaultValue string) time.Duration {
		if val := os.Getenv(key); val != "" {
			d, err := time.ParseDuration(val)
			if err == nil {
				return d
			}
		}
		d, _ := time.ParseDuration(defaultValue)
		return d
	}

	getOrDefaultBool := func(key string, defaultValue bool) bool {
		if val := os.Getenv(key); val != "" {
			return val == "true"
		}
		return defaultValue
	}

	getOrDefaultInt := func(key string, defaultValue int) int {
		if val := os.Getenv(key); val != "" {
			var n int
			fmt.Sscanf(val, "%d", &n)
			return n
		}
		return defaultValue
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getOrDefault("SERVER_PORT", "8080"),
			Host: getOrDefault("SERVER_HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getOrDefault("DB_HOST", "localhost"),
			Port:     getOrDefault("DB_PORT", "5432"),
			User:     getOrDefault("DB_USER", "postgres"),
			Password: getOrDefault("DB_PASSWORD", "postgres"),
			Name:     getOrDefault("DB_NAME", "auth_service"),
			SSLMode:  getOrDefault("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getOrDefault("REDIS_HOST", "localhost"),
			Port:     getOrDefault("REDIS_PORT", "6379"),
			Password: getOrDefault("REDIS_PASSWORD", ""),
			DB:       getOrDefaultInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:             getOrDefault("JWT_SECRET", "change-this-secret-in-production"),
			AccessTokenExpiry:  getOrDefaultDuration("JWT_ACCESS_TOKEN_EXPIRY", "15m"),
			RefreshTokenExpiry: getOrDefaultDuration("JWT_REFRESH_TOKEN_EXPIRY", "168h"),
		},
		Argon2: Argon2Config{
			Memory:      uint32(getOrDefaultInt("ARGON2_MEMORY", 65536)),
			Iterations:  uint32(getOrDefaultInt("ARGON2_ITERATIONS", 3)),
			Parallelism: uint8(getOrDefaultInt("ARGON2_PARALLELISM", 4)),
			SaltLength:  uint32(getOrDefaultInt("ARGON2_SALTLENGTH", 16)),
			KeyLength:   uint32(getOrDefaultInt("ARGON2_KEYLENGTH", 32)),
		},
		RateLimit: RateLimitConfig{
			Enabled:  getOrDefaultBool("RATE_LIMIT_ENABLED", true),
			Requests: getOrDefaultInt("RATE_LIMIT_REQUESTS", 10),
			Window:   getOrDefaultDuration("RATE_LIMIT_WINDOW", "60s"),
		},
	}

	return cfg, nil
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}