package models

type Country struct {
	Code      string `json:"code"` // unique identifier
	Name      string `json:"name"`
	Continent string `json:"continent"`
}

type CountryBorders struct {
	Code    string  `json:"code"` // foreign key
	Borders GeoJson `json:"borders"`
}

type CountryBordersLocal struct {
	Code    string `json:"code"` // foreign key
	Borders string `json:"borders"`
}

type GeoJsonGeometry struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

type GeoJsonFeature struct {
	Type       string          `json:"type"`
	Geometry   GeoJsonGeometry `json:"geometry"`
	Properties struct{}        `json:"properties"`
}

type GeoJson struct {
	Type     string           `json:"type"`
	Features []GeoJsonFeature `json:"features"`
}
