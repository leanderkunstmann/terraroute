// Represents an Airport with its details.
export interface Airport {
  label?: string
  iata: string // unique identifier (IATA code)
  name: string
  city: string
  country: string
  continent: string
  latitude: number
  longitude: number
}

// Represents the request structure for distance calculation.
export interface DistanceRequest {
  departure: string // The starting point for the route.
  destination: string // The ending point for the route.
  borders: string[] // A list of borders to consider or avoid.
}

// Represents geographical coordinates.
export interface PointCoords {
  lat: number // Latitude coordinate.
  lng: number // Longitude coordinate.
}

// Represents the data returned after calculating distances and routes.
export interface DistanceData {
  route: DistanceRequest // The original route request details.
  distances: Record<string, number> // A map of distances, keyed by some identifier (e.g., border name).
  path: PointCoords[] // An array of coordinates representing the calculated path.
  midpoint: PointCoords // The calculated midpoint of the route.
}

export interface GeoLabel {
  lat: number
  lng: number
  text: string
  size: number
  radius: number
  dot: boolean
  color: string
  alt: number
}

export interface GeoPath {
  coords: number[][]
  properties: GeoPathProperties
}

export interface GeoPathProperties {
  label: string
  color: string | string[]
}
