package bahn // import "github.com/octo/icestat/bahn"

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const inputStr = `
{
   "trip" : {
      "stopInfo" : {
         "actualNext" : "8000261_00",
         "finalStationEvaNr" : "8000261_00",
         "scheduledNext" : "8000261_00",
         "actualLast" : "8000284_00",
         "actualLastStarted" : "8000261",
         "finalStationName" : "München Hbf"
      },
      "trainType" : "ICE",
      "totalDistance" : 503640,
      "distanceFromLastStop" : 136911,
      "actualPosition" : 354328,
      "tripDate" : "2018-08-02",
      "stops" : [
         {
            "timetable" : {
               "showActualArrivalTime" : null,
               "actualArrivalTime" : null,
               "scheduledArrivalTime" : null,
               "departureDelay" : "",
               "scheduledDepartureTime" : 1533176520000,
               "showActualDepartureTime" : true,
               "actualDepartureTime" : 1533176520000,
               "arrivalDelay" : ""
            },
            "delayReasons" : null,
            "info" : {
               "distance" : 0,
               "passed" : true,
               "status" : 0,
               "distanceFromStart" : 0
            },
            "track" : {
               "actual" : "5",
               "scheduled" : "5"
            },
            "station" : {
               "evaNr" : "8000207_00",
               "name" : "Köln Hbf",
               "geocoordinates" : {
                  "latitude" : 50.94303,
                  "longitude" : 6.958729
               }
            }
         },
         {
            "delayReasons" : null,
            "timetable" : {
               "showActualArrivalTime" : true,
               "arrivalDelay" : "",
               "showActualDepartureTime" : true,
               "scheduledDepartureTime" : 1533177360000,
               "actualDepartureTime" : 1533177420000,
               "actualArrivalTime" : 1533177300000,
               "departureDelay" : "+1",
               "scheduledArrivalTime" : 1533177300000
            },
            "info" : {
               "distanceFromStart" : 23857,
               "status" : 0,
               "distance" : 23857,
               "passed" : true
            },
            "track" : {
               "actual" : "6",
               "scheduled" : "6"
            },
            "station" : {
               "geocoordinates" : {
                  "longitude" : 7.203026,
                  "latitude" : 50.793915
               },
               "name" : "Siegburg/Bonn",
               "evaNr" : "8005556_00"
            }
         },
         {
            "station" : {
               "geocoordinates" : {
                  "longitude" : 7.825333,
                  "latitude" : 50.444834
               },
               "evaNr" : "8000667_00",
               "name" : "Montabaur"
            },
            "delayReasons" : null,
            "timetable" : {
               "showActualArrivalTime" : true,
               "actualArrivalTime" : 1533178560000,
               "departureDelay" : "+1",
               "scheduledArrivalTime" : 1533178560000,
               "showActualDepartureTime" : true,
               "scheduledDepartureTime" : 1533178620000,
               "actualDepartureTime" : 1533178680000,
               "arrivalDelay" : ""
            },
            "info" : {
               "distanceFromStart" : 82475,
               "distance" : 58618,
               "status" : 0,
               "passed" : true
            },
            "track" : {
               "scheduled" : "1",
               "actual" : "4"
            }
         },
         {
            "station" : {
               "geocoordinates" : {
                  "latitude" : 50.382498,
                  "longitude" : 8.096112
               },
               "evaNr" : "8003680_00",
               "name" : "Limburg Süd"
            },
            "track" : {
               "scheduled" : "1",
               "actual" : "4"
            },
            "info" : {
               "distanceFromStart" : 102881,
               "distance" : 20406,
               "passed" : true,
               "status" : 0
            },
            "timetable" : {
               "showActualArrivalTime" : true,
               "showActualDepartureTime" : true,
               "scheduledDepartureTime" : 1533179280000,
               "actualDepartureTime" : 1533179280000,
               "actualArrivalTime" : 1533179220000,
               "departureDelay" : "",
               "scheduledArrivalTime" : 1533179220000,
               "arrivalDelay" : ""
            },
            "delayReasons" : null
         },
         {
            "track" : {
               "scheduled" : "Fern 4",
               "actual" : "Fern 4"
            },
            "info" : {
               "distance" : 49801,
               "status" : 0,
               "passed" : true,
               "distanceFromStart" : 152682
            },
            "delayReasons" : null,
            "timetable" : {
               "showActualArrivalTime" : true,
               "arrivalDelay" : "",
               "showActualDepartureTime" : true,
               "scheduledDepartureTime" : 1533180600000,
               "actualDepartureTime" : 1533180600000,
               "actualArrivalTime" : 1533180420000,
               "departureDelay" : "",
               "scheduledArrivalTime" : 1533180420000
            },
            "station" : {
               "evaNr" : "8070003_00",
               "name" : "Frankfurt (M) Flughafen Fernbf",
               "geocoordinates" : {
                  "longitude" : 8.570185,
                  "latitude" : 50.053167
               }
            }
         },
         {
            "station" : {
               "geocoordinates" : {
                  "latitude" : 50.107145,
                  "longitude" : 8.663789
               },
               "evaNr" : "8000105_00",
               "name" : "Frankfurt (Main) Hbf"
            },
            "timetable" : {
               "departureDelay" : "+2",
               "scheduledArrivalTime" : 1533181200000,
               "actualArrivalTime" : 1533181320000,
               "actualDepartureTime" : 1533182160000,
               "scheduledDepartureTime" : 1533182040000,
               "showActualDepartureTime" : true,
               "arrivalDelay" : "+2",
               "showActualArrivalTime" : true
            },
            "delayReasons" : null,
            "track" : {
               "actual" : "7",
               "scheduled" : "7"
            },
            "info" : {
               "distanceFromStart" : 161664,
               "distance" : 8982,
               "status" : 0,
               "passed" : true
            }
         },
         {
            "station" : {
               "geocoordinates" : {
                  "latitude" : 50.120953,
                  "longitude" : 8.929
               },
               "name" : "Hanau Hbf",
               "evaNr" : "8000150_00"
            },
            "info" : {
               "distanceFromStart" : 180642,
               "passed" : true,
               "distance" : 18978,
               "status" : 0
            },
            "track" : {
               "actual" : "103",
               "scheduled" : "103"
            },
            "delayReasons" : null,
            "timetable" : {
               "scheduledDepartureTime" : 1533183000000,
               "showActualDepartureTime" : true,
               "actualDepartureTime" : 1533183180000,
               "actualArrivalTime" : 1533183120000,
               "departureDelay" : "+3",
               "scheduledArrivalTime" : 1533182940000,
               "arrivalDelay" : "+3",
               "showActualArrivalTime" : true
            }
         },
         {
            "station" : {
               "geocoordinates" : {
                  "latitude" : 49.980557,
                  "longitude" : 9.143697
               },
               "evaNr" : "8000010_00",
               "name" : "Aschaffenburg Hbf"
            },
            "timetable" : {
               "showActualArrivalTime" : true,
               "departureDelay" : "+3",
               "scheduledArrivalTime" : 1533183780000,
               "actualArrivalTime" : 1533183900000,
               "actualDepartureTime" : 1533184020000,
               "showActualDepartureTime" : true,
               "scheduledDepartureTime" : 1533183840000,
               "arrivalDelay" : "+2"
            },
            "delayReasons" : null,
            "track" : {
               "actual" : "6",
               "scheduled" : "6"
            },
            "info" : {
               "distance" : 21885,
               "passed" : true,
               "status" : 0,
               "distanceFromStart" : 202527
            }
         },
         {
            "station" : {
               "geocoordinates" : {
                  "latitude" : 49.801796,
                  "longitude" : 9.93578
               },
               "evaNr" : "8000260_00",
               "name" : "Würzburg Hbf"
            },
            "delayReasons" : null,
            "timetable" : {
               "arrivalDelay" : "+1",
               "scheduledArrivalTime" : 1533186120000,
               "departureDelay" : "+1",
               "actualArrivalTime" : 1533186180000,
               "actualDepartureTime" : 1533186360000,
               "showActualDepartureTime" : true,
               "scheduledDepartureTime" : 1533186300000,
               "showActualArrivalTime" : true
            },
            "info" : {
               "status" : 0,
               "distance" : 60139,
               "passed" : true,
               "distanceFromStart" : 262666
            },
            "track" : {
               "scheduled" : "5",
               "actual" : "5"
            }
         },
         {
            "timetable" : {
               "showActualArrivalTime" : true,
               "arrivalDelay" : "",
               "actualDepartureTime" : 1533189720000,
               "scheduledDepartureTime" : 1533189720000,
               "showActualDepartureTime" : true,
               "departureDelay" : "",
               "scheduledArrivalTime" : 1533189540000,
               "actualArrivalTime" : 1533189540000
            },
            "delayReasons" : null,
            "info" : {
               "passed" : true,
               "distance" : 91662,
               "status" : 0,
               "distanceFromStart" : 354328
            },
            "track" : {
               "scheduled" : "9",
               "actual" : "8"
            },
            "station" : {
               "evaNr" : "8000284_00",
               "name" : "Nürnberg Hbf",
               "geocoordinates" : {
                  "latitude" : 49.445616,
                  "longitude" : 11.082989
               }
            }
         },
         {
            "delayReasons" : null,
            "timetable" : {
               "showActualDepartureTime" : null,
               "scheduledDepartureTime" : null,
               "actualDepartureTime" : null,
               "actualArrivalTime" : 1533193620000,
               "departureDelay" : "",
               "scheduledArrivalTime" : 1533193620000,
               "arrivalDelay" : "",
               "showActualArrivalTime" : true
            },
            "track" : {
               "actual" : "26",
               "scheduled" : "23"
            },
            "info" : {
               "passed" : false,
               "distance" : 149312,
               "status" : 0,
               "distanceFromStart" : 503640
            },
            "station" : {
               "geocoordinates" : {
                  "latitude" : 48.140232,
                  "longitude" : 11.558335
               },
               "evaNr" : "8000261_00",
               "name" : "München Hbf"
            }
         }
      ],
      "vzn" : "521"
   },
   "selectedRoute" : {
      "mobility" : null,
      "conflictInfo" : {
         "text" : null,
         "status" : "NO_CONFLICT"
      }
   },
   "connection" : null
}
`

