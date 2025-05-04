// models/distance_request.go
package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type DistanceRequest struct {
	Departure   string   `json:"departure"`
	Destination string   `json:"destination"`
	Borders     []string `json:"borders"`
}

func NewDistanceRequest(b io.Reader) (*DistanceRequest, error) {
	var req DistanceRequest
	if err := json.NewDecoder(b).Decode(&req); err != nil {
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}

	if err := req.validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	return &req, nil
}

func (req *DistanceRequest) validate() error {
	var err error
	if req.Departure == "" {
		err = errors.New("departure IATA code is required")
	}
	if req.Destination == "" {
		err = errors.Join(err, errors.New("destination IATA code is required"))
	}
	return err
}

type PointCoords struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type DistanceData struct {
	Route     *DistanceRequest   `json:"route"`
	Distances map[string]float64 `json:"distances"`
	Path      []PointCoords      `json:"path"`
	Midpoint  PointCoords        `json:"midpoint"`
}
