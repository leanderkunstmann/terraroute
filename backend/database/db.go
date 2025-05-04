package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/leanderkunstmann/terraroute/backend/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// db is the global database connection
var db *bun.DB

type Config struct {
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
	Host     string `json:"host" mapstructure:"host"`
	Port     string `json:"port" mapstructure:"port"`
	Name     string `json:"name" mapstructure:"name"`
	LocalDB  bool   `json:"local" mapstructure:"local"`
}

func (cfg *Config) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name)
}

// New creates a new database connection based on the provided configuration.
// If LocalDB is true, it creates an in-memory SQLite database.
// Otherwise, it creates a PostgreSQL database connection using the provided credentials.
func New(ctx context.Context, cfg *Config) (*bun.DB, error) {
	if cfg.LocalDB {
		return newLocalDB(ctx)
	}

	return newPostgresDB(ctx, cfg)
}

func newPostgresDB(ctx context.Context, cfg *Config) (*bun.DB, error) {
	dsn := cfg.GetDSN()
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db = bun.NewDB(sqldb, pgdialect.New())

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil
}

func newLocalDB(ctx context.Context) (*bun.DB, error) {
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
		{IATA: "PVG", Name: "Shanghai Pudong International Airport", City: "Shanghai", Country: "China", Continent: "Asia", Latitude: 31.1434, Longitude: 121.805},
		{IATA: "NGO", Name: "Chubu Centrair International Airport", City: "Nagoya", Country: "Japan", Continent: "Asia", Latitude: 34.8583, Longitude: 136.805},
		{IATA: "AKL", Name: "Auckland Airport", City: "Auckland", Country: "New Zealand", Continent: "Oceania", Latitude: -37.0081, Longitude: 174.792},
		{IATA: "ADD", Name: "Addis Ababa Bole International Airport", City: "Addis Ababa", Country: "Ethiopia", Continent: "Africa", Latitude: 8.97789, Longitude: 38.799301},
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

	if err = initLocalDB(ctx, airports, aircrafts, flights); err != nil {
		return nil, fmt.Errorf("creating database: %w", err)
	}

	return db, nil
}

func initLocalDB(ctx context.Context, airports []models.Airport, aircrafts []models.Aircraft, flights []models.Flight) error {
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
