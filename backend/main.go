package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/leanderkunstmann/terraroute/backend/config"
	"github.com/leanderkunstmann/terraroute/backend/database"
	"github.com/leanderkunstmann/terraroute/backend/handlers"
	"github.com/leanderkunstmann/terraroute/backend/services"
	"github.com/rs/cors"
)

const readHeaderTimeout = 5 * time.Second

// configPath is the path to the config file.
// It is set by the -config flag.
var configPath string

type svcs struct {
	Aircraft *services.AircraftService
	Airport  *services.AirportService
	Distance *services.DistanceCalculator
}

func main() {
	flag.StringVar(&configPath, "config", "./config.json", "Path to the config file")
	flag.Parse()

	ctx := context.Background()
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatal("Error loading config file:", err)
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

	for _, handler := range handlers {
		handler.Register(r)
	}

	corsOptions := cors.New(cors.Options{
		AllowedOrigins: cfg.AllowedOrigins,
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		// AllowCredentials: true,
	})

	// Use a custom [http.Server] to set a read header timeout
	// to prevent slowloris attacks.
	srv := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           corsOptions.Handler(r),
		ReadHeaderTimeout: readHeaderTimeout,
	}

	fmt.Printf("Server listening on %s\n", cfg.ListenAddr)
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
