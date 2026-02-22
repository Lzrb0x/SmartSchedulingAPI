package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Lzrb0x/SmartSchedulingAPI/internal/config"
	"github.com/Lzrb0x/SmartSchedulingAPI/internal/database"
	"github.com/Lzrb0x/SmartSchedulingAPI/internal/server"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.Connect(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	srv := server.New(cfg, db)

	if err := srv.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
