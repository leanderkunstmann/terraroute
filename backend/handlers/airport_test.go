package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leanderkunstmann/terraroute/backend/database"
	"github.com/leanderkunstmann/terraroute/backend/models"
	"github.com/leanderkunstmann/terraroute/backend/services"
	"github.com/stretchr/testify/assert"
)

func TestGetAirports(t *testing.T) {
	ctx := t.Context()
	db, err := database.New(ctx, &database.Config{LocalDB: true})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
	}()

	service := services.NewAirportService(db)
	handler := NewAirportHandler(service)

	const path = "/airports"

	// Test case 1: No filters
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, path, http.NoBody)
	rr := httptest.NewRecorder()
	handler.getAirports(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var airports []models.Airport
	json.Unmarshal(rr.Body.Bytes(), &airports)
	assert.Equal(t, 4, len(airports)) // Check if all airports are returned

	// Test case 2: Filter by country
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s?country=USA", path), http.NoBody)
	rr = httptest.NewRecorder()
	handler.getAirports(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	json.Unmarshal(rr.Body.Bytes(), &airports)
	assert.Equal(t, 2, len(airports)) // Check if correct number of airports are returned

	// Test case 3: Filter by iata
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s?iata=JFK", path), http.NoBody)
	rr = httptest.NewRecorder()
	handler.getAirports(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	json.Unmarshal(rr.Body.Bytes(), &airports)
	assert.Equal(t, 1, len(airports))
	assert.Equal(t, "JFK", airports[0].IATA)

	// Test case 4: Filter by continent
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s?iata=JFK", path), http.NoBody)
	rr = httptest.NewRecorder()
	handler.getAirports(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	json.Unmarshal(rr.Body.Bytes(), &airports)
	assert.Equal(t, 1, len(airports))
	assert.Equal(t, "JFK", airports[0].IATA)

}
