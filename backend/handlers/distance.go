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

func (dc *DistanceCalculator) calculateMidPoint(coords []models.PointCoords) models.PointCoords {
	if len(coords) == 0 {
		// Return a zero-value PointCoords if the slice is empty
		return models.PointCoords{}
	}

	var totalLat float64
	var totalLng float64

	// Sum up all latitudes and longitudes
	for _, p := range coords {
		totalLat += p.Lat
		totalLng += p.Lng
	}

	// Calculate the average latitude and longitude
	avgLat := totalLat / float64(len(coords))
	avgLng := totalLng / float64(len(coords))

	// Return the calculated midpoint
	return models.PointCoords{Lat: avgLat, Lng: avgLng}
}

func (dc *DistanceCalculator) calculateDirectDistance(departure, destination models.PointCoords) map[string]float64 {
	// Convert latitude and longitude to radians
	lat1Rad := departure.Lat * math.Pi / 180
	lon1Rad := departure.Lng * math.Pi / 180
	lat2Rad := destination.Lat * math.Pi / 180
	lon2Rad := destination.Lng * math.Pi / 180

	// Haversine formula
	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad
	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return calculateDistanceValues(c)
}

func (dc *DistanceCalculator) calculateAdjustedDistance(departure, destination models.PointCoords, borders []string) (map[string]float64, []models.PointCoords) {
	// This should take into account the borders and calculate the optimal distance avoiding those borders

	// For now, let's return the direct distance as a placeholder

	// Placeholder for the adjusted distance calculation logic
	// c := 1000.0

	var distanceKm float64
	var distanceMiles float64
	var distanceNauticalMiles float64
	var path []models.PointCoords

	return map[string]float64{
		"km":    distanceKm,
		"miles": distanceMiles,
		"nm":    distanceNauticalMiles,
	}, path
}

func (dc *DistanceCalculator) CalculateDistance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Departure   string   `json:"departure"`
		Destination string   `json:"destination"`
		Borders     []string `json:"borders"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Departure == "" || req.Destination == "" {
		http.Error(w, "Departure and destination IATA codes are required", http.StatusBadRequest)
		return
	}

	var departureAirport, destinationAirport models.Airport
	if err := dc.DB.NewSelect().Model(&departureAirport).Where("iata = ?", req.Departure).Scan(r.Context()); err != nil {
		http.Error(w, "Departure airport not found", http.StatusNotFound)
		return
	}
	if err := dc.DB.NewSelect().Model(&destinationAirport).Where("iata = ?", req.Destination).Scan(r.Context()); err != nil {
		http.Error(w, "Destination airport not found", http.StatusNotFound)
		return
	}

	var distances map[string]float64
	var res models.DistanceData
	var path []models.PointCoords

	if len(req.Borders) != 0 {
		var path []models.PointCoords
		distances, path = dc.calculateAdjustedDistance(models.PointCoords{Lat: departureAirport.Latitude, Lng: departureAirport.Longitude}, models.PointCoords{Lat: destinationAirport.Latitude, Lng: destinationAirport.Longitude}, req.Borders)
		res = models.DistanceData{Route: req, Distances: distances, Path: path, Midpoint: dc.calculateMidPoint(path)}

	} else {
		distances = dc.calculateDirectDistance(models.PointCoords{Lat: departureAirport.Latitude, Lng: departureAirport.Longitude}, models.PointCoords{Lat: destinationAirport.Latitude, Lng: destinationAirport.Longitude})
		path = []models.PointCoords{{Lat: departureAirport.Latitude, Lng: departureAirport.Longitude}, {Lat: destinationAirport.Latitude, Lng: destinationAirport.Longitude}}
		res = models.DistanceData{Route: req, Distances: distances, Path: path, Midpoint: dc.calculateMidPoint(path)}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)

}
