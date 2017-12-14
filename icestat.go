package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
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

			return fmt.Errorf("stop %q not found. Valid stops are: ",
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

func formatRSSI(rssi float64) string {
	symbols := []string{"█", "▇", "▆", "▅", "▄", "▃", "▂", "▁"}

	// TODO(octo): these levels assume 3G/HSPA and should be adapted for 4G/LTE.
	lowerBounds := []float64{-67.5, -75.0, -82.5, -90.0, -95.0, -100.0, -105.0}

	for i, lb := range lowerBounds {
		if rssi >= lb {
			return symbols[i]
		}
	}
	return symbols[len(symbols)-1]
}

func printConnectivity() error {
	c, err := bahn.ConnectivityInfo()
	if err != nil {
		return err
	}

	state := "offline"
	if c.Online {
		state = "online"
	}

	var linksUp int
	var rssiIndicators []string
	for _, link := range c.Links {
		if link.Up() {
			linksUp++
			rssiIndicators = append(rssiIndicators, formatRSSI(link.RSSI))
		} else {
			rssiIndicators = append(rssiIndicators, " ")
		}
	}

	fmt.Printf(", wifi=%s [%s] (%d/%d)", state, strings.Join(rssiIndicators, ""), linksUp, len(c.Links))
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

	if err := printConnectivity(); err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	// Patch net/http's DefaultTransport to ignore invalid TLS certificates.
	// This is required to work around Deutsche Bahn's broken TLS setup
	// (they're serving their on-train content with the bahn.de certificate.
	// TODO(octo): Remove once DB fixes their TLS setup.
	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		if t.TLSClientConfig == nil {
			t.TLSClientConfig = &tls.Config{}
		}
		t.TLSClientConfig.InsecureSkipVerify = true
		log.Print("disabled TLS certificate verification")
	}

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
