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
	var path []models.Coordinate

	if len(req.Borders) != 0 {
		var path []models.Coordinate
		distances, path = dc.calculateAdjustedDistance(models.Coordinate{Lat: departureAirport.Latitude, Lng: departureAirport.Longitude}, models.Coordinate{Lat: destinationAirport.Latitude, Lng: destinationAirport.Longitude}, req.Borders)
		return models.DistanceData{
			Route:     req,
			Distances: distances,
			Path:      path,
			Midpoint:  dc.calculateMidPoint(path),
		}, nil
	}

	var coords = []models.Coordinate{{Lat: departureAirport.Latitude, Lng: departureAirport.Longitude}, {Lat: destinationAirport.Latitude, Lng: destinationAirport.Longitude}}
	// var coords = []models.Coordinate{{Lat: departureAirport.Latitude, Lng: departureAirport.Longitude}, {Lat: 33.9416, Lng: -118.4085}, {Lat: destinationAirport.Latitude, Lng: destinationAirport.Longitude}}
	distances = dc.calculateDirectDistance(coords)
	path = coords
	return models.DistanceData{
		Route:     req,
		Distances: distances,
		Path:      path,
		Midpoint:  dc.calculateMidPoint(path),
	}, nil
}

func (dc *DistanceCalculator) calculateAdjustedDistance(departure, destination models.Coordinate, borders []string) (map[string]float64, []models.Coordinate) {
	// This should take into account the borders and calculate the optimal distance avoiding those borders

	// For now, let's return the direct distance as a placeholder
	_, _, _ = departure, destination, borders

	// Placeholder for the adjusted distance calculation logic
	// c := 1000.0

	var distanceKm float64
	var distanceMiles float64
	var distanceNauticalMiles float64
	var path []models.Coordinate

	return map[string]float64{
		"km":    distanceKm,
		"miles": distanceMiles,
		"nm":    distanceNauticalMiles,
	}, path
}

func (dc *DistanceCalculator) calculateDirectDistance(points []models.Coordinate) map[string]float64 {

	var c float64 = 0

	for i := range len(points) - 1 {

		// Convert latitude and longitude to radians
		lat1Rad := points[i].Lat * math.Pi / 180
		lon1Rad := points[i].Lng * math.Pi / 180
		lat2Rad := points[i+1].Lat * math.Pi / 180
		lon2Rad := points[i+1].Lng * math.Pi / 180

		// Haversine formula
		dlat := lat2Rad - lat1Rad
		dlon := lon2Rad - lon1Rad
		a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
		c += 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	}

	return dc.calculateDistanceValues(c)
}

func (dc *DistanceCalculator) calculateMidPoint(coords []models.Coordinate) models.Coordinate {
	if len(coords) == 0 {
		// Return a zero-value Coordinate if the slice is empty
		return models.Coordinate{}
	}

	var x, y, z float64

	// Convert geographical coordinates to 3D Cartesian coordinates (assuming a unit sphere)
	for _, p := range coords {
		// Convert degrees to radians
		latRad := p.Lat * math.Pi / 180
		lngRad := p.Lng * math.Pi / 180

		// Calculate Cartesian coordinates
		x += math.Cos(latRad) * math.Cos(lngRad)
		y += math.Cos(latRad) * math.Sin(lngRad)
		z += math.Sin(latRad)
	}

	// Calculate the average Cartesian coordinates
	numCoords := float64(len(coords))
	avgX := x / numCoords
	avgY := y / numCoords
	avgZ := z / numCoords

	// Convert average Cartesian coordinates back to geographical coordinates
	// atan2(y, x) gives the angle in radians
	lngRad := math.Atan2(avgY, avgX)
	// atan2(z, sqrt(x^2 + y^2)) gives the elevation angle (latitude)
	latRad := math.Atan2(avgZ, math.Sqrt(avgX*avgX+avgY*avgY))

	// Convert radians back to degrees
	avgLat := latRad * 180 / math.Pi
	avgLng := lngRad * 180 / math.Pi

	// Return the calculated midpoint (spherical centroid)
	return models.Coordinate{Lat: avgLat, Lng: avgLng}
}

func (dc *DistanceCalculator) calculateDistanceValues(c float64) map[string]float64 {
	distanceKm := earthRadiusKm * c

	return map[string]float64{
		"km":    distanceKm,
		"miles": distanceKm * milesPerKm,
		"nm":    distanceKm * nauticalMilesPerKm,
	}
}
