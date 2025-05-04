package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/leanderkunstmann/terraroute/backend/models"
	"github.com/uptrace/bun"
)

type DistanceCalculator struct {
	db *bun.DB
}

func NewDistanceCalculator(db *bun.DB) *DistanceCalculator {
	return &DistanceCalculator{db: db}
}

const (
	// earthRadiusKm is the radius of the Earth in kilometers.
	earthRadiusKm = 6371
	// milesPerKm is the conversion factor from kilometers to miles.
	milesPerKm = 0.621371
	// nauticalMilesPerKm is the conversion factor from kilometers to nautical miles.
	nauticalMilesPerKm = 0.539957
)

type NotFoundError string

func (e NotFoundError) Error() string {
	return string(e)
}

func (e NotFoundError) Code() int {
	return http.StatusNotFound
}

func (dc *DistanceCalculator) CalculateDistance(ctx context.Context, req *models.DistanceRequest) (models.DistanceData, error) {
	var departureAirport, destinationAirport models.Airport
	if err := dc.db.NewSelect().Model(&departureAirport).Where("iata = ?", req.Departure).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.DistanceData{}, NotFoundError(fmt.Sprintf("departure airport not found: %s", req.Departure))
		}
		return models.DistanceData{}, fmt.Errorf("failed to find departure airport: %w", err)
	}
	if err := dc.db.NewSelect().Model(&destinationAirport).Where("iata = ?", req.Destination).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.DistanceData{}, NotFoundError(fmt.Sprintf("destination airport not found: %s", req.Destination))
		}
		return models.DistanceData{}, fmt.Errorf("failed to find destination airport: %w", err)
	}

	var distances map[string]float64
	var path []models.PointCoords

	if len(req.Borders) != 0 {
		var path []models.PointCoords
		distances, path = dc.calculateAdjustedDistance(models.PointCoords{Lat: departureAirport.Latitude, Lng: departureAirport.Longitude}, models.PointCoords{Lat: destinationAirport.Latitude, Lng: destinationAirport.Longitude}, req.Borders)
		return models.DistanceData{
			Route:     req,
			Distances: distances,
			Path:      path,
			Midpoint:  dc.calculateMidPoint(path),
		}, nil
	}

	distances = dc.calculateDirectDistance(models.PointCoords{Lat: departureAirport.Latitude, Lng: departureAirport.Longitude}, models.PointCoords{Lat: destinationAirport.Latitude, Lng: destinationAirport.Longitude})
	path = []models.PointCoords{{Lat: departureAirport.Latitude, Lng: departureAirport.Longitude}, {Lat: destinationAirport.Latitude, Lng: destinationAirport.Longitude}}
	return models.DistanceData{
		Route:     req,
		Distances: distances,
		Path:      path,
		Midpoint:  dc.calculateMidPoint(path),
	}, nil
}

func (dc *DistanceCalculator) calculateAdjustedDistance(departure, destination models.PointCoords, borders []string) (map[string]float64, []models.PointCoords) {
	// This should take into account the borders and calculate the optimal distance avoiding those borders

	// For now, let's return the direct distance as a placeholder
	_, _, _ = departure, destination, borders

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

	return dc.calculateDistanceValues(c)
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

func (dc *DistanceCalculator) calculateDistanceValues(c float64) map[string]float64 {
	distanceKm := earthRadiusKm * c

	return map[string]float64{
		"km":    distanceKm,
		"miles": distanceKm * milesPerKm,
		"nm":    distanceKm * nauticalMilesPerKm,
	}
}
