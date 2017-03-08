package bahn // import "github.com/octo/icestat/bahn"

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// PositionURL is the URL of JSONP-encoded information about the train's location and speed.
const PositionURL = "http://www.ombord.info/api/jsonp/position/"

// Position holds the position of the train as well as some other GPS related
// data (speed, #satellites, ...).
type Position struct {
	Version   string
	Time      time.Time
	Latitude  float64
	Longitude float64
	Altitude  float64
	// Speed in km/h.
	Speed      float64
	Satellites int
}

// quotedFloat64 helps parsing quoted strings as floats.
type quotedFloat64 float64

func (f *quotedFloat64) UnmarshalText(text []byte) error {
	parsed, err := strconv.ParseFloat(string(text), 64)
	if err != nil {
		return err
	}

	*f = quotedFloat64(parsed)
	return nil
}

// quotedInt helps parsing quoted strings as ints.
type quotedInt int

func (i *quotedInt) UnmarshalText(text []byte) error {
	parsed, err := strconv.ParseInt(string(text), 0, 0)
	if err != nil {
		return err
	}

	*i = quotedInt(parsed)
	return nil
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (p *Position) UnmarshalJSON(b []byte) error {
	var parsed struct {
		Version   string
		Time      quotedInt
		Age       string
		Latitude  quotedFloat64
		Longitude quotedFloat64
		Altitude  quotedFloat64
		// Speed in m/s.
		Speed quotedFloat64
		// Gyroskop?
		CMG        string
		Satellites quotedInt
		Mode       string
	}
	if err := json.Unmarshal(b, &parsed); err != nil {
		return err
	}

	*p = Position{
		Version:   parsed.Version,
		Time:      time.Unix(int64(parsed.Time), 0),
		Latitude:  float64(parsed.Latitude),
		Longitude: float64(parsed.Longitude),
		Altitude:  float64(parsed.Altitude),
		// convert speed from m/s to km/h
		Speed:      float64(parsed.Speed) * 3600.0 / 1000.0,
		Satellites: int(parsed.Satellites),
	}
	return nil
}

// PositionInfo returns information about the train's position and speed.
func PositionInfo() (*Position, error) {
	res, err := http.Get(PositionURL)
	if err != nil {
		return nil, err
	}

	var p Position
	if err := unmarshalJSONP(res.Body, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
