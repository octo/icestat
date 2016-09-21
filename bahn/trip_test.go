package bahn // import "github.com/octo/icestat/bahn"

import (
	"encoding/json"
	"testing"
	"time"
)

const inputStr = `
{
   "actualPosition" : 0,
   "vzn" : "1522",
   "stopInfo" : {
      "finalStationName" : "Dortmund Hbf",
      "scheduledNext" : "8000183_00",
      "actualNext" : "8000183_00",
      "finalStationEvaNr" : "8000080_00",
      "actualLastStarted" : "8000183",
      "actualLast" : "8000261_00"
   },
   "distanceFromLastStop" : 1318,
   "totalDistance" : 595608,
   "trainType" : "ICE",
   "tripDate" : "2016-04-05",
   "stops" : [
      {
         "station" : {
            "evaNr" : "8000261_00",
            "geocoordinates" : {
               "latitude" : 48.140232,
               "longitude" : 11.558335
            },
            "name" : "München Hbf"
         },
         "track" : {
            "actual" : "",
            "scheduled" : "18"
         },
         "info" : {
            "distanceFromStart" : 0,
            "passed" : true,
            "distance" : 0,
            "status" : 0
         },
         "timetable" : {
            "scheduledDepartureTime" : 1459779480000,
            "arrivalDelay" : "",
            "scheduledArrivalTime" : null,
            "actualArrivalTime" : null,
            "departureDelay" : "+1",
            "actualDepartureTime" : 1459779540000
         }
      },
      {
         "station" : {
            "name" : "Ingolstadt Hbf",
            "geocoordinates" : {
               "latitude" : 48.744541,
               "longitude" : 11.437337
            },
            "evaNr" : "8000183_00"
         },
         "track" : {
            "actual" : "",
            "scheduled" : "4"
         },
         "info" : {
            "status" : 0,
            "distance" : 67805,
            "distanceFromStart" : 67805,
            "passed" : false
         },
         "timetable" : {
            "actualDepartureTime" : 1459781880000,
            "departureDelay" : "+1",
            "actualArrivalTime" : 1459781760000,
            "scheduledArrivalTime" : 1459781700000,
            "arrivalDelay" : "+1",
            "scheduledDepartureTime" : 1459781820000
         }
      },
      {
         "station" : {
            "name" : "Nürnberg Hbf",
            "evaNr" : "8000284_00",
            "geocoordinates" : {
               "latitude" : 49.445616,
               "longitude" : 11.082989
            }
         },
         "track" : {
            "actual" : "8",
            "scheduled" : "6"
         },
         "info" : {
            "passed" : false,
            "distanceFromStart" : 149942,
            "status" : 0,
            "distance" : 82137
         },
         "timetable" : {
            "arrivalDelay" : "+1",
            "scheduledDepartureTime" : 1459784520000,
            "actualArrivalTime" : 1459784340000,
            "scheduledArrivalTime" : 1459784280000,
            "actualDepartureTime" : 1459784520000,
            "departureDelay" : ""
         }
      },
      {
         "timetable" : {
            "arrivalDelay" : "",
            "scheduledDepartureTime" : 1459787700000,
            "actualArrivalTime" : 1459787460000,
            "scheduledArrivalTime" : 1459787460000,
            "actualDepartureTime" : 1459787700000,
            "departureDelay" : ""
         },
         "station" : {
            "evaNr" : "8000260_00",
            "geocoordinates" : {
               "longitude" : 9.93578,
               "latitude" : 49.801796
            },
            "name" : "Würzburg Hbf"
         },
         "info" : {
            "status" : 0,
            "distance" : 91662,
            "distanceFromStart" : 241604,
            "passed" : false
         },
         "track" : {
            "actual" : "",
            "scheduled" : "6"
         }
      },
      {
         "station" : {
            "geocoordinates" : {
               "longitude" : 8.929,
               "latitude" : 50.120953
            },
            "evaNr" : "8000150_00",
            "name" : "Hanau Hbf"
         },
         "track" : {
            "scheduled" : "102",
            "actual" : ""
         },
         "info" : {
            "distanceFromStart" : 321912,
            "passed" : false,
            "status" : 0,
            "distance" : 80308
         },
         "timetable" : {
            "actualDepartureTime" : 1459790640000,
            "departureDelay" : "",
            "actualArrivalTime" : 1459790520000,
            "scheduledArrivalTime" : 1459790520000,
            "arrivalDelay" : "",
            "scheduledDepartureTime" : 1459790640000
         }
      },
      {
         "station" : {
            "geocoordinates" : {
               "latitude" : 50.107145,
               "longitude" : 8.663789
            },
            "evaNr" : "8000105_00",
            "name" : "Frankfurt(Main)Hbf"
         },
         "info" : {
            "distance" : 18978,
            "status" : 0,
            "passed" : false,
            "distanceFromStart" : 340890
         },
         "track" : {
            "actual" : "",
            "scheduled" : "6"
         },
         "timetable" : {
            "departureDelay" : "",
            "actualDepartureTime" : 1459791900000,
            "scheduledArrivalTime" : 1459791600000,
            "actualArrivalTime" : 1459791600000,
            "scheduledDepartureTime" : 1459791900000,
            "arrivalDelay" : ""
         }
      },
      {
         "timetable" : {
            "actualDepartureTime" : 1459792680000,
            "departureDelay" : "",
            "actualArrivalTime" : 1459792560000,
            "scheduledArrivalTime" : 1459792560000,
            "arrivalDelay" : "",
            "scheduledDepartureTime" : 1459792680000
         },
         "info" : {
            "distanceFromStart" : 349872,
            "passed" : false,
            "distance" : 8982,
            "status" : 0
         },
         "track" : {
            "actual" : "",
            "scheduled" : "Fern 7"
         },
         "station" : {
            "name" : "Frankfurt(M) Flughafen Fernbf",
            "geocoordinates" : {
               "longitude" : 8.570185,
               "latitude" : 50.053167
            },
            "evaNr" : "8070003_00"
         }
      },
      {
         "timetable" : {
            "arrivalDelay" : "",
            "scheduledDepartureTime" : 1459794000000,
            "actualArrivalTime" : 1459793880000,
            "scheduledArrivalTime" : 1459793880000,
            "actualDepartureTime" : 1459794000000,
            "departureDelay" : ""
         },
         "station" : {
            "name" : "Mainz Hbf",
            "geocoordinates" : {
               "longitude" : 8.25872,
               "latitude" : 50.001117
            },
            "evaNr" : "8000240_00"
         },
         "info" : {
            "passed" : false,
            "distanceFromStart" : 372868,
            "distance" : 22996,
            "status" : 0
         },
         "track" : {
            "scheduled" : "3",
            "actual" : ""
         }
      },
      {
         "timetable" : {
            "actualArrivalTime" : 1459797060000,
            "scheduledArrivalTime" : 1459797060000,
            "arrivalDelay" : "",
            "scheduledDepartureTime" : 1459797180000,
            "actualDepartureTime" : 1459797180000,
            "departureDelay" : ""
         },
         "station" : {
            "evaNr" : "8000206_00",
            "geocoordinates" : {
               "latitude" : 50.350929,
               "longitude" : 7.588345
            },
            "name" : "Koblenz Hbf"
         },
         "track" : {
            "scheduled" : "3",
            "actual" : ""
         },
         "info" : {
            "status" : 0,
            "distance" : 61596,
            "passed" : false,
            "distanceFromStart" : 434464
         }
      },
      {
         "timetable" : {
            "arrivalDelay" : "",
            "scheduledDepartureTime" : 1459799100000,
            "actualArrivalTime" : 1459798980000,
            "scheduledArrivalTime" : 1459798980000,
            "actualDepartureTime" : 1459799100000,
            "departureDelay" : ""
         },
         "info" : {
            "passed" : false,
            "distanceFromStart" : 489256,
            "distance" : 54792,
            "status" : 0
         },
         "track" : {
            "actual" : "",
            "scheduled" : "1"
         },
         "station" : {
            "name" : "Bonn Hbf",
            "evaNr" : "8000044_00",
            "geocoordinates" : {
               "longitude" : 7.097136,
               "latitude" : 50.732008
            }
         }
      },
      {
         "track" : {
            "actual" : "",
            "scheduled" : "4"
         },
         "info" : {
            "distanceFromStart" : 514661,
            "passed" : false,
            "distance" : 25405,
            "status" : 0
         },
         "station" : {
            "name" : "Köln Hbf",
            "evaNr" : "8000207_00",
            "geocoordinates" : {
               "latitude" : 50.94303,
               "longitude" : 6.958729
            }
         },
         "timetable" : {
            "departureDelay" : "",
            "actualDepartureTime" : 1459800600000,
            "scheduledArrivalTime" : 1459800300000,
            "actualArrivalTime" : 1459800300000,
            "scheduledDepartureTime" : 1459800600000,
            "arrivalDelay" : ""
         }
      },
      {
         "info" : {
            "distance" : 24426,
            "status" : 0,
            "passed" : false,
            "distanceFromStart" : 539087
         },
         "track" : {
            "actual" : "",
            "scheduled" : "3"
         },
         "station" : {
            "evaNr" : "8000087_00",
            "geocoordinates" : {
               "latitude" : 51.160766,
               "longitude" : 7.004187
            },
            "name" : "Solingen Hbf"
         },
         "timetable" : {
            "actualDepartureTime" : 1459801740000,
            "departureDelay" : "",
            "arrivalDelay" : "",
            "scheduledDepartureTime" : 1459801740000,
            "actualArrivalTime" : 1459801620000,
            "scheduledArrivalTime" : 1459801620000
         }
      },
      {
         "timetable" : {
            "arrivalDelay" : "",
            "scheduledDepartureTime" : 1459802580000,
            "actualArrivalTime" : 1459802460000,
            "scheduledArrivalTime" : 1459802460000,
            "actualDepartureTime" : 1459802580000,
            "departureDelay" : ""
         },
         "track" : {
            "scheduled" : "2",
            "actual" : ""
         },
         "info" : {
            "status" : 0,
            "distance" : 14525,
            "distanceFromStart" : 553612,
            "passed" : false
         },
         "station" : {
            "name" : "Wuppertal Hbf",
            "geocoordinates" : {
               "latitude" : 51.254363,
               "longitude" : 7.149543
            },
            "evaNr" : "8000266_00"
         }
      },
      {
         "timetable" : {
            "actualArrivalTime" : 1459803600000,
            "scheduledArrivalTime" : 1459803600000,
            "arrivalDelay" : "",
            "scheduledDepartureTime" : 1459803720000,
            "actualDepartureTime" : 1459803720000,
            "departureDelay" : ""
         },
         "track" : {
            "actual" : "",
            "scheduled" : "6"
         },
         "info" : {
            "status" : 0,
            "distance" : 24739,
            "passed" : false,
            "distanceFromStart" : 578351
         },
         "station" : {
            "evaNr" : "8000142_00",
            "geocoordinates" : {
               "latitude" : 51.362747,
               "longitude" : 7.460249
            },
            "name" : "Hagen Hbf"
         }
      },
      {
         "track" : {
            "actual" : "",
            "scheduled" : "10"
         },
         "info" : {
            "distanceFromStart" : 595608,
            "passed" : false,
            "distance" : 17257,
            "status" : 0
         },
         "station" : {
            "evaNr" : "8000080_00",
            "geocoordinates" : {
               "longitude" : 7.45929,
               "latitude" : 51.517896
            },
            "name" : "Dortmund Hbf"
         },
         "timetable" : {
            "actualArrivalTime" : 1459804920000,
            "scheduledArrivalTime" : 1459804920000,
            "arrivalDelay" : "",
            "scheduledDepartureTime" : null,
            "actualDepartureTime" : null,
            "departureDelay" : ""
         }
      }
   ]
}
`

