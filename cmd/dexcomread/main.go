package main

import (
	"fmt"
	"os"

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
	fmt.Println("Sample Time:", glucoreading.Time.String())
	fmt.Println("Glucose:", glucoreading.Value.String())

}
