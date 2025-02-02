package cache

import (
	"Test-Task-Go/internal/db"
	"Test-Task-Go/internal/entity"
	"Test-Task-Go/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

var ctx = context.Background()

func logExecutionTime(start time.Time, operation string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", operation, elapsed)
}

func GetFromCache(key string, result interface{}) error {
	start := time.Now()
	defer logExecutionTime(start, fmt.Sprintf("Getting data from cache for key '%s'", key))

	val, err := db.Rdb.Get(ctx, key).Result()
	if err != nil {
		log.Printf("Error getting data from cache for key '%s': %v", key, err)
		return err
	}
	log.Printf("Data successfully retrieved from cache for key: %s", key)

	err = json.Unmarshal([]byte(val), result)
	if err != nil {
		log.Printf("Error deserializing data for key '%s': %v", key, err)
		return fmt.Errorf("cache deserialization error: %w", err)
	}
	return nil
}

func SetInCache(key string, value interface{}, ttl time.Duration) error {
	start := time.Now()
	defer logExecutionTime(start, fmt.Sprintf("Saving data to cache for key '%s'", key))

	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("Error serializing data for key '%s': %v", key, err)
		return fmt.Errorf("serialization error: %w", err)
	}
	err = db.Rdb.Set(ctx, key, data, ttl).Err()
	if err != nil {
		log.Printf("Error saving data to Redis for key '%s': %v", key, err)
		return fmt.Errorf("redis save error: %w", err)
	}

	log.Printf("Data successfully saved to Redis for key '%s' with TTL=%v", key, ttl)
	return nil
}

func acquireLock(lockKey string) bool {
	start := time.Now()
	defer logExecutionTime(start, fmt.Sprintf("Acquiring lock for '%s'", lockKey))

	result, err := db.Rdb.SetNX(ctx, lockKey, "locked", 10*time.Second).Result()
	if err != nil {
		log.Printf("Error acquiring lock for '%s': %v", lockKey, err)
		return false
	}
	return result
}

func releaseLock(lockKey string) {
	start := time.Now()
	defer logExecutionTime(start, fmt.Sprintf("Releasing lock for '%s'", lockKey))

	_, err := db.Rdb.Del(ctx, lockKey).Result()
	if err != nil {
		log.Printf("Error releasing lock for '%s': %v", lockKey, err)
	}
}

func UpdateCacheForSources(sourceRepo repository.SourceRepository) {
	start := time.Now()
	defer logExecutionTime(start, "Updating cache for sources")

	lockKey := "lock:active_sources"
	if !acquireLock(lockKey) {
		log.Println("Lock for sources is already active, skipping update")
		return
	}
	defer releaseLock(lockKey)

	var sources []entity.Source
	err := GetFromCache("active_sources", &sources)
	if err == nil {
		log.Println("Data successfully retrieved from cache for sources")
		return
	}

	sources, err = sourceRepo.GetActiveSources(ctx)
	if err != nil {
		log.Printf("Error getting active sources: %v", err)
		return
	}
	log.Printf("Found %d active sources", len(sources))

	err = SetInCache("active_sources", sources, 10*time.Minute)
	if err != nil {
		log.Printf("Error updating cache for sources: %v", err)
	} else {
		log.Println("Cache for sources successfully updated")
	}
}

func UpdateCacheForCampaigns(sourceRepo repository.SourceRepository, campaignRepo repository.CampaignRepository) {
	start := time.Now()
	defer logExecutionTime(start, "Updating cache for campaigns")

	lockKey := "lock:active_campaigns"
	if !acquireLock(lockKey) {
		log.Println("Lock for campaigns is already active, skipping update")
		return
	}
	defer releaseLock(lockKey)

	var allCampaigns []entity.Campaign
	err := GetFromCache("active_campaigns", &allCampaigns)
	if err == nil {
		log.Println("Data successfully retrieved from cache for campaigns")
		return
	}

	sources, err := sourceRepo.GetActiveSources(ctx)
	if err != nil {
		log.Printf("Error getting active sources: %v", err)
		return
	}

	for _, source := range sources {
		log.Printf("Getting active campaigns for SourceID=%d", source.ID)

		campaigns, err := campaignRepo.GetActiveCampaigns(ctx, source.ID)
		if err != nil {
			log.Printf("Error getting active campaigns for SourceID=%d: %v", source.ID, err)
			continue
		}
		log.Printf("Found %d active campaigns for SourceID=%d", len(campaigns), source.ID)

		allCampaigns = append(allCampaigns, campaigns...)
	}

	err = SetInCache("active_campaigns", allCampaigns, 10*time.Minute)
	if err != nil {
		log.Printf("Error updating cache for campaigns: %v", err)
	} else {
		log.Println("Cache for campaigns successfully updated")
	}
}

func UpdateCacheForCreatives(sourceRepo repository.SourceRepository, creativeRepo repository.CreativeRepository) {
	start := time.Now()
	defer logExecutionTime(start, "Updating cache for creatives")

	lockKey := "lock:active_creatives"
	if !acquireLock(lockKey) {
		log.Println("Lock for creatives is already active, skipping update")
		return
	}
	defer releaseLock(lockKey)

	var allCreatives []entity.Creative
	err := GetFromCache("active_creatives", &allCreatives)
	if err == nil {
		log.Println("Data successfully retrieved from cache for creatives")
		return
	}

	sources, err := sourceRepo.GetActiveSources(ctx)
	if err != nil {
		log.Printf("Error getting active sources: %v", err)
		return
	}

	for _, source := range sources {
		log.Printf("Getting creatives for SourceID=%d", source.ID)

		creatives, err := creativeRepo.GetCreativesBySourceID(ctx, source.ID, 30)
		if err != nil {
			log.Printf("Error getting creatives for SourceID=%d: %v", source.ID, err)
			continue
		}

		allCreatives = append(allCreatives, creatives...)
	}

	err = SetInCache("active_creatives", allCreatives, 10*time.Minute)
	if err != nil {
		log.Printf("Error updating cache for creatives: %v", err)
	} else {
		log.Println("Cache for creatives successfully updated")
	}
}
