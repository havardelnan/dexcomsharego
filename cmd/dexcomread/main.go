package main

import (
	"fmt"
	"os"

	"github.com/havardelnan/dexcomsharego/pkg/glucose"
	"github.com/havardelnan/dexcomsharego/pkg/shareclient"
)

func main() {
	shareclient, err := shareclient.NewSharesession(shareclient.ShareAuthConfig{
		ApplicationId: os.Getenv("DEXCOM_APPLICATION_ID"),
		Username:      os.Getenv("DEXCOM_USERNAME"),
		Password:      os.Getenv("DEXCOM_PASSWORD"),
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	glucoreading := shareclient.GetGlucoseReading()
	fmt.Println("Last Sample:", glucoreading.Time.Format("15:04 02.01.2006"))
	fmt.Println("Glucose:", glucoreading.Value.String())

	readings := shareclient.GetGlucoseReadings(60, 15)
	fmt.Print("\nLast hour:\n")
	fmt.Print("Time\tValue\t\tTrend\n")
	for _, reading := range readings {
		fmt.Printf("%s\t%s\t%s\n", reading.Time.Format("15:04:05"), reading.Value.String(), reading.Trend)
	}

	readings.Stats().Print()
	fmt.Println("")
	glucose.GetHourlyStats(shareclient.GetGlucoseReadings(1440, 288)).Print()
}
