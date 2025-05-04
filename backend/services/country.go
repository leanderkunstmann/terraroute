package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/leanderkunstmann/terraroute/backend/models"
	"github.com/uptrace/bun"
)

var ErrCountryNotFound = fmt.Errorf("country not found")
var ErrCountriesNotFound = fmt.Errorf("countries not found")

type CountryService struct {
	db *bun.DB
}

func NewCountryService(db *bun.DB) *CountryService {
	return &CountryService{db: db}
}

func (as *CountryService) ListCountries(ctx context.Context, continent string) ([]models.Country, error) {
	var countries []models.Country
	query := as.db.NewSelect().Model(&countries)

	if continent != "" {
		query.Where("continent = ?", continent)
	}

	if err := query.Scan(ctx); err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to list countries: %w", err)
	}

	if len(countries) == 0 {
		return nil, ErrCountriesNotFound
	}

	return countries, nil
}

func (as *CountryService) ListCountry(ctx context.Context, code string) (models.FullCountry, error) {
	var country models.Country
	query := as.db.NewSelect().Model(&country)

	//placeholder borders
	borders := [][]models.Coordinate{
		{
			{Lat: 49.0, Lng: -114.0},
			{Lat: 49.0, Lng: -80.0},
			{Lat: 83.0, Lng: -80.0},
			{Lat: 83.0, Lng: -141.0},
			{Lat: 49.0, Lng: -141.0},
		},
		{
			{Lat: 46.0, Lng: -55.0},
			{Lat: 51.0, Lng: -55.0},
			{Lat: 51.0, Lng: -60.0},
			{Lat: 46.0, Lng: -60.0},
			{Lat: 46.0, Lng: -55.0},
		},
	}

	if code != "" {
		query.Where("code = ?", code)
	}

	if err := query.Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.FullCountry{}, ErrCountryNotFound
		}
		return models.FullCountry{}, fmt.Errorf("failed to list countries: %w", err)
	}

	return models.FullCountry{Country: country, Borders: borders}, nil
}
