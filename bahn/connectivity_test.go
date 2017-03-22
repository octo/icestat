package bahn

import (
	"io/ioutil"
	"net"
	"reflect"
	"strings"
	"testing"
)

func TestConnectivity_UnmarshalJSON(t *testing.T) {
	input := ioutil.NopCloser(strings.NewReader(`({
    "version":"1.9",
    "online":"1",
    "bundleid":"84106380353",
    "bundleip":"10.7.19.1",
    "links":[
        {
            "index":"101",
            "device_type":"modem",
            "device_subtype":"mc7304",
            "device_state":"up",
            "link_state":"available",
            "rssi":"-32",
            "technology":"lte",
            "operator_id":"26201",
            "apninfo":"railnet.telekom,t-mobile,tm",
            "umts_info":{
                "net_status":"-1",
                "lac":"-1",
                "cellid":"01B59302"
            }
        },
        {
            "index":"102",
            "device_type":"modem",
            "device_subtype":"mc7304",
            "device_state":"up",
            "link_state":"available",
            "rssi":"-33",
            "technology":"lte",
            "operator_id":"26202",
            "apninfo":"fv1.deutschebahn.com,-1,-1",
            "umts_info":{
                "net_status":"-1",
                "lac":"-1",
                "cellid":"0149C001"
            }
        },
        {
            "index":"103",
            "device_type":"modem",
            "device_subtype":"mc7455",
            "device_state":"up",
            "link_state":"available",
            "rssi":"-79",
            "technology":"dc-hspa+",
            "operator_id":"26201",
            "apninfo":"railnet.telekom,t-mobile,tm",
            "umts_info":{
                "net_status":"-1",
                "lac":"44BC",
                "cellid":"0000D74D"
            }
        },
        {
            "index":"104",
            "device_type":"modem",
            "device_subtype":"mc7455",
            "device_state":"up",
            "link_state":"available",
            "rssi":"-34",
            "technology":"lte",
            "operator_id":"26201",
            "apninfo":"railnet.telekom,t-mobile,tm",
            "umts_info":{
                "net_status":"-1",
                "lac":"-1",
                "cellid":"01B59302"
            }
        },
        {
            "index":"105",
            "device_type":"modem",
            "device_subtype":"mc7455",
            "device_state":"up",
            "link_state":"available",
            "rssi":"-68",
            "technology":"dc-hspa+",
            "operator_id":"26202",
            "apninfo":"fv1.deutschebahn.com,-1,-1",
            "umts_info":{
                "net_status":"-1",
                "lac":"03BB",
                "cellid":"0000AC85"
            }
        },
        {
            "index":"106",
            "device_type":"modem",
            "device_subtype":"mc7304",
            "device_state":"down",
            "link_state":"disconnected",
            "rssi":"-69",
            "technology":"dc-hspa+",
            "operator_id":"26201",
            "apninfo":"railnet.telekom,t-mobile,tm",
            "umts_info":{
                "net_status":"-1",
                "lac":"44B2",
                "cellid":"0000FFFF"
            }
        }
    ]
});
`))

	var got Connectivity
	if err := unmarshalJSONP(input, &got); err != nil {
		t.Fatal(err)
	}

	want := Connectivity{
		Version:  "1.9",
		Online:   true,
		BundleID: "84106380353",
		BundleIP: net.IPv4(10, 7, 19, 1),
		Links: []Link{
			{
				Index:         101,
				DeviceType:    "modem",
				DeviceSubtype: "mc7304",
				DeviceState:   Up,
				LinkState:     Available,
				RSSI:          -32,
				Technology:    "lte",
				Operator:      optionalString("T-Mobile"),
				APN: &AccessPointName{
					Name:     "railnet.telekom",
					User:     optionalString("t-mobile"),
					Password: optionalString("tm"),
				},
				UMTS: &UMTSInfo{
					CellID: optionalString("01B59302"),
				},
			},
			{
				Index:         102,
				DeviceType:    "modem",
				DeviceSubtype: "mc7304",
				DeviceState:   Up,
				LinkState:     Available,
				RSSI:          -33,
				Technology:    "lte",
				Operator:      optionalString("Vodafone"),
				APN: &AccessPointName{
					Name: "fv1.deutschebahn.com",
				},
				UMTS: &UMTSInfo{
					CellID: optionalString("0149C001"),
				},
			},
			{
				Index:         103,
				DeviceType:    "modem",
				DeviceSubtype: "mc7455",
				DeviceState:   Up,
				LinkState:     Available,
				RSSI:          -79,
				Technology:    "dc-hspa+",
				Operator:      optionalString("T-Mobile"),
				APN: &AccessPointName{
					Name:     "railnet.telekom",
					User:     optionalString("t-mobile"),
					Password: optionalString("tm"),
				},
				UMTS: &UMTSInfo{
					LAC:    optionalString("44BC"),
					CellID: optionalString("0000D74D"),
				},
			},
			{
				Index:         104,
				DeviceType:    "modem",
				DeviceSubtype: "mc7455",
				DeviceState:   Up,
				LinkState:     Available,
				RSSI:          -34,
				Technology:    "lte",
				Operator:      optionalString("T-Mobile"),
				APN: &AccessPointName{
					Name:     "railnet.telekom",
					User:     optionalString("t-mobile"),
					Password: optionalString("tm"),
				},
				UMTS: &UMTSInfo{
					CellID: optionalString("01B59302"),
				},
			},
			{
				Index:         105,
				DeviceType:    "modem",
				DeviceSubtype: "mc7455",
				DeviceState:   Up,
				LinkState:     Available,
				RSSI:          -68,
				Technology:    "dc-hspa+",
				Operator:      optionalString("Vodafone"),
				APN: &AccessPointName{
					Name: "fv1.deutschebahn.com",
				},
				UMTS: &UMTSInfo{
					LAC:    optionalString("03BB"),
					CellID: optionalString("0000AC85"),
				},
			},
			{
				Index:         106,
				DeviceType:    "modem",
				DeviceSubtype: "mc7304",
				DeviceState:   Down,
				LinkState:     Disconnected,
				RSSI:          -69,
				Technology:    "dc-hspa+",
				Operator:      optionalString("T-Mobile"),
				APN: &AccessPointName{
					Name:     "railnet.telekom",
					User:     optionalString("t-mobile"),
					Password: optionalString("tm"),
				},
				UMTS: &UMTSInfo{
					LAC:    optionalString("44B2"),
					CellID: optionalString("0000FFFF"),
				},
			},
		},
	}

	for i := range got.Links {
		if !reflect.DeepEqual(got.Links[i], want.Links[i]) {
			t.Errorf("got.Links[%d] != want.Links[%d]\ngot  = %#v\nwant = %#v", i, i, got.Links[i], want.Links[i])
		}
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("json.Unmarshal() = %#v, want %#v", got, want)
	}
}
