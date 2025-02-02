package repository

import (
	"Test-Task-Go/internal/db"
	"Test-Task-Go/internal/entity"
	"context"
	"database/sql"
	"log"
	"time"
)

type CampaignRepository interface {
	GetActiveCampaigns(ctx context.Context, sourceID int) ([]entity.Campaign, error)
}

type campaignRepo struct{}

func NewCampaignRepository() CampaignRepository {
	return &campaignRepo{}
}

func (r *campaignRepo) GetActiveCampaigns(ctx context.Context, sourceID int) ([]entity.Campaign, error) {
	query := `
        SELECT c.id, c.name, c.start_time, c.end_time
        FROM campaigns c
        JOIN source_campaigns sc ON c.id = sc.campaign_id
        WHERE sc.source_id = ? AND c.start_time <= ? AND c.end_time >= ?`

	now := time.Now()
	rows, err := db.Mdb.QueryContext(ctx, query, sourceID, now, now)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}(rows)

	var campaigns []entity.Campaign
	for rows.Next() {
		var campaign entity.Campaign
		if err := rows.Scan(&campaign.ID, &campaign.Name, &campaign.StartTime, &campaign.EndTime); err != nil {
			return nil, err
		}
		campaigns = append(campaigns, campaign)
	}
	return campaigns, nil
}
