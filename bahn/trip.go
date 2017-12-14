package bahn // import "github.com/octo/icestat/bahn"

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Legacy URL, may be used by old systems.
// const TripInfoURL = "http://ice.portal/jetty/api/v1/tripInfo"

// TripInfoURL is the URL of JSON encoded information about the train's schedule.
const TripInfoURL = "https://portal.imice.de/api1/rs/tripInfo"

// Station is a train station.
type Station struct {
	ID        string
	Name      string
	Latitude  float64
	Longitude float64
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (s *Station) UnmarshalJSON(b []byte) error {
	var parsed struct {
		EvaNr          string
		Geocoordinates struct {
			Latitude, Longitude float64
		}
		Name string
	}

	if err := json.Unmarshal(b, &parsed); err != nil {
		return err
	}

	*s = Station{
		ID:        parsed.EvaNr,
		Name:      parsed.Name,
		Latitude:  parsed.Geocoordinates.Latitude,
		Longitude: parsed.Geocoordinates.Longitude,
	}

	return nil
}

func (s Station) String() string {
	return s.Name
}

// Stop is a scheduled stop along the route.
type Stop struct {
	Station              *Station
	Platform             string
	DistanceFromStart    float64
	DistanceFromLastStop float64
	Passed               bool
	ScheduledArrival     time.Time
	ActualArrival        time.Time
	ScheduledDeparture   time.Time
	ActualDeparture      time.Time
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (s *Stop) UnmarshalJSON(b []byte) error {
	var parsed struct {
		Station *Station
		Track   struct {
			Actual, Scheduled string
		}
		Info struct {
			DistanceFromStart int
			Passed            bool
			Distance          int
			Status            int
		}
		Timetable struct {
			ScheduledDepartureTime int
			ArrivalDelay           string
			ScheduledArrivalTime   int
			ActualArrivalTime      int
			DepartureDelay         string
			ActualDepartureTime    int
		}
	}

	if err := json.Unmarshal(b, &parsed); err != nil {
		return err
	}

	*s = Stop{
		Station:              parsed.Station,
		Platform:             parsed.Track.Scheduled,
		DistanceFromStart:    float64(parsed.Info.DistanceFromStart) / 1000.0,
		DistanceFromLastStop: float64(parsed.Info.Distance) / 1000.0,
		Passed:               parsed.Info.Passed,
		ScheduledArrival:     time.Unix(int64(parsed.Timetable.ScheduledArrivalTime/1000), 0),
		ActualArrival:        time.Unix(int64(parsed.Timetable.ActualArrivalTime/1000), 0),
		ScheduledDeparture:   time.Unix(int64(parsed.Timetable.ScheduledDepartureTime/1000), 0),
		ActualDeparture:      time.Unix(int64(parsed.Timetable.ActualDepartureTime/1000), 0),
	}

	if parsed.Track.Actual != "" {
		s.Platform = parsed.Track.Actual
	}

	return nil
}

// Delay returns the estimated delay of arrival for upcoming stops and the
// actual delay of departure for past stops.
func (s Stop) Delay() time.Duration {
	if s.Passed {
		return s.ActualDeparture.Sub(s.ScheduledDeparture)
	}

	return s.ActualArrival.Sub(s.ScheduledArrival)
}

// ETA returns the relative estimated time of arrival, i.e. a duration.
func (s Stop) ETA() time.Duration {
	if s.Passed {
		return time.Duration(0)
	}

	return s.ActualArrival.Sub(time.Now())
}

func (s Stop) String() string {
	return fmt.Sprintf("%v P:%s (%.0fm delay)", s.Station, s.Platform, s.Delay().Minutes())
}

// Trip holds information about the trip. This is modeled after the JSON
// interface and mixes static information (list of stops, train number) with
// volatile live information (distance from last stop, delays).
type Trip struct {
	TrainID              string
	TrainType            string
	Date                 time.Time
	ActualPosition       int
	DistanceFromLastStop float64
	TotalDistance        float64
	Stops                []*Stop
	NextStop             *Stop
	PreviousStop         *Stop
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (t *Trip) UnmarshalJSON(b []byte) error {
	var parsed struct {
		ActualPosition       int
		DistanceFromLastStop int
		VZN                  string
		TrainType            string
		TripDate             string
		StopInfo             struct {
			FinalStationName  string
			ScheduledNext     string
			ActualNext        string
			FinalStationEvaNr string
			ActualLastStarted string
			ActualLast        string
		}
		TotalDistance int
		Stops         []*Stop
	}

	if err := json.Unmarshal(b, &parsed); err != nil {
		return err
	}

	*t = Trip{
		TrainID:              parsed.VZN,
		TrainType:            parsed.TrainType,
		ActualPosition:       parsed.ActualPosition,
		DistanceFromLastStop: float64(parsed.DistanceFromLastStop) / 1000.0,
		TotalDistance:        float64(parsed.TotalDistance) / 1000.0,
		Stops:                parsed.Stops,
	}

	t.Date, _ = time.Parse("2006-01-02", parsed.TripDate)
	t.NextStop = t.findStop(parsed.StopInfo.ActualNext)
	t.PreviousStop = t.findStop(parsed.StopInfo.ActualLast)

	return nil
}

func (t *Trip) findStop(evaNr string) *Stop {
	for _, s := range t.Stops {
		if s.Station.ID == evaNr {
			return s
		}
	}

	return nil
}

// FindStop returns the first Stop in t where the station name matches name.
// This function uses a submatch to find a matching station, so that you can
// specify a city name, e.g. "Basel", and don't have to specify the exact
// station name, e.g. "Basel Bad Bf".
func (t *Trip) FindStop(name string) (*Stop, bool) {
	for _, stop := range t.Stops {
		if strings.Contains(stop.Station.Name, name) {
			return stop, true
		}
	}

	return nil, false
}

// DistanceFromStart returns the distance, in kilometers, from the beginning of the trip.
func (t *Trip) DistanceFromStart() float64 {
	if t.PreviousStop != nil {
		return t.PreviousStop.DistanceFromStart + t.DistanceFromLastStop
	}
	return t.DistanceFromLastStop
}

// DistanceTo returns the distance, in kilometers, between the current position and s.
func (t *Trip) DistanceTo(s *Stop) float64 {
	return s.DistanceFromStart - t.DistanceFromStart()
}

// TripInfo calls the tripInfo API and returns the parsed data.
func TripInfo() (*Trip, error) {
	res, err := http.Get(TripInfoURL)
	if err != nil {
		return nil, err
	}

	var t Trip
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		return nil, err
	}

	return &t, nil
}
