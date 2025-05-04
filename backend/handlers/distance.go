package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leanderkunstmann/terraroute/backend/models"
	"github.com/leanderkunstmann/terraroute/backend/services"
)

var _ Handler = (*DistanceHandler)(nil)

type DistanceHandler struct {
	service *services.DistanceCalculator
}

func NewDistanceHandler(svc *services.DistanceCalculator) *DistanceHandler {
	return &DistanceHandler{service: svc}
}

func (dc *DistanceHandler) Register(r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("%s/distance", basePathV1), dc.CalculateDistance).
		Methods(http.MethodPost, http.MethodOptions).
		Name("CalculateDistance")
}

func (dc *DistanceHandler) CalculateDistance(w http.ResponseWriter, r *http.Request) {
	req, err := models.NewDistanceRequest(r.Body)
	if err != nil {
		newErrorResponse(w, fmt.Errorf("failed to parse request: %w", err), http.StatusBadRequest)
		return
	}

	res, err := dc.service.CalculateDistance(r.Context(), req)
	if err != nil {
		newErrorResponse(w, fmt.Errorf("failed to calculate distance: %w", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(res); err != nil {
		newErrorResponse(w, fmt.Errorf("failed to encode response: %w", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
