// models/distance_request.go
package models

// both IATA codes
type DistanceRequest struct {
	Origin      string   `json:"origin"`
	Destination string   `json:"destination"`
	Borders     []string `json:"borders"`
}
