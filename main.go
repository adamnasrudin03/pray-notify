package main

import (
	"flag"

	"github.com/adamnasrudin03/pray-notify/scheduler"
)

var testingMode bool

func main() {
	// Pendaftaran flag: -test
	flag.BoolVar(&testingMode, "test", false, "Aktifkan mode uji coba, kirim semua notifikasi ke testWebhookURL.")
	flag.Parse()

	scheduler.RunTesting(testingMode)
	scheduler.Run(testingMode)
}
