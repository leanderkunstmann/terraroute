package handlers

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/leanderkunstmann/terraroute/backend/models"
	"github.com/uptrace/bun"
)

type DistanceCalculator struct {
	DB *bun.DB
}

func NewDistanceCalculator(db *bun.DB) *DistanceCalculator {
	return &DistanceCalculator{DB: db}
}

func calculateDistanceValues(c float64) map[string]float64 {
	// Radius of Earth in kilometers
	const R = 6371

	// Distance in kilometers
	distanceKm := R * c

	// Convert distance to miles and nautical miles
	distanceMiles := distanceKm * 0.621371
	distanceNauticalMiles := distanceKm * 0.539957

	return map[string]float64{
		"km":    distanceKm,
		"miles": distanceMiles,
		"nm":    distanceNauticalMiles,
	}
}

func (dc *DistanceCalculator) calculateDirectDistance(lat1, lon1, lat2, lon2 float64) map[string]float64 {
	// Convert latitude and longitude to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Haversine formula
	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad
	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return calculateDistanceValues(c)
}

func (dc *DistanceCalculator) calculateAdjustedDistance(lat1, lon1, lat2, lon2 float64, borders []string) map[string]float64 {
	// This should take into account the borders and calculate the optimal distance avoiding those borders

	// For now, let's return the direct distance as a placeholder
	if len(borders) == 0 {
		return dc.calculateDirectDistance(lat1, lon1, lat2, lon2)
	}

	// Placeholder for the adjusted distance calculation logic
	c := 1000.0

	return calculateDistanceValues(c)
}

func (dc *DistanceCalculator) CalculateDistance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Origin      string   `json:"origin"`
		Destination string   `json:"destination"`
		Borders     []string `json:"borders"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Origin == "" || req.Destination == "" {
		http.Error(w, "Origin and destination IATA codes are required", http.StatusBadRequest)
		return
	}

	var originAirport, destinationAirport models.Airport
	if err := dc.DB.NewSelect().Model(&originAirport).Where("iata = ?", req.Origin).Scan(r.Context()); err != nil {
		http.Error(w, "Origin airport not found", http.StatusNotFound)
		return
	}
	if err := dc.DB.NewSelect().Model(&destinationAirport).Where("iata = ?", req.Destination).Scan(r.Context()); err != nil {
		http.Error(w, "Destination airport not found", http.StatusNotFound)
		return
	}

	var distance map[string]float64

	if len(req.Borders) == 0 {
		distance = dc.calculateAdjustedDistance(originAirport.Latitude, originAirport.Longitude, destinationAirport.Latitude, destinationAirport.Longitude, req.Borders)
	} else {
		distance = dc.calculateDirectDistance(originAirport.Latitude, originAirport.Longitude, destinationAirport.Latitude, destinationAirport.Longitude)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(distance)

}
