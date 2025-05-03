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

func TestAircraft(t *testing.T) {
	ctx := context.Background()
	db, err := InitDB(ctx, "true")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
	}()

	handler := GetAircraft(db)

	path := "/aircrafts"

	// Test case 1: No filters
	req, _ := http.NewRequest("GET", path, nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var aircrafts []models.Aircraft
	json.Unmarshal(rr.Body.Bytes(), &aircrafts)
	assert.Equal(t, 4, len(aircrafts)) // Check if all aircrafts are returned

	// Test case 2: Filter by manufacturer
	req, _ = http.NewRequest("GET", fmt.Sprintf("%s?manufacturer=Boeing", path), nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	json.Unmarshal(rr.Body.Bytes(), &aircrafts)
	assert.Equal(t, 1, len(aircrafts)) // Check if correct number of aircrafts are returned
	assert.Equal(t, models.Manufacturer("Boeing"), aircrafts[0].Manufacturer)

	// Test case 3: Filter by aircraftType
	req, _ = http.NewRequest("GET", fmt.Sprintf("%s?aircraftType=commercial", path), nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	json.Unmarshal(rr.Body.Bytes(), &aircrafts)
	assert.Equal(t, 2, len(aircrafts))
	assert.Equal(t, models.AircraftType("commercial"), aircrafts[0].Type)
	assert.Equal(t, models.AircraftType("commercial"), aircrafts[1].Type)

	// Test case 4: Filter by aircraftType and manufacturer
	req, _ = http.NewRequest("GET", fmt.Sprintf("%s?manufacturer=Airbus&aircraftType=commercial", path), nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	json.Unmarshal(rr.Body.Bytes(), &aircrafts)
	assert.Equal(t, 1, len(aircrafts))
	assert.Equal(t, models.Manufacturer("Airbus"), aircrafts[0].Manufacturer)



}
