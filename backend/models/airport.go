package models

type Airport struct {
	IATA      string  `json:"iata"` // unique identifier
	Name      string  `json:"name"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Continent string  `json:"continent"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
