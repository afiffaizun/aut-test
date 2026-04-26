package main

import (
	"context"
	"log"

	"auth-service/internal/config"
	"auth-service/internal/infrastructure/database"
	"auth-service/internal/infrastructure/redis"
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
	db.Close()
	log.Println("Database connected successfully!")

	rdb, err := redis.Connect(context.Background(), cfg.Redis.Addr(), cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	rdb.Close()
	log.Println("Redis connected successfully!")

	log.Printf("Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Database: %s:%s/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	log.Printf("Redis: %s", cfg.Redis.Addr())
}