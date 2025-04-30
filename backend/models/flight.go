package models

type Flight struct {
	FlightNumber  string `json:"flight_number"` // unique identifier
	AircraftId    int    `json:"aircraft_id"`   // foreign key
	Origin        string `json:"origin"`        // foreign key
	Destination   string `json:"destination"`   // foreign key
	DepartureTime string `json:"departure_time"`
	ArrivalTime   string `json:"arrival_time"`
}
