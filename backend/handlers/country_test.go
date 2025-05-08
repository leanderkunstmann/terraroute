package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/leanderkunstmann/terraroute/backend/database"
	"github.com/leanderkunstmann/terraroute/backend/models"
	"github.com/leanderkunstmann/terraroute/backend/services"
	"github.com/stretchr/testify/assert"
)

func TestGetCountries(t *testing.T) {
	ctx := t.Context()
	db, err := database.New(ctx, &database.Config{LocalDB: true})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
	}()

	service := services.NewCountryService(db)
	handler := NewCountryHandler(service)

	const path = "/countries"

	// Test /countries

	// Test case 1: No filters

	var countries []models.Country

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, path, http.NoBody)
	rr := httptest.NewRecorder()
	handler.getCountries(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	json.Unmarshal(rr.Body.Bytes(), &countries)
	assert.Equal(t, 4, len(countries)) // Check if all countries are returned

	// Test case 2: Filter by continent
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s?continent=Europe", path), http.NoBody)
	rr = httptest.NewRecorder()
	handler.getCountries(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	json.Unmarshal(rr.Body.Bytes(), &countries)
	assert.Equal(t, 2, len(countries)) // Check if correct number of countries are returned

	// Test case 3: Filter by continent wrong
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s?continent=Europee", path), http.NoBody)
	rr = httptest.NewRecorder()
	handler.getCountries(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Test /countries/code

	// Test case 1: No filters

	var country models.CountryBorders

	// needed for handling of path variables, TODO: implement like that in all test files
	router := mux.NewRouter()
	router.HandleFunc("/countries/{code}", handler.getCountry).Methods("GET")

	newPath := "/countries/de"
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, newPath, http.NoBody)
	rr = httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	json.Unmarshal(rr.Body.Bytes(), &country)
	assert.Equal(t, "DE", country.Code) // Check if all countries are returned

	newPath = "/countries/test"
	// Test case 2: Filter by continent
	req, _ = http.NewRequestWithContext(ctx, http.MethodGet, newPath, http.NoBody)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}
