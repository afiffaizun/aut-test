package main

import (
	"context"
	"log"

	"auth-service/internal/config"
	"auth-service/internal/infrastructure/database"
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

	log.Printf("Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Database: %s:%s/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	log.Printf("Redis: %s", cfg.Redis.Addr())
	log.Println("Database connected successfully!")
}