package glucose

import (
	"fmt"
	"math"
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
