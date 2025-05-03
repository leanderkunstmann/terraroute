package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/leanderkunstmann/terraroute/backend/models"
	"github.com/uptrace/bun"
)

func GetAircraft(db *bun.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var aircrafts []models.Aircraft

		// Get query parameters

		manufacturer := r.URL.Query().Get("manufacturer")
		aircraftType := r.URL.Query().Get("aircraftType")

		// Build the query with optional filters

		query := db.NewSelect().Model(&aircrafts)
		if manufacturer != "" {
			query.Where("manufacturer = ?", manufacturer)
		}
		if aircraftType != "" {
			query.Where("type = ?", aircraftType)
		}

		// Execute the query
		if err := query.Scan(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(aircrafts)
	}
}
