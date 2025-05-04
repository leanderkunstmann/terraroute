package models

type Country struct {
	Code      string `json:"code"` // unique identifier
	Name      string `json:"name"`
	Continent string `json:"continent"`
}

type FullCountry struct {
	Country
	Borders [][]Coordinate `json:"borders"`
}
