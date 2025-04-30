package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/leanderkunstmann/terraroute/backend/models"
	"github.com/uptrace/bun"
)

func GetAirports(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var airports []models.Airport

		// Get query parameters
		country := r.URL.Query().Get("country")
		iata := r.URL.Query().Get("iata")
		continent := r.URL.Query().Get("continent")

		// Build the query with optional filters
		query := db.NewSelect().Model(&airports)
		if country != "" {
			query.Where("country = ?", country)
		}
		if iata != "" {
			query.Where("iata = ?", iata)
		}
		if continent != "" {
			query.Where("continent = ?", continent)
		}

		// Execute the query
		if err := query.Scan(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(airports)
	}
}
