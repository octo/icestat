package bahn

import (
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestPosition_UnmarshalJSON(t *testing.T) {
	input := ioutil.NopCloser(strings.NewReader(`({
    "version":"1.9",
    "time":"1488959212",
    "age":"0",
    "latitude":"48.694882",
    "longitude":"11.45607",
    "altitude":"371.4",
    "speed":"44.206",
    "cmg":"171.07",
    "satellites":"10",
    "mode":"3"
});
`))

	var got Position
	if err := unmarshalJSONP(input, &got); err != nil {
		t.Fatal(err)
	}

	want := Position{
		Version:    "1.9",
		Time:       time.Unix(1488959212, 0),
		Latitude:   48.694882,
		Longitude:  11.45607,
		Altitude:   371.4,
		Speed:      44.206 * 3.6,
		Satellites: 10,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("json.Unmarshal() = %#v, want %#v", got, want)
	}
}
