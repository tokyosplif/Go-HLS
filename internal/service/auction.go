package service

import (
	"Test-Task-Go/internal/entity"
	"Test-Task-Go/internal/repository"
	"context"
	"fmt"
	"log"
)

type AuctionService interface {
	ProcessAuction(ctx context.Context, sourceID, maxDuration int) ([]entity.Creative, error)
}

type auctionService struct {
	sourceRepo   repository.SourceRepository
	creativeRepo repository.CreativeRepository
}

func NewAuctionService(sourceRepo repository.SourceRepository, creativeRepo repository.CreativeRepository) AuctionService {
	return &auctionService{sourceRepo: sourceRepo, creativeRepo: creativeRepo}
}

func (s *auctionService) ProcessAuction(ctx context.Context, sourceID, maxDuration int) ([]entity.Creative, error) {
	activeSources, err := s.sourceRepo.GetActiveSources(ctx)
	if err != nil {
		log.Printf("Error retrieving active sources: %v", err)
		return nil, fmt.Errorf("unable to retrieve active sources")
	}

	activeSourceMap := make(map[int]struct{})
	for _, source := range activeSources {
		activeSourceMap[source.ID] = struct{}{}
	}

	if _, found := activeSourceMap[sourceID]; !found {
		log.Printf("Source with ID=%d is inactive", sourceID)
		return nil, fmt.Errorf("source with ID %d is inactive", sourceID)
	}

	creatives, err := s.creativeRepo.GetCreativesBySourceID(ctx, sourceID, maxDuration)
	if err != nil {
		log.Printf("Error retrieving creatives for SourceID=%d: %v", sourceID, err)
		return nil, fmt.Errorf("unable to retrieve creatives")
	}

	if len(creatives) == 0 {
		log.Printf("No creatives found for SourceID=%d with CueOutDuration=%d", sourceID, maxDuration)
		return nil, fmt.Errorf("no creatives found for SourceID %d with CueOutDuration %d", sourceID, maxDuration)
	}

	s.creativeRepo.SortCreativesByPrice(creatives)

	return creatives, nil
}
