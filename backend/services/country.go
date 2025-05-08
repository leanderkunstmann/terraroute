package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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

func (as *CountryService) ListCountry(ctx context.Context, code string) (models.CountryBorders, error) {

	var countryBorders models.CountryBorders
	var countryBordersLocal models.CountryBordersLocal

	var query *bun.SelectQuery

	if as.db.Dialect().Name().String() == "sqlite" {
		query = as.db.NewSelect().Model(&countryBordersLocal)
	} else {
		query = as.db.NewSelect().Model(&countryBorders)
	}

	if code != "" {
		code = strings.ToUpper(code)
		query.Where("code = ?", code)
	}

	if err := query.Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.CountryBorders{}, ErrCountryNotFound
		}
		return models.CountryBorders{}, fmt.Errorf("failed to list countries: %w", err)
	}

	if as.db.Dialect().Name().String() == "sqlite" {

		// Create a variable of the struct type to hold the unmarshaled data
		var localBorders models.GeoJson

		err := json.Unmarshal([]byte(countryBordersLocal.Borders), &localBorders)
		if err != nil {
			// Handle the error if unmarshaling fails (e.g., invalid JSON, type mismatch)
			fmt.Println("Error unmarshaling JSON to struct:", err)
			return models.CountryBorders{}, fmt.Errorf("failed to list countries: %w", err)
		}

		return models.CountryBorders{Code: countryBordersLocal.Code, Borders: localBorders}, nil
	}
	return models.CountryBorders{Code: countryBorders.Code, Borders: countryBorders.Borders}, nil
}
