package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/leanderkunstmann/terraroute/backend/handlers"
)

// use viper for configuration
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	isLocalDB := os.Getenv("LOCAL_DB")

	// Check if the application is running in test mode
	ctx := context.Background()
	db, err := handlers.InitDB(ctx, isLocalDB)

	if err != nil {
		log.Fatal(err)
	}

	distanceCalculator := handlers.NewDistanceCalculator(db)

	r := mux.NewRouter()

	r.HandleFunc("/airports", handlers.GetAirports(db)).Methods("GET")
	r.HandleFunc("/distances", distanceCalculator.CalculateDistance).Methods("POST")
	r.HandleFunc("/aircrafts", handlers.GetAircraft(db)).Methods("GET")
	//r.HandleFunc("/flights", handlers.GetFlights(db)).Methods("GET")

	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
