package bahn // import "github.com/octo/icestat/bahn"

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// StatusURL is the URL of JSON encoded information about the train's location and speed.
const StatusURL = "https://iceportal.de/api1/rs/status"

// Status holds the information returned by the status API call.
type Status struct {
	Connection   bool
	ServiceLevel string
	Speed        float64
	Longitude    float64
	Latitude     float64
	ServerTime   time.Time
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (s *Status) UnmarshalJSON(b []byte) error {
	var parsed struct {
		Connection   bool
		ServiceLevel string
		Speed        float64
		Longitude    float64
		Latitude     float64
		ServerTime   int
	}
	if err := json.Unmarshal(b, &parsed); err != nil {
		return err
	}

	*s = Status{
		Connection:   parsed.Connection,
		ServiceLevel: parsed.ServiceLevel,
		Speed:        parsed.Speed,
		Longitude:    parsed.Longitude,
		Latitude:     parsed.Latitude,
		ServerTime:   time.Unix(int64(parsed.ServerTime/1000), 0),
	}

	return nil
}

// StatusInfo calls the status API and returns the parsed data.
func StatusInfo(ctx context.Context) (*Status, error) {
	req, err := http.NewRequest(http.MethodGet, StatusURL, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var s Status
	if err := json.NewDecoder(res.Body).Decode(&s); err != nil {
		return nil, err
	}

	return &s, nil
}
