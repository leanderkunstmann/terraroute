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
)

type Config struct {
	Database   database.Config `mapstructure:"database"`
	ListenAddr string          `mapstructure:"listen_addr"`
}

type svcs struct {
	Aircraft *services.AircraftService
	Airport  *services.AirportService
	Distance *services.DistanceCalculator
}

const (
	listenAddr        = ":8080"
	readHeaderTimeout = 5 * time.Second
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// TODO: read config via viper from both .env and config file
	// optionally use viper with cobra to also read from CLI args
	cfg := Config{
		ListenAddr: listenAddr,
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
	}
	handlers := []handlers.Handler{
		handlers.NewAircraftHandler(s.Aircraft),
		handlers.NewAirportHandler(s.Airport),
		handlers.NewDistanceHandler(s.Distance),
	}

	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))
	for _, handler := range handlers {
		handler.Register(r)
	}

	// Use a custom [http.Server] to set a read header timeout
	// to prevent slowloris attacks.
	srv := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           r,
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
