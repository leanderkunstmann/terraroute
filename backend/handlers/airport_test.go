package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leanderkunstmann/terraroute/backend/models"
	"github.com/stretchr/testify/assert"
)

func TestGetAirports(t *testing.T) {
	ctx := context.Background()
	db, err := InitDB(ctx, "true")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
	}()

	handler := GetAirports(db)

	// TODO: use url package to build the path
	path := "/airports"

	// Test case 1: No filters
	req, _ := http.NewRequest("GET", path, nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var airports []models.Airport
	json.Unmarshal(rr.Body.Bytes(), &airports)
	assert.Equal(t, 3, len(airports)) // Check if all airports are returned

	// Test case 2: Filter by country
	req, _ = http.NewRequest("GET", fmt.Sprintf("%s?country=USA", path), nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	json.Unmarshal(rr.Body.Bytes(), &airports)
	assert.Equal(t, 2, len(airports)) // Check if correct number of airports are returned

	// Test case 3: Filter by iata
	req, _ = http.NewRequest("GET", fmt.Sprintf("%s?iata=JFK", path), nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	json.Unmarshal(rr.Body.Bytes(), &airports)
	assert.Equal(t, 1, len(airports))
	assert.Equal(t, "JFK", airports[0].IATA)

}
