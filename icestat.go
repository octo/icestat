package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"sort"
	"strings"
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

func printTrip(trip *bahn.Trip) error {
	finalStop := trip.Stops[len(trip.Stops)-1]
	nextStop := trip.NextStop
	if nextStop == nil {
		return fmt.Errorf("train arrived in %v", finalStop)
	}

	destinationStop := finalStop
	if *destination != "" {
		var ok bool
		destinationStop, ok = trip.FindStop(*destination)

		if !ok {
			var stops []string
			for _, stop := range trip.Stops {
				stops = append(stops, stop.Station.Name)
			}

			return fmt.Errorf("stop %q not found. Valid stops are: %s",
				*destination, strings.Join(stops, ", "))
		}
	}

	if destinationStop.Passed {
		return fmt.Errorf("train has passed %v", destinationStop)
	}

	if destinationStop != nextStop {
		fmt.Printf("%s%s to %q (via %q): "+
			"distance=%.0f(%.0f) km, "+
			"eta=%s(%s), "+
			"delay=%s(%s)",
			trip.TrainType, trip.TrainID, destinationStop.Station, nextStop.Station,
			trip.DistanceTo(destinationStop), trip.DistanceTo(nextStop),
			formatDuration(destinationStop.ETA()), formatDuration(nextStop.ETA()),
			formatDuration(destinationStop.Delay()), formatDuration(nextStop.Delay()))
	} else {
		fmt.Printf("%s%s to %q: "+
			"distance=%.0f km, "+
			"eta=%s, "+
			"delay=%s",
			trip.TrainType, trip.TrainID, destinationStop.Station,
			trip.DistanceTo(destinationStop),
			formatDuration(destinationStop.ETA()),
			formatDuration(destinationStop.Delay()))
	}

	return nil
}

var speed speedDistribution

func printSpeed(ctx context.Context) error {
	s, err := bahn.StatusInfo(ctx)
	if err != nil {
		return err
	}
	speed.add(s.Speed)

	fmt.Printf(", speed=%.0f/%.0f/%.0f [km/h] (cur/avg/max)",
		s.Speed, speed.average(), speed.max())

	return nil
}

func printUpdate(ctx context.Context) error {
	trip, err := bahn.TripInfo(ctx)
	if err != nil {
		return err
	}

	defer fmt.Println()

	if err := printTrip(trip); err != nil {
		return err
	}

	if err := printSpeed(ctx); err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	for *count != 0 {
		if *count > 0 {
			*count -= 1
		}

		ctx, cancel := context.WithTimeout(ctx, *interval)
		if err := printUpdate(ctx); err != nil {
			cancel()
			log.Println(err)
			time.Sleep(*interval)
			continue
		}
		cancel()

		if *count != 0 {
			time.Sleep(*interval)
		}
	}
}
