package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leanderkunstmann/terraroute/backend/services"
)

var _ Handler = (*CountryHandler)(nil)

type CountryHandler struct {
	service *services.CountryService
}

func NewCountryHandler(svc *services.CountryService) *CountryHandler {
	return &CountryHandler{service: svc}
}

func (ah *CountryHandler) Register(r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("%s/countries", basePathV1), ah.getCountries).
		Methods(http.MethodGet, http.MethodOptions).
		Name("GetCountries")
	r.HandleFunc(fmt.Sprintf("%s/countries/{code}", basePathV1), ah.getCountry).
		Methods(http.MethodGet, http.MethodOptions).
		Name("GetCountry")
}

func (ah *CountryHandler) getCountries(w http.ResponseWriter, r *http.Request) {
	continent := r.URL.Query().Get("continent")

	res, err := ah.service.ListCountries(r.Context(), continent)

	if err != nil {
		if errors.Is(err, services.ErrCountriesNotFound) {
			newErrorResponse(w, err, http.StatusNotFound)
		} else {
			newErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(res); err != nil {
		newErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (ah *CountryHandler) getCountry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code, ok := vars["code"]

	if !ok || code == "" {
		newErrorResponse(w, fmt.Errorf("country code not found in path"), http.StatusBadRequest)
		return
	}

	res, err := ah.service.ListCountry(r.Context(), code)
	if err != nil {
		if errors.Is(err, services.ErrCountryNotFound) {
			newErrorResponse(w, err, http.StatusNotFound)
		} else {
			newErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(res); err != nil {
		newErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
