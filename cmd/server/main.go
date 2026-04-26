package main

import (
	"log"

	"auth-service/internal/config"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Database: %s:%s/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	log.Printf("Redis: %s", cfg.Redis.Addr())
}
