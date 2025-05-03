package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leanderkunstmann/terraroute/backend/models"
	"github.com/stretchr/testify/assert"
)


func TestCalculateDistance(t *testing.T) {
	ctx := context.Background()
	db, err := InitDB(ctx, "true")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
	}()

	handler := NewDistanceCalculator(db)

	tests := []struct {
		name           string
		requestBody    models.DistanceRequest
		expectedStatus int
		expectedBody   models.DistanceData
	}{
		{
			name:           "Valid request",
			requestBody:    models.DistanceRequest{Departure: "JFK", Destination: "LAX"},
			expectedStatus: http.StatusOK,
			expectedBody: models.DistanceData{
				Route: models.DistanceRequest{
					Departure: "JFK",
					Destination: "LAX",
					Borders: []string{},
				},
				Distances: map[string]float64{
					"km":    3974,
					"miles": 2470,
					"nm":    2145,
				},
				Path: []models.PointCoords{
					{Lat: 40.6413, Lng: -73.7781},
					{Lat: 33.9416, Lng: -118.4085},
				},
				Midpoint: models.PointCoords{
					Lat: 37.29145,
					Lng: -96.0933,
				},
			},
		},
		{
			name:           "Invalid request body",
			requestBody:    models.DistanceRequest{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Departure airport not found",
			requestBody:    models.DistanceRequest{Departure: "XXX", Destination: "LAX"},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Destination airport not found",
			requestBody:    models.DistanceRequest{Departure: "JFK", Destination: "XXX"},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/distances", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()
			handler.CalculateDistance(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				var response models.DistanceData
				json.NewDecoder(rr.Body).Decode(&response)
				assert.InDeltaMapValues(t, tt.expectedBody.Distances, response.Distances, 1)
				assert.InDelta(t, tt.expectedBody.Midpoint.Lat, response.Midpoint.Lat, 1)
				assert.InDelta(t, tt.expectedBody.Midpoint.Lng, response.Midpoint.Lng, 1)
			}
		})
	}
}
