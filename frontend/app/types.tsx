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
export interface Coordinate {
  lat: number // Latitude coordinate.
  lng: number // Longitude coordinate.
}

export interface Country {
  code: string
  name: string
  continent: string
  borders?: Coordinate[][]
}

// Represents the data returned after calculating distances and routes.
export interface DistanceData {
  route: DistanceRequest // The original route request details.
  distances: Record<string, number> // A map of distances, keyed by some identifier (e.g., border name).
  path: Coordinate[] // An array of coordinates representing the calculated path.
  midpoint: Coordinate // The calculated midpoint of the route.
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

export interface GeoArc {
  startLat: number
  startLng: number
  endLat: number
  endLng: number
  label: string
  color: string | string[]
  alt: number
}

export interface GeoPathProperties {
  label: string
  color: string | string[]
}
