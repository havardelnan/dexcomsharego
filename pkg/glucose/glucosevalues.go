package glucose

import (
	"fmt"
	"maps"
	"math"
	"slices"
	"time"
)

type GlucoseReadings []GlucoseReading

type GlucoseReading struct {
	Time  time.Time
	Value GlucoseValue
	Trend string
}

type GlucoseValue float64

// String returns the string representation of the GlucoseValue in mmol/L
func MgDlTommolL(val int) float64 {
	mmol := float64(val) * 0.0555
	return (math.Round(mmol*10) / 10)
}

func (v GlucoseValue) String() string {
	return fmt.Sprintf("%.1f mmol/L", v)
}

type GlucoseStats struct {
	Samples int
	Min     GlucoseValue
	Max     GlucoseValue
	Avg     GlucoseValue
}

func (s GlucoseStats) String() string {
	return fmt.Sprintf("Samples: %d, Min: %s, Max: %s, Avg: %s\n", s.Samples, s.Min, s.Max, s.Avg)
}

func (s GlucoseStats) Print() {
	fmt.Print(s.String())
}

func (r GlucoseReadings) Stats() GlucoseStats {
	return GetDetailedStats(r)
}

func GetDetailedStats(r GlucoseReadings) GlucoseStats {
	var min, max, avg GlucoseValue
	min = 1000
	max = 0
	avg = 0
	for _, reading := range r {
		if reading.Value < min {
			min = reading.Value
		}
		if reading.Value > max {
			max = reading.Value
		}
		avg += reading.Value
	}
	avg = GlucoseValue(float64(avg) / float64(len(r)))
	return GlucoseStats{Samples: len(r), Min: min, Max: max, Avg: avg}
}

type HourlyStats map[int64]GlucoseReadings

func (h HourlyStats) Print() {

	fmt.Print("Last 24 hours:\n")
	fmt.Print("Hour\t\tMin\t\tMax\t\tAvg\n")
	sortedhours := slices.Sorted(maps.Keys(h))
	r := GlucoseReadings{}
	for _, hour := range sortedhours {
		r = append(r, h[hour]...)
		fmt.Printf("%s\t%s\t%s\t%s\n", h[hour][0].Time.Format("02.01 15:00"), h[hour].Stats().Min, h[hour].Stats().Max, h[hour].Stats().Avg)
	}
	fmt.Println("")

	GetDetailedStats(r).Print()
}

func GetHourlyStats(r GlucoseReadings) HourlyStats {
	hours := make(map[int64]GlucoseReadings)
	for _, reading := range r {
		timebyhour, _ := time.Parse("2006.01.02 15:00", reading.Time.Format("2006.01.02 15:00"))
		hours[timebyhour.Unix()] = append(hours[timebyhour.Unix()], reading)
	}
	return hours
}
