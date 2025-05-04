import type { Route } from './+types/home'
import Globe from 'react-globe.gl'
import type { GlobeMethods } from 'react-globe.gl'
import { useEffect, useState, useMemo, useRef } from 'react'
import * as React from 'react'
import { Button, Grid, Skeleton } from '@mui/material'
import TextField from '@mui/material/TextField'
import Autocomplete from '@mui/material/Autocomplete'
import axios from 'axios'

import type {
  Airport,
  Country,
  DistanceData,
  GeoLabel,
  GeoPath,
  GeoArc,
} from '~/types'
import { colors } from '~/constants'

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
  const [countries, setCountries] = useState<Country[]>([])
  const [route, setRoute] = useState<DistanceData | null>(null)

  const [departureAirport, setDepartureAirport] = useState<Airport | null>(null)
  const [destinationAirport, setDestinationAirport] = useState<Airport | null>(
    null
  )

  const [countriesToExclude, setCountriesToExclude] = useState<Country[]>([])

  // State to hold data for the paths (lines) drawn on the globe
  const [pathsData, setPathsData] = useState<GeoPath[]>([])
  const [arcsData, setArcsData] = useState<GeoArc[]>([])
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

    const getCountries = async () => {
      const res = await axios.get('http://192.168.0.178:8080/api/v1/countries')
      setCountries(res.data)
    }
    getCountries()

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

  const countryOptions = useMemo(() => {
    return countries.map((country) => ({
      ...country,
      label: country.name,
    }))
  }, [countries]) // Re-run if the original airports list changes

  const baseAirportOptions = useMemo(() => {
    return airports.map((airport) => ({
      ...airport,
      label:
        airport.iata +
        ' | ' +
        airport.name +
        (airport.name.toLowerCase().includes(airport.city.toLowerCase())
          ? ''
          : ' ' + airport.city),
    }))
  }, [airports]) // Re-run if the original airports list changes

  // Filter options for Departure: exclude the selected Destination airport
  const departureOptions = useMemo(() => {
    if (!destinationAirport) {
      return baseAirportOptions.sort((a, b) => (a.iata < b.iata ? -1 : 1))
    }
    return baseAirportOptions
      .filter((option) => option.iata !== destinationAirport.iata)
      .sort((a, b) => (a.iata < b.iata ? -1 : 1))
  }, [baseAirportOptions, destinationAirport])

  const destinationOptions = useMemo(() => {
    if (!departureAirport) {
      return baseAirportOptions.sort((a, b) => (a.iata < b.iata ? -1 : 1))
    }
    return baseAirportOptions
      .filter((option) => option.iata !== departureAirport.iata)
      .sort((a, b) => (a.iata < b.iata ? -1 : 1))
  }, [baseAirportOptions, departureAirport])

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
          'http://192.168.0.178:8080/api/v1/routes',
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
      const arcs: GeoArc[] = []

      const depAirport: Airport | undefined = airports.find(
        (airport) => airport.iata === route.route.departure
      )
      const desAirport: Airport | undefined = airports.find(
        (airport) => airport.iata === route.route.destination
      )
      let textLabel = ''
      if (depAirport && desAirport) {
        textLabel = `from (${depAirport.iata}) ${depAirport.city} in ${depAirport.country} to (${desAirport.iata}) ${desAirport.city} in ${desAirport.country}`
      } else {
        textLabel = route.route.departure + ' --> ' + route.route.destination
      }

      let colorNum = Math.floor(Math.random() * colors.length)

      for (let i = 0; i < route.path.length; i++) {
        pathPoints.push([route.path[i].lat, route.path[i].lng])
        if (i < route.path.length - 1) {
          const arc: GeoArc = {
            startLat: route.path[i].lat,
            startLng: route.path[i].lng,
            endLat: route.path[i + 1].lat,
            endLng: route.path[i + 1].lng,
            label: textLabel,
            color: colors[colorNum],
            alt: 0,
          }
          colorNum === colors.length - 1 ? (colorNum = 0) : colorNum++
          arcs.push(arc)
        }
      }

      setArcsData(arcs)
      console.log(arcsData)

      paths.push({
        coords: pathPoints,
        properties: {
          label: textLabel,
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
          size:
            route.distances.km > 5000
              ? 3
              : route.distances.km > 3000
              ? 2
              : route.distances.km > 3000
              ? 1.25
              : 0.75,
          radius: 0.5,
          dot: true,
          color: 'white',
          alt: 0,
        },
        {
          lat: route.path[route.path.length - 1].lat,
          lng: route.path[route.path.length - 1].lng,
          text: `${route.route.destination}`,
          size:
            route.distances.km > 5000
              ? 3
              : route.distances.km > 3000
              ? 2
              : route.distances.km > 3000
              ? 1.25
              : 0.75,
          radius: 0.5,
          dot: true,
          color: 'white',
          alt: 0,
        },
        {
          lat:
            route.distances.km < 1000
              ? route.midpoint.lat + 3
              : route.midpoint.lat > 55
              ? 55
              : route.midpoint.lat < -64
              ? -64
              : route.midpoint.lat,
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
          size:
            route.distances.km > 5000
              ? 3
              : route.distances.km > 3000
              ? 2
              : route.distances.km > 3000
              ? 1.25
              : 0.75,
          radius: 0.5,
          dot: false,
          color:
            route.midpoint.lat > 60 || route.midpoint.lat < -60
              ? 'white'
              : 'white',
          alt:
            route.distances.km > 5000
              ? 0.2
              : route.distances.km > 3000
              ? 0.03
              : 0,
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

  // pathsData={pathsData.length > 0 ? pathsData : undefined}
  // pathPoints={'coords'}
  // pathColor={(path: any) => path.properties.color}
  // pathDashLength={0.1}
  // pathStroke={1}
  // pathDashGap={0.02}
  // pathDashAnimateTime={10000}
  // pathLabel={(path: any) => path.properties.label}

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
          showAtmosphere={true}
          atmosphereColor="orange"
          atmosphereAltitude={0.05}
          arcsData={arcsData}
          arcColor={'color'}
          arcDashLength={0.1}
          arcStroke={1}
          arcDashGap={0.02}
          arcAltitude={'alt'}
          arcLabel={'label'}
          arcDashAnimateTime={10000}
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
                  blurOnSelect={'touch'}
                  autoSelect
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
                  blurOnSelect={'touch'}
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
                  blurOnSelect={'touch'}
                  id="excluded-countries-autocomplete"
                  multiple
                  value={countriesToExclude}
                  onChange={(event, newValue) => {
                    setCountriesToExclude(newValue) // Update the state on selection
                  }}
                  size={width > 1024 ? 'medium' : 'small'}
                  options={countryOptions} // Use the FILTERED options for destination
                  isOptionEqualToValue={(option, value) =>
                    option.code === value.code
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
