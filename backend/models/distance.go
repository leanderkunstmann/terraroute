// models/distance_request.go
package models

// both IATA codes
type DistanceRequest struct {
	Departure   string   `json:"departure"`
	Destination string   `json:"destination"`
	Borders     []string `json:"borders"`
}

type PointCoords struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type DistanceData struct {
	Route     DistanceRequest    `json:"route"`
	Distances map[string]float64 `json:"distances"`
	Path      []PointCoords      `json:"path"`
	Midpoint  PointCoords        `json:"midpoint"`
}
