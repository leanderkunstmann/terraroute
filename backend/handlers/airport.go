package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leanderkunstmann/terraroute/backend/services"
)

var _ Handler = (*AirportHandler)(nil)

type AirportHandler struct {
	service *services.AirportService
}

func NewAirportHandler(svc *services.AirportService) *AirportHandler {
	return &AirportHandler{service: svc}
}

func (ah *AirportHandler) Register(r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("%s/airports", basePathV1), ah.getAirports).
		Methods(http.MethodGet, http.MethodOptions).
		Name("GetAirports")
}

func (ah *AirportHandler) getAirports(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	iata := r.URL.Query().Get("iata")
	continent := r.URL.Query().Get("continent")

	res, err := ah.service.ListAirports(r.Context(), iata, continent, country)
	if err != nil {
		newErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(res); err != nil {
		newErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