func TestTrip(t *testing.T) {
	var trip Trip
	if err := json.Unmarshal([]byte(inputStr), &trip); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if trip.DistanceFromLastStop != 1.318 {
		t.Errorf("trip.DistanceFromLastStop = %g, want %g", trip.DistanceFromLastStop, 1.318)
	}

	if trip.TotalDistance != 595.608 {
		t.Errorf("trip.TotalDistance = %g, want %g", trip.TotalDistance, 595.608)
	}

	if len(trip.Stops) != 15 {
		t.Errorf("len(trip.Stops) = %d, want %d", len(trip.Stops), 15)
	}

	wantDate := time.Date(2016, 4, 5, 0, 0, 0, 0, time.UTC)
	if !trip.Date.Equal(wantDate) {
		t.Errorf("trip.Date = %v, want %v", trip.Date, wantDate)
	}

	if got := trip.DistanceFromStart(); got != 1.318 {
		t.Errorf("DistanceFromStart() = %g, want %g", got, 1.318)
	}

	if got := trip.DistanceTo(trip.Stops[2]); got != 148.624 {
		t.Errorf("trip.DistanceTo(%v) = %g, want %g", trip.Stops[2], got, 148.624)
	}
}

func TestStop(t *testing.T) {
	var trip Trip
	if err := json.Unmarshal([]byte(inputStr), &trip); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if trip.Stops[1].Platform != "4" {
		t.Errorf("trip.Stops[1].Platform = %q, want %q", trip.Stops[1].Platform, "4")
	}
	// "actual" overrides "scheduled" if both are given.
	if trip.Stops[2].Platform != "8" {
		t.Errorf("trip.Stops[2].Platform = %q, want %q", trip.Stops[1].Platform, "8")
	}

	// when passed == true, departure delay is used
	if got := trip.Stops[0].Delay(); got != time.Minute {
		t.Errorf("trip.Stops[0].Delay() = %v, want %v", got, time.Minute)
	}
	// else, arrival delay is used
	if got := trip.Stops[2].Delay(); got != time.Minute {
		t.Errorf("trip.Stops[2].Delay() = %v, want %v", got, time.Minute)
	}
}
