package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/leanderkunstmann/terraroute/backend/database"
	"github.com/leanderkunstmann/terraroute/backend/handlers"
	"github.com/leanderkunstmann/terraroute/backend/services"
	"github.com/rs/cors"
)

type Config struct {
	Database    database.Config `mapstructure:"database"`
	AllowedCors []string        `mapstructure:"allowed_cors"`
	ListenAddr  string          `mapstructure:"listen_addr"`
}

type svcs struct {
	Aircraft *services.AircraftService
	Airport  *services.AirportService
	Distance *services.DistanceCalculator
	Country  *services.CountryService
}

const (
	listenAddr        = ":8080"
	readHeaderTimeout = 5 * time.Second
	defaultCORSOrigin = "http://localhost:5173"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	allowedCORSConfig := os.Getenv("ALLOWED_CORS")
	var allowedOrigins []string = []string{defaultCORSOrigin}

	if allowedCORSConfig != "" {
		allowedOrigins = strings.Split(allowedCORSConfig, ",")
		log.Printf("Allowed CORS origins read from environment variable: %v", allowedOrigins)
	} else {
		log.Printf("Allowed CORS origins set to default: %v", allowedOrigins)
	}

	// TODO: read config via viper from both .env and config file
	// optionally use viper with cobra to also read from CLI args
	cfg := Config{
		ListenAddr:  listenAddr,
		AllowedCors: allowedOrigins,
		Database: database.Config{
			LocalDB:  strings.EqualFold(os.Getenv("LOCAL_DB"), "true"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
		},
	}

	db, err := database.New(ctx, &cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	s := svcs{
		Aircraft: services.NewAircraftService(db),
		Airport:  services.NewAirportService(db),
		Distance: services.NewDistanceCalculator(db),
		Country:  services.NewCountryService(db),
	}
	handlers := []handlers.Handler{
		handlers.NewAircraftHandler(s.Aircraft),
		handlers.NewAirportHandler(s.Airport),
		handlers.NewDistanceHandler(s.Distance),
		handlers.NewCountryHandler(s.Country),
	}

	r := mux.NewRouter()

	for _, handler := range handlers {
		handler.Register(r)
	}

	corsOptions := cors.New(cors.Options{
		AllowedOrigins: cfg.AllowedCors,
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		// AllowCredentials: true,
	})

	handler := corsOptions.Handler(r)

	// Use a custom [http.Server] to set a read header timeout
	// to prevent slowloris attacks.
	srv := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	fmt.Println("Server listening on :8080")
	fmt.Println("Available routes:")
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println(pathTemplate)
		}
		return nil
	})
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
