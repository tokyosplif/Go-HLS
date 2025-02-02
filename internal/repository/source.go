package repository

import (
	"Test-Task-Go/internal/db"
	"Test-Task-Go/internal/entity"
	"context"
	"database/sql"
	"log"
)

type SourceRepository interface {
	GetActiveSources(ctx context.Context) ([]entity.Source, error)
}

type sourceRepo struct{}

func NewSourceRepository() SourceRepository {
	return &sourceRepo{}
}

func (r *sourceRepo) GetActiveSources(ctx context.Context) ([]entity.Source, error) {
	query := "SELECT id, name, status FROM sources WHERE status = 'active'"

	rows, err := db.Mdb.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}(rows)

	var sources []entity.Source
	for rows.Next() {
		var source entity.Source
		if err := rows.Scan(&source.ID, &source.Name, &source.Status); err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}
	return sources, nil
}
