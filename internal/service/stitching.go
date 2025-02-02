package service

import (
	"Test-Task-Go/internal/entity"
	"Test-Task-Go/internal/playlist"
	"Test-Task-Go/internal/repository"
	"context"
	"fmt"
	"log"
	"math/rand"
)

type StitchingService interface {
	ProcessStitching(ctx context.Context, sourceID int, playlistContent string) (string, error)
}

type stitchingService struct {
	sourceRepo   repository.SourceRepository
	creativeRepo repository.CreativeRepository
}

func NewStitchingService(sourceRepo repository.SourceRepository, creativeRepo repository.CreativeRepository) StitchingService {
	return &stitchingService{sourceRepo: sourceRepo, creativeRepo: creativeRepo}
}

func (s *stitchingService) ProcessStitching(ctx context.Context, sourceID int, playlistContent string) (string, error) {
	activeSources, err := s.sourceRepo.GetActiveSources(ctx)
	if err != nil {
		log.Printf("Error retrieving active sources: %v", err)
		return "", fmt.Errorf("unable to retrieve active sources")
	}

	activeSourceMap := make(map[int]struct{})
	for _, source := range activeSources {
		activeSourceMap[source.ID] = struct{}{}
	}

	if _, found := activeSourceMap[sourceID]; !found {
		log.Printf("Source with ID=%d is inactive", sourceID)
		return "", fmt.Errorf("source with ID %d is inactive", sourceID)
	}

	cueOutDuration, err := playlist.GetCueOutDuration(playlistContent)
	if err != nil || cueOutDuration == 0 {
		if err != nil {
			log.Printf("Error extracting cue-out duration: %v", err)
			return "", fmt.Errorf("unable to extract cue-out duration")
		}
		log.Printf("No valid cue-out duration found in playlist")
		return "", fmt.Errorf("no valid cue-out duration found")
	}

	creatives, err := s.creativeRepo.GetCreativesBySourceID(ctx, sourceID, cueOutDuration)
	if err != nil {
		log.Printf("Error retrieving creatives for SourceID=%d: %v", sourceID, err)
		return "", fmt.Errorf("unable to retrieve creatives")
	}

	if len(creatives) == 0 {
		log.Printf("No creatives found for SourceID=%d with CueOutDuration=%d", sourceID, cueOutDuration)
		return "", fmt.Errorf("no creatives found for SourceID %d with CueOutDuration %d", sourceID, cueOutDuration)
	}

	var selectedCreative entity.Creative
	maxDurationClosest := -1
	var closestCreatives []entity.Creative

	for _, creative := range creatives {
		if maxDurationClosest == -1 || creative.Duration > maxDurationClosest {
			maxDurationClosest = creative.Duration
			closestCreatives = []entity.Creative{creative}
		} else if creative.Duration == maxDurationClosest {
			closestCreatives = append(closestCreatives, creative)
		}
	}

	if len(closestCreatives) > 1 {
		selectedCreative = closestCreatives[rand.Intn(len(closestCreatives))]
	} else {
		selectedCreative = closestCreatives[0]
	}

	modifiedPlaylist, err := playlist.InsertAdsIntoPlaylist(playlistContent, selectedCreative)
	if err != nil {
		log.Printf("Error inserting ads into playlist: %v", err)
		return "", fmt.Errorf("unable to insert ads into playlist")
	}

	return modifiedPlaylist, nil
}
