package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/adamnasrudin03/pray-notify/scheduler"
)

var (
	testingMode  bool
	prayerWindow string
)

func main() {
	flag.BoolVar(&testingMode, "test", false, "Aktifkan mode uji coba.")
	flag.StringVar(&prayerWindow, "prayer", "", "Mode window: subuh, dhuha, dzuhur, ashar, maghrib, isya")
	flag.Parse()

	scheduler.RunTesting(testingMode && prayerWindow == "")

	if prayerWindow != "" {
		if err := scheduler.RunPrayerWindow(prayerWindow, testingMode); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}
		return
	}

	scheduler.Run(testingMode)
}
