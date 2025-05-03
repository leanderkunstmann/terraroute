import type { Route } from './+types/home'
import Globe from 'react-globe.gl'
import type { GlobeMethods } from 'react-globe.gl'
import { useEffect, useState, useMemo, useRef } from 'react'
import * as React from 'react'
import { Button, Grid, Skeleton } from '@mui/material'
import TextField from '@mui/material/TextField'
import Autocomplete from '@mui/material/Autocomplete'
import axios from 'axios'

import type { Airport, DistanceData, GeoLabel, GeoPath } from '~/types'

export function meta({}: Route.MetaArgs) {
  return [
    { title: 'Globe View' },
    { name: 'description', content: 'Get your optimal flight routes' },
  ]
}

export default function Routes() {
  const [width, setWidth] = useState(window.innerWidth)
  const [height, setHeight] = useState(window.innerHeight)
  const globeContainerRef = React.useRef<HTMLDivElement | null>(null)
  const globeEl = useRef<GlobeMethods | undefined>(undefined)

  const [airports, setAirports] = useState<Airport[]>([])
  const [route, setRoute] = useState<DistanceData | null>(null)

  const [departureAirport, setDepartureAirport] = useState<Airport | null>(null)
  const [destinationAirport, setDestinationAirport] = useState<Airport | null>(
    null
  )

  // State to hold data for the paths (lines) drawn on the globe
  const [pathsData, setPathsData] = useState<GeoPath[]>([])
  // State to hold data for the labels displayed on the globe
  const [labelsData, setLabelsData] = useState<GeoLabel[]>([])

  useEffect(() => {
    const handleResize = () => {
      if (globeContainerRef.current) {
        const { clientWidth, clientHeight } = globeContainerRef.current

        setWidth(clientWidth)
        setHeight(clientHeight)
      }
    }

    handleResize() // Set initial dimensions

    const getAirports = async () => {
      const res = await axios.get('http://192.168.0.178:8080/api/v1/airports')
      setAirports(res.data)
    }
    getAirports()

    window.addEventListener('resize', handleResize)
    return () => {
      window.removeEventListener('resize', handleResize)
    }
  }, [globeContainerRef])

  useMemo(() => {
    if (labelsData.length > 0) {
      globeEl.current?.pointOfView({
        lat: route?.midpoint.lat,
        lng: route?.midpoint.lng,
        altitude: width > 1024 ? 2 : 3,
      })
    } else {
      globeEl.current?.pointOfView({
        lat: 40,
        lng: 0,
        altitude: width > 1024 ? 2 : 3,
      })
    }
  }, [width])

  // Create the base options list once, including the original airport data
  // Use useMemo if airports list is large or changes frequently, otherwise a simple map outside component is fine
  const baseAirportOptions = useMemo(() => {
    return airports.map((airport) => ({
      // Keep the original airport data in the option object
      ...airport,
      // Create the label format you want to display
      label: airport.iata + ' | ' + airport.name,
    }))
  }, [airports]) // Re-run if the original airports list changes

  // Filter options for Departure: exclude the selected Destination airport
  const departureOptions = useMemo(() => {
    if (!destinationAirport) {
      return baseAirportOptions // If no destination is selected, all airports are available for departure
    }
    return baseAirportOptions.filter(
      (option) => option.iata !== destinationAirport.iata // Compare using a unique identifier like IATA
    )
  }, [baseAirportOptions, destinationAirport]) // Re-run if base options or destination changes

  // Filter options for Destination: exclude the selected Departure airport
  const destinationOptions = useMemo(() => {
    if (!departureAirport) {
      return baseAirportOptions // If no departure is selected, all airports are available for destination
    }
    return baseAirportOptions.filter(
      (option) => option.iata !== departureAirport.iata // Compare using a unique identifier like IATA
    )
  }, [baseAirportOptions, departureAirport]) // Re-run if base options or departure changes

  useMemo(async () => {
    if (departureAirport && destinationAirport) {
      // const getRoute = async () => {
      //   const res = await axios.post(
      //     'http://192.168.0.178:8080/api/v1/distances',
      //     {
      //       departure: departureAirport.iata,
      //       destination: destinationAirport.iata,
      //       borders: [],
      //     }
      //   )
      //   setRoute(await res.data)
      // }
      // await getRoute()
    } else {
      setPathsData([])
      setLabelsData([])
    }
  }, [departureAirport, destinationAirport])

  const handleRequest = async () => {
    if (departureAirport && destinationAirport) {
      const getRoute = async () => {
        const res = await axios.post(
          'http://192.168.0.178:8080/api/v1/distances',
          {
            departure: departureAirport.iata,
            destination: destinationAirport.iata,
            borders: [],
          }
        )
        setRoute(await res.data)
      }

      await getRoute()
    }
  }

  useMemo(() => {
    if (
      route &&
      route.path &&
      route.path.length >= 2 &&
      departureAirport &&
      destinationAirport
    ) {
      const paths = []

      const pathPoints = []
      for (let r of route.path) {
        pathPoints.push([r.lat, r.lng])
      }

      const colors: string[][] = [
        ['rgba(0,0,255,1)', 'rgba(255,0,0,1)'],
        ['rgb(113, 126, 0)', 'rgba(255,0,0,1)'],
        ['rgba(255,0,255,1)', 'rgba(0,255,0,1)'],
      ]
      paths.push({
        coords: pathPoints,
        properties: {
          label: route.route.departure + '-->' + route.route.destination,
          color: colors[Math.floor(Math.random() * colors.length)],
        },
      })

      setPathsData(paths)

      const labels: GeoLabel[] = []

      // push departure
      labels.push(
        {
          lat: route.path[0].lat,
          lng: route.path[0].lng,
          text: `${route.route.departure}`,
          size: route.distances.km > 3000 ? 3 : 1,
          radius: 0.5,
          dot: true,
          color: 'white',
          alt: 0,
        },
        {
          lat: route.path[route.path.length - 1].lat,
          lng: route.path[route.path.length - 1].lng,
          text: `${route.route.destination}`,
          size: route.distances.km > 3000 ? 3 : 1,
          radius: 0.5,
          dot: true,
          color: 'white',
          alt: 0,
        },
        {
          lat:
            route.midpoint.lat < 85
              ? route.midpoint.lat + 5
              : route.midpoint.lat - 5,
          lng: route.midpoint.lng,
          text:
            route.distances.km > 3000
              ? `${route.distances.km.toFixed(
                  0
                )} km / ${route.distances.nm.toFixed(0)} nm \n  from ${
                  route.route.departure
                } to ${route.route.destination}`
              : `${route.distances.km.toFixed(
                  0
                )} km / ${route.distances.nm.toFixed(0)} NM`,
          size: 1.8,
          radius: 0.5,
          dot: false,
          color: '#F0F8FF',
          alt: route.distances.km > 3000 ? 0 : 0,
        }
      )

      setLabelsData(labels)
      globeEl.current?.pointOfView({
        lat: route.midpoint.lat,
        lng: route.midpoint.lng,
        altitude:
          width > 1024
            ? route.distances.km > 3000
              ? 2
              : route.distances.km > 1024
              ? 1.7
              : 1
            : route.distances.km > 3000
            ? 3
            : route.distances.km > 1024
            ? 1.5
            : 1,
      })
    } else {
      // If distanceData is not valid, clear the paths and labels from the globe
      setPathsData([])
      setLabelsData([])
    }
  }, [route])

  return (
    <div>
      <div
        ref={globeContainerRef}
        style={{
          width: '100% !important',
          height: width > 1024 ? '80dvh' : '66dvh',
          alignContent: 'center !important',
        }}
      >
        <Globe
          ref={globeEl}
          width={width}
          height={height}
          globeImageUrl="earth-blue-marble.jpg"
          bumpImageUrl="earth-topology.png"
          backgroundColor="rgba(0, 0, 0, 0)"
          showAtmosphere={false}
          pathsData={pathsData.length > 0 ? pathsData : undefined}
          pathPoints={'coords'}
          pathColor={(path: any) => path.properties.color}
          pathDashLength={0.1}
          pathStroke={1}
          pathDashGap={0.02}
          pathDashAnimateTime={10000}
          pathLabel={(path: any) => path.properties.label}
          labelsData={labelsData.length > 0 ? labelsData : undefined}
          labelResolution={10}
          labelColor={'color'}
          labelSize={'size'}
          labelAltitude={'alt'}
          labelDotRadius={'radius'}
          labelIncludeDot={'dot'}
        />
      </div>
      <div
        style={{
          width: '100%',
          height: width > 1024 ? '20dvh' : '34dvh',
          alignContent: 'center',
        }}
      >
        <Grid
          container
          spacing={1}
          sx={{
            justifyContent: 'top',
            alignItems: 'top',
            mt: '8px',
          }}
          direction="row"
        >
          <Grid size={{ xs: 12, sm: 4 }}>
            <div style={{ display: 'flex', justifyContent: 'center' }}>
              {airports.length !== 0 ? (
                <Autocomplete
                  disablePortal
                  blurOnSelect
                  id="departure-airport-autocomplete"
                  size={width > 1024 ? 'medium' : 'small'}
                  options={departureOptions} // Use the FILTERED options for departure
                  value={departureAirport} // Control the value
                  onChange={(event, newValue) => {
                    setDepartureAirport(newValue) // Update the state on selection
                  }}
                  isOptionEqualToValue={(option, value) =>
                    option.iata === value.iata
                  } // Helps determine if an option matches the current value
                  sx={{ width: 300 }}
                  renderInput={(params) => (
                    <TextField {...params} label="Departure" />
                  )}
                />
              ) : (
                <Skeleton variant="rectangular" width={280} height={50} />
              )}
            </div>
          </Grid>
          <Grid size={{ xs: 12, sm: 4 }}>
            <div style={{ display: 'flex', justifyContent: 'center' }}>
              {airports.length !== 0 ? (
                <Autocomplete
                  disablePortal
                  blurOnSelect
                  id="destination-airport-autocomplete"
                  size={width > 1024 ? 'medium' : 'small'}
                  options={destinationOptions} // Use the FILTERED options for destination
                  value={destinationAirport} // Control the value
                  onChange={(event, newValue) => {
                    setDestinationAirport(newValue) // Update the state on selection
                  }}
                  isOptionEqualToValue={(option, value) =>
                    option.iata === value.iata
                  } // Helps determine if an option matches the current value
                  sx={{ width: 300 }}
                  renderInput={(params) => (
                    <TextField {...params} label="Destination" />
                  )}
                />
              ) : (
                <Skeleton variant="rectangular" width={280} height={50} />
              )}
            </div>
          </Grid>
          <Grid size={{ xs: 12, sm: 4 }}>
            <div style={{ display: 'flex', justifyContent: 'center' }}>
              {airports.length !== 0 ? (
                <Autocomplete
                  disablePortal
                  blurOnSelect
                  id="excluded-countries-autocomplete"
                  size={width > 1024 ? 'medium' : 'small'}
                  options={destinationOptions} // Use the FILTERED options for destination
                  isOptionEqualToValue={(option, value) =>
                    option.iata === value.iata
                  } // Helps determine if an option matches the current value
                  sx={{ width: 300 }}
                  renderInput={(params) => (
                    <TextField {...params} label="Exclude Countries" />
                  )}
                />
              ) : (
                <Skeleton variant="rectangular" width={280} height={50} />
              )}
            </div>
          </Grid>
          <Grid size={{ xs: 12, sm: 12 }}>
            <div style={{ display: 'flex', justifyContent: 'center' }}>
              <Button
                variant="outlined"
                disabled={!destinationAirport || !departureAirport}
                onClick={handleRequest}
              >
                Generate Route
              </Button>
            </div>
          </Grid>
        </Grid>
      </div>
    </div>
  )
}
