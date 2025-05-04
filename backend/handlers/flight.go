package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var _ Handler = (*FlightHandler)(nil)

type FlightHandler struct{}

func NewFlightHandler() *FlightHandler {
	return &FlightHandler{}
}

func (fh *FlightHandler) Register(r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("%s/flights", basePathV1), fh.GetFlights).
		Methods(http.MethodGet, http.MethodOptions).
		Name("GetFlights")
}

func (fh *FlightHandler) GetFlights(w http.ResponseWriter, r *http.Request) {
	// Placeholder for handling flight requests
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Flight handler placeholder"))
}
