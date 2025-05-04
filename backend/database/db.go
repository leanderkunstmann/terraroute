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
	LocalDB  bool   `mapstructure:"local_db"`
	Username string `mapstructure:"db_username"`
	Password string `mapstructure:"db_password"`
	Host     string `mapstructure:"db_host"`
	Port     string `mapstructure:"db_port"`
	Name     string `mapstructure:"db_name"`
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
		{IATA: "CPT", Name: "Cape Town International Airport", City: "Cape Town", Country: "South Africa", Continent: "Africa", Latitude: -33.965, Longitude: 18.602},
		{IATA: "SCL", Name: "Arturo Merino Benítez International Airport", City: "Santiago de Chile", Country: "Chile", Continent: "South America", Latitude: -33.393056, Longitude: -70.785833},
		{IATA: "USH", Name: "Ushuaia Malvinas Argentinas International Airport", City: "Ushuaia", Country: "Argentina", Continent: "South America", Latitude: -54.8433, Longitude: -68.2944},
		{IATA: "GUA", Name: "La Aurora International Airport", City: "Guatemala City", Country: "Guatemala", Continent: "North America", Latitude: 14.5817, Longitude: -90.5267},
		{IATA: "ANC", Name: "Ted Stevens Anchorage International Airport", City: "Anchorage", Country: "USA", Continent: "North America", Latitude: 61.1744444, Longitude: -149.99639},
		{IATA: "MNL", Name: "Ninoy Aquino International Airport", City: "Manila", Country: "Philippines", Continent: "Asia", Latitude: 14.5086, Longitude: 121.0199966},
		{IATA: "CTS", Name: "New Chitose Airport", City: "Sapporo", Country: "Japan", Continent: "Asia", Latitude: 42.7753, Longitude: 141.692},
		{IATA: "PER", Name: "Perth Airport", City: "Perth", Country: "Australia", Continent: "Oceania", Latitude: -31.9403, Longitude: 115.967},
		{IATA: "LAG", Name: "Laughtale Grandline International Airport", City: "One Piece", Country: "Bermuda", Continent: "North America", Latitude: 25.0, Longitude: -71.0}, // Fictional airport
	}

	aircrafts := []models.Aircraft{
		{Id: 1, Type: models.Commercial, Name: "Boeing 737", Manufacturer: "Boeing", Range: 3510},
		{Id: 2, Type: models.Heavy, Name: "Gulfstream G650", Manufacturer: "Gulfstream", Range: 7500},
		{Id: 3, Type: models.Cargo, Name: "Antonov An-225", Manufacturer: "Antonov", Range: 9700},
		{Id: 4, Type: models.Commercial, Name: "Airbus A320", Manufacturer: "Airbus", Range: 3200},
	}

	countries := []models.Country{
		{Code: "us", Name: "United States of America", Continent: "North America"},
		{Code: "ger", Name: "Germany", Continent: "Europe"},
		{Code: "ru", Name: "Russia", Continent: "Europe"},
		{Code: "cn", Name: "China", Continent: "Asia"},
	}

	flights := []models.Flight{
		{FlightNumber: "AA100", AircraftId: 1, Origin: "JFK", Destination: "LAX", DepartureTime: "2023-10-01T08:00:00Z", ArrivalTime: "2023-10-01T11:00:00Z"},
		{FlightNumber: "AF200", AircraftId: 2, Origin: "CDG", Destination: "JFK", DepartureTime: "2023-10-02T09:00:00Z", ArrivalTime: "2023-10-02T12:00:00Z"},
		{FlightNumber: "DL300", AircraftId: 3, Origin: "LAX", Destination: "CDG", DepartureTime: "2023-10-03T10:00:00Z", ArrivalTime: "2023-10-03T18:00:00Z"},
	}

	if err = initLocalDB(ctx, airports, aircrafts, countries, flights); err != nil {
		return nil, fmt.Errorf("creating database: %w", err)
	}

	return db, nil
}

func initLocalDB(ctx context.Context, airports []models.Airport, aircrafts []models.Aircraft, countries []models.Country, flights []models.Flight) error {
	_, err := db.NewCreateTable().Model((*models.Airport)(nil)).Exec(ctx)
	if err != nil {
		return err
	}
	_, err = db.NewCreateTable().Model((*models.Aircraft)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewCreateTable().Model((*models.Country)(nil)).Exec(ctx)
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

	_, err = db.NewInsert().Model(&countries).Exec(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewInsert().Model(&flights).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
