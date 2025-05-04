package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/leanderkunstmann/terraroute/backend/models"
	"github.com/uptrace/bun"
)

type AirportService struct {
	db *bun.DB
}

func NewAirportService(db *bun.DB) *AirportService {
	return &AirportService{db: db}
}

func (as *AirportService) ListAirports(ctx context.Context, iata, continent, country string) ([]models.Airport, error) {
	var airports []models.Airport
	query := as.db.NewSelect().Model(&airports)

	if iata != "" {
		query.Where("iata = ?", iata)
	}
	if continent != "" {
		query.Where("continent = ?", continent)
	}
	if country != "" {
		query.Where("country = ?", country)
	}

	if err := query.Scan(ctx); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to list airports: %w", err)
	}

	return airports, nil
}