func TestTrip(t *testing.T) {
	var trip Trip
	if err := json.Unmarshal([]byte(inputStr), &trip); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	// Distances, returned as fractional km, are precise up to 1 meter.
	approxFloat := cmpopts.EquateApprox(0, .001)

	if got, want := trip.DistanceFromLastStop, 136.911; !cmp.Equal(got, want, approxFloat) {
		t.Errorf("trip.DistanceFromLastStop = %g, want %g", got, want)
	}

	if got, want := trip.TotalDistance, 503.640; got != want {
		t.Errorf("trip.TotalDistance = %g, want %g", got, want)
	}

	if got, want := len(trip.Stops), 11; got != want {
		t.Fatalf("len(trip.Stops) = %d, want %d", got, want)
	}

	if got, want := trip.Date, time.Date(2018, 8, 2, 0, 0, 0, 0, time.UTC); !got.Equal(want) {
		t.Errorf("trip.Date = %v, want %v", got, want)
	}

	if got, want := trip.DistanceFromStart(), 354.328+136.911; !cmp.Equal(got, want, approxFloat) {
		t.Errorf("DistanceFromStart() = %g, want %g", got, want)
	}

	if got, want := trip.DistanceTo(trip.Stops[10]), 503.640-(354.328+136.911); !cmp.Equal(got, want, approxFloat) {
		t.Errorf("trip.DistanceTo(%v) = %g, want %g", trip.Stops[10], got, want)
	}
}
