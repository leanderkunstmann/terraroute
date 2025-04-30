package models

type AircraftType string

const (
	Ultralight AircraftType = "ultralight"
	Light      AircraftType = "light"
	Heavy      AircraftType = "heavy"
	Commercial AircraftType = "commercial"
	Cargo      AircraftType = "cargo"
	Military   AircraftType = "military"
)

type Manufacturer string

const (
	Airbus     Manufacturer = "Airbus"
	Antonov    Manufacturer = "Antonov"
	Boeing     Manufacturer = "Boeing"
	Bombardier Manufacturer = "Bombardier"
	Gulfstream Manufacturer = "Gulfstream"
)

type Aircraft struct {
	Id           int          `json:"id"` // unique identifier
	Type         AircraftType `json:"type"`
	Name         string       `json:"name"`
	Manufacturer Manufacturer `json:"manufacturer"`
	Range        int          `json:"range"`
}
