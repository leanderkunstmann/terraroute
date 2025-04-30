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
		expectedBody   map[string]float64
	}{
		{
			name:           "Valid request",
			requestBody:    models.DistanceRequest{Origin: "JFK", Destination: "LAX"},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]float64{
				"km":    3974,
				"miles": 2470,
				"nm":    2145,
			},
		},
		{
			name:           "Invalid request body",
			requestBody:    models.DistanceRequest{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Origin airport not found",
			requestBody:    models.DistanceRequest{Origin: "XXX", Destination: "LAX"},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Destination airport not found",
			requestBody:    models.DistanceRequest{Origin: "JFK", Destination: "XXX"},
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
				var response map[string]float64
				json.NewDecoder(rr.Body).Decode(&response)
				assert.InDeltaMapValues(t, tt.expectedBody, response, 1)
			}
		})
	}
}
