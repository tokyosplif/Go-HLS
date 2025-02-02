package repository

import (
	"Test-Task-Go/internal/db"
	"Test-Task-Go/internal/entity"
	"context"
	"database/sql"
	"fmt"
	"log"
	"sort"
)

type CreativeRepository interface {
	GetCreativesBySourceID(ctx context.Context, sourceID int, maxDuration int) ([]entity.Creative, error)
	SortCreativesByPrice(creatives []entity.Creative)
}

type creativeRepo struct{}

func NewCreativeRepository() CreativeRepository {
	return &creativeRepo{}
}

func (r *creativeRepo) GetCreativesBySourceID(ctx context.Context, sourceID int, maxDuration int) ([]entity.Creative, error) {
	query := `
        SELECT id, campaign_id, price, duration, playlist_hls
        FROM creatives
        WHERE campaign_id IN (
            SELECT id FROM campaigns
            WHERE source_id = ? AND NOW() BETWEEN start_time AND end_time
        ) AND duration <= ? 
        ORDER BY duration DESC`

	rows, err := db.Mdb.QueryContext(ctx, query, sourceID, maxDuration)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %w", err)
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}(rows)

	var creatives []entity.Creative
	for rows.Next() {
		var creative entity.Creative
		if err := rows.Scan(&creative.ID, &creative.CampaignID, &creative.Price, &creative.Duration, &creative.PlaylistHLS); err != nil {
			return nil, fmt.Errorf("data scanning error: %w", err)
		}
		creatives = append(creatives, creative)
	}

	return creatives, nil
}

func (r *creativeRepo) SortCreativesByPrice(creatives []entity.Creative) {
	sort.SliceStable(creatives, func(i, j int) bool {
		return creatives[i].Price < creatives[j].Price
	})
}
