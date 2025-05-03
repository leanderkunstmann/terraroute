package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/leanderkunstmann/terraroute/backend/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func CreateDB(ctx context.Context, db *bun.DB, airports []models.Airport, aircrafts []models.Aircraft, flights []models.Flight) error {
	_, err := db.NewCreateTable().Model((*models.Airport)(nil)).Exec(ctx)
	if err != nil {
		return err
	}
	_, err = db.NewCreateTable().Model((*models.Aircraft)(nil)).Exec(ctx)
	if err != nil {
		return err
	}
	_, err = db.NewCreateTable().Model((*models.Flight)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewInsert().Model(&airports).Exec(ctx)
	if err != nil {
		return err
	}
	_, err = db.NewInsert().Model(&aircrafts).Exec(ctx)
	if err != nil {
		return err
	}
	_, err = db.NewInsert().Model(&flights).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func InitDB(ctx context.Context, lcl string) (*bun.DB, error) {
	var db *bun.DB

	if lcl == "true" {
		sqldb, err := sql.Open("sqlite3", "file::memory:?cache=shared")
		if err != nil {
			return nil, fmt.Errorf("failed to open SQLite database: %w", err)
		}

		db = bun.NewDB(sqldb, sqlitedialect.New())

		if err := db.PingContext(ctx); err != nil {
			return nil, fmt.Errorf("database ping failed: %w", err)
		}

		airports := []models.Airport{
			{IATA: "JFK", Name: "John F. Kennedy International Airport", City: "New York", Country: "USA", Continent: "North America", Latitude: 40.6413, Longitude: -73.7781},
			{IATA: "LAX", Name: "Los Angeles International Airport", City: "Los Angeles", Country: "USA", Continent: "North America", Latitude: 33.9416, Longitude: -118.4085},
			{IATA: "CDG", Name: "Charles de Gaulle Airport", City: "Paris", Country: "France", Continent: "Europe", Latitude: 49.0097, Longitude: 2.5479},
			{IATA: "FRA", Name: "Frankfurt Airport", City: "Frankfurt", Country: "Germany", Continent: "Europe", Latitude: 50.0333, Longitude: 8.5706},
		}
		aircrafts := []models.Aircraft{
			{Id: 1, Type: models.Commercial, Name: "Boeing 737", Manufacturer: "Boeing", Range: 3510},
			{Id: 2, Type: models.Heavy, Name: "Gulfstream G650", Manufacturer: "Gulfstream", Range: 7500},
			{Id: 3, Type: models.Cargo, Name: "Antonov An-225", Manufacturer: "Antonov", Range: 9700},
			{Id: 4, Type: models.Commercial, Name: "Airbus A320", Manufacturer: "Airbus", Range: 3200},
		}
		flights := []models.Flight{
			{FlightNumber: "AA100", AircraftId: 1, Origin: "JFK", Destination: "LAX", DepartureTime: "2023-10-01T08:00:00Z", ArrivalTime: "2023-10-01T11:00:00Z"},
			{FlightNumber: "AF200", AircraftId: 2, Origin: "CDG", Destination: "JFK", DepartureTime: "2023-10-02T09:00:00Z", ArrivalTime: "2023-10-02T12:00:00Z"},
			{FlightNumber: "DL300", AircraftId: 3, Origin: "LAX", Destination: "CDG", DepartureTime: "2023-10-03T10:00:00Z", ArrivalTime: "2023-10-03T18:00:00Z"},
		}

		err = CreateDB(ctx, db, airports, aircrafts, flights)
		if err != nil {
			return nil, fmt.Errorf("creating database: %w", err)
		}
	} else {
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"))

		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
		db = bun.NewDB(sqldb, pgdialect.New())

		if err := db.PingContext(ctx); err != nil {
			return nil, fmt.Errorf("database ping failed: %w", err)
		}
	}

	return db, nil
}
