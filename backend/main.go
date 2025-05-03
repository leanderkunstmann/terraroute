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

func serveRoute(router *mux.Router, api_version string, path string, method string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
	router.HandleFunc("/api" + api_version + path, f).Methods("OPTIONS")
	return router.HandleFunc("/api" + api_version + path, f).Methods(method)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust as needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // Often needed for cookies/auth headers

		// Allow specific headers, or dynamically allow requested headers for preflight
		allowedHeaders := "Content-Type, Authorization" // Start with commonly needed headers
		if r.Method == "OPTIONS" {
			// For preflight requests, echo back the requested headers if they exist
			requestedHeaders := r.Header.Get("Access-Control-Request-Headers")
			if requestedHeaders != "" {
				allowedHeaders = requestedHeaders
			}
		}
		w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)


		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

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

	r.Use(corsMiddleware)


	//v1
	serveRoute(r, "/v1", "/airports", "GET", handlers.GetAirports(db))
	serveRoute(r, "/v1", "/distances", "POST", distanceCalculator.CalculateDistance)
	serveRoute(r, "/v1", "/aircrafts", "GET", handlers.GetAircraft(db))

	//serveRoute(r, "/v1", "/flights", "GET", handlers.GetFlights(db))
	fmt.Println("Server listening on :8080")
	fmt.Println("Available routes:")
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println(pathTemplate)
		}
		return nil
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
