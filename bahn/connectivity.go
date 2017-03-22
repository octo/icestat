package bahn // import "github.com/octo/icestat/bahn"

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// ConnectivityURL is the URL of JSONP-encoded information about the train's upstream connections.
const ConnectivityURL = "http://www.ombord.info/api/jsonp/connectivity/"

// DeviceState is the state of a physical connection; either Up or Down.
type DeviceState int

const (
	Down DeviceState = iota
	Up               = iota
)

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (s *DeviceState) UnmarshalText(b []byte) error {
	if string(b) == "up" {
		*s = Up
	} else {
		*s = Down
	}
	return nil
}

// LinkState is the state of a logical connection; either Available (connected)
// or Disconnected.
type LinkState int

const (
	Disconnected LinkState = iota
	Available              = iota
)

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (s *LinkState) UnmarshalText(b []byte) error {
	if string(b) == "available" {
		*s = Available
	} else {
		*s = Disconnected
	}
	return nil
}

type UMTSInfo struct {
	NetStatus *string
	// LAC is the Location Area Code
	LAC    *string
	CellID *string
}

// AccessPointName is the APN of the mobile network.
type AccessPointName struct {
	Name     string
	User     *string
	Password *string
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (apn *AccessPointName) UnmarshalText(b []byte) error {
	fields := strings.Split(string(b), ",")

	apn.Name = fields[0]
	if len(fields) > 1 && fields[1] != "-1" {
		apn.User = &fields[1]
	}
	if len(fields) > 2 && fields[2] != "-1" {
		apn.Password = &fields[2]
	}
	return nil
}

// plmnCodes maps the Mobile Network Codes of the Public Land Mobile Networks
// to their human readable names.
var plmnCodes = map[int]string{
	26201: "T-Mobile",
	26202: "Vodafone",
	26204: "Vodafone",
	26209: "Vodafone",
	26203: "E-plus",
	26205: "E-plus",
	26277: "E-plus",
	26207: "O2",
	26208: "O2",
	26211: "O2",
}

// Link holds information about a single internet uplink.
type Link struct {
	Index         int
	DeviceType    string
	DeviceSubtype string
	DeviceState   DeviceState
	LinkState     LinkState
	RSSI          float64
	Technology    string
	Operator      *string
	APN           *AccessPointName
	UMTS          *UMTSInfo
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (l *Link) UnmarshalJSON(b []byte) error {
	var parsed struct {
		Index         int         `json:",string"`
		DeviceType    string      `json:"device_type"`
		DeviceSubtype string      `json:"device_subtype"`
		DeviceState   DeviceState `json:"device_state"`
		LinkState     LinkState   `json:"link_state"`
		RSSI          float64     `json:",string"`
		Technology    string
		OperatorID    int `json:"operator_id,string"`
		Apninfo       *AccessPointName
		UmtsInfo      struct {
			NetStatus string `json:"net_status"`
			LAC       string
			Cellid    string
		} `json:"umts_info"`
	}
	if err := json.Unmarshal(b, &parsed); err != nil {
		return err
	}

	*l = Link{
		Index:         parsed.Index,
		DeviceType:    parsed.DeviceType,
		DeviceSubtype: parsed.DeviceSubtype,
		DeviceState:   parsed.DeviceState,
		LinkState:     parsed.LinkState,
		RSSI:          parsed.RSSI,
		Technology:    parsed.Technology,
		APN:           parsed.Apninfo,
		UMTS: &UMTSInfo{
			NetStatus: optionalString(parsed.UmtsInfo.NetStatus),
			LAC:       optionalString(parsed.UmtsInfo.LAC),
			CellID:    optionalString(parsed.UmtsInfo.Cellid),
		},
	}

	if parsed.OperatorID > 0 {
		if op, ok := plmnCodes[parsed.OperatorID]; ok {
			l.Operator = &op
		} else {
			op := strconv.Itoa(parsed.OperatorID)
			l.Operator = &op
		}
	}

	return nil
}

// Connectivity holds information about the train's upstream internet
// connections.
type Connectivity struct {
	Version  string
	Online   bool
	BundleID string
	BundleIP net.IP
	Links    []Link
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (c *Connectivity) UnmarshalJSON(b []byte) error {
	var parsed struct {
		Version  string
		Online   int `json:",string"`
		BundleID string
		BundleIP net.IP
		Links    []Link
	}
	if err := json.Unmarshal(b, &parsed); err != nil {
		return err
	}

	*c = Connectivity{
		Version:  parsed.Version,
		Online:   parsed.Online == 1,
		BundleID: parsed.BundleID,
		BundleIP: parsed.BundleIP,
		Links:    parsed.Links,
	}

	return nil
}

// PositionInfo returns information about the train's upstream internet
// connections.
func ConnectivityInfo() (*Connectivity, error) {
	res, err := http.Get(ConnectivityURL)
	if err != nil {
		return nil, err
	}

	var c Connectivity
	if err := unmarshalJSONP(res.Body, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func optionalString(s string) *string {
	if s == "-1" {
		return nil
	}
	return &s
}
