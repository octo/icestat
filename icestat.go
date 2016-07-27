package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"time"

	"github.com/octo/icestat/bahn"
)

var (
	interval    = flag.Duration("interval", 10*time.Second, "Interval in which to report statistics.")
	count       = flag.Int("count", -1, "Number of iterations.")
	destination = flag.String("destination", "", "Optional destination to anticipate.")
)

type speedDistribution struct {
	data []float64
}

func (s *speedDistribution) add(kmh float64) {
	s.data = append(s.data, kmh)
	sort.Float64s(s.data)
}

func (s speedDistribution) max() float64 {
	if len(s.data) == 0 {
		return math.NaN()
	}

	return s.data[len(s.data)-1]
}

func (s speedDistribution) percentile(p float64) float64 {
	if len(s.data) == 0 {
		return math.NaN()
	}

	idx := float64(len(s.data)) * p / 100.0
	return s.data[int(idx-0.5)]
}

func (s speedDistribution) median() float64 {
	return s.percentile(50.0)
}

func (s speedDistribution) average() float64 {
	var sum float64
	var num int

	for _, speed := range s.data {
		if !math.IsNaN(speed) {
			sum += speed
			num++
		}
	}

	if num == 0 {
		return math.NaN()
	}

	return sum / float64(num)
}

func formatDuration(d time.Duration) string {
	h, m := math.Modf(d.Hours())
	m *= 60.0

	return fmt.Sprintf("%.0f:%02.0f", h, m)
}

func main() {
	var sd speedDistribution

	flag.Parse()

	for *count != 0 {
		if *count > 0 {
			*count -= 1
		}

		trip, err := bahn.TripInfo()
		if err != nil {
			log.Print("TripInfo: ", err)
			continue
		}

		s, err := bahn.StatusInfo()
		if err != nil {
			log.Print("StatusInfo: ", err)
			continue
		}

		sd.add(s.Speed)

		nextStop := trip.NextStop
		finalStop := trip.Stops[len(trip.Stops)-1]
		if *destination != "" {
			var ok bool
			finalStop, ok = trip.FindStop(*destination)

			if !ok {
				log.Printf("stop %q not found. Valid stops are:", *destination)
				for _, stop := range trip.Stops {
					log.Printf("  * %q", stop.Station)
				}
				os.Exit(1)
			}
		}

		if finalStop.Passed {
			fmt.Printf("Train has passed %v\n", finalStop)
			*count = 0
			continue
		}

		if nextStop == nil {
			fmt.Printf("Train arrived in %v\n", trip.Stops[len(trip.Stops)-1])
			*count = 0
			continue
		}

		if nextStop == nil || finalStop == nil {
			continue
		}

		fmt.Printf("%s%s to %q (via %q): "+
			"distance=%.0f(%.0f) km, "+
			"eta=%s(%s), "+
			"delay=%s(%s), "+
			"speed=%.0f/%.0f/%.0f km/h (cur/avg/max)\n",
			trip.TrainType, trip.TrainID, finalStop.Station, nextStop.Station,
			trip.DistanceTo(finalStop), trip.DistanceTo(nextStop),
			formatDuration(finalStop.ETA()), formatDuration(nextStop.ETA()),
			formatDuration(finalStop.Delay()), formatDuration(nextStop.Delay()),
			s.Speed, sd.average(), sd.max())

		if *count != 0 {
			time.Sleep(*interval)
		}
	}
}
