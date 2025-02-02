package api

import (
	"Test-Task-Go/internal/handler"
	"Test-Task-Go/internal/repository"
	"Test-Task-Go/internal/service"
	"net/http"
)

func InitRoutes() {
	sourceRepo := repository.NewSourceRepository()
	creativeRepo := repository.NewCreativeRepository()

	auctionService := service.NewAuctionService(sourceRepo, creativeRepo)
	stitchingService := service.NewStitchingService(sourceRepo, creativeRepo)

	http.HandleFunc("/auction", handler.NewAuctionHandler(auctionService).HandleAuction)
	http.HandleFunc("/stitching.m3u8", handler.NewStitchingHandler(stitchingService).HandleStitching)
}
