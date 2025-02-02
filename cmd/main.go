package main

import (
	"Test-Task-Go/internal/api"
	"Test-Task-Go/internal/repository"
	"Test-Task-Go/logger"
	"Test-Task-Go/pkg/config"
	"Test-Task-Go/pkg/db"
	"Test-Task-Go/pkg/db/cache"
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.LoadConfig()

	logger.InitLogger()

	db.InitRedis(cfg.RedisAddr)
	db.InitMySql(cfg.MySQLDSN)

	sourceRepo := repository.NewSourceRepository()
	campaignRepo := repository.NewCampaignRepository()
	creativeRepo := repository.NewCreativeRepository()

	api.InitRoutes()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("Starting cache updates...")

		cache.UpdateCacheForSources(sourceRepo)
		cache.UpdateCacheForCampaigns(sourceRepo, campaignRepo)
		cache.UpdateCacheForCreatives(sourceRepo, creativeRepo)
	}()

	go func() {
		log.Println("Starting HTTP server on port 8080...")
		fmt.Println("Starting HTTP server on port 8080...")
		if err := api.Server(ctx, nil); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")
}
