package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/leanderkunstmann/terraroute/backend/models"
	"github.com/uptrace/bun"
)

type AircraftService struct {
	db *bun.DB
}

func NewAircraftService(db *bun.DB) *AircraftService {
	return &AircraftService{db: db}
}

func (as *AircraftService) ListAircraft(ctx context.Context, manufacturer, aircraftType string) ([]models.Aircraft, error) {
	var aircrafts []models.Aircraft
	query := as.db.NewSelect().Model(&aircrafts)

	if manufacturer != "" {
		query.Where("manufacturer = ?", manufacturer)
	}
	if aircraftType != "" {
		query.Where("type = ?", aircraftType)
	}

	if err := query.Scan(ctx); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to list aircrafts: %w", err)
	}

	return aircrafts, nil
}
