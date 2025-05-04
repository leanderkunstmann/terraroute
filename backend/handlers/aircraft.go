package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leanderkunstmann/terraroute/backend/services"
)

var _ Handler = (*AircraftHandler)(nil)

type AircraftHandler struct {
	service *services.AircraftService
}

func NewAircraftHandler(svc *services.AircraftService) *AircraftHandler {
	return &AircraftHandler{service: svc}
}

func (ah *AircraftHandler) Register(r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("%s/aircraft", basePathV1), ah.GetAircraft).
		Methods(http.MethodGet, http.MethodOptions).
		Name("GetAircraft")
}

func (ah *AircraftHandler) GetAircraft(w http.ResponseWriter, r *http.Request) {
	manufacturer := r.URL.Query().Get("manufacturer")
	aircraftType := r.URL.Query().Get("aircraftType")

	res, err := ah.service.ListAircraft(r.Context(), manufacturer, aircraftType)
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
