package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/adamnasrudin03/pray-notify/models"
)

const (
	MyUsername = "Adam Nasrudin (@adamnasrudin)"
	Test       = "Testing"
	DivisiIT   = "Divisi IT"
	MyChannel  = "My Channel"
)

// WebhookTargets adalah daftar target webhook dengan lokasi spesifik
var WebhookTargets = []models.WebhookTarget{}

func init() {
	// Overwrite WebhookTargets if set by environment variable (JSON string)
	// Example: [{"Name":"Custom","URL":"https://...","Lat":-6.175,"Lon":106.827,"Tz":"Asia/Jakarta","City":"Jakarta"}]
	if env := os.Getenv("WEBHOOK_TARGETS"); env != "" {
		var targets []models.WebhookTarget
		if err := json.Unmarshal([]byte(env), &targets); err != nil {
			fmt.Printf("❌ ERROR: Failed to parse WEBHOOK_TARGETS: %v\n", err)
		} else if len(targets) > 0 {
			for i := range targets {
				if targets[i].Lat == 0 {
					targets[i].Lat = models.DefaultLat
				}
				if targets[i].Lon == 0 {
					targets[i].Lon = models.DefaultLon
				}
				if targets[i].Tz == "" {
					targets[i].Tz = models.DefaultTz
				}
				if targets[i].City == "" {
					targets[i].City = models.DefaultCity
				}
			}
			WebhookTargets = targets
		}
	}

	if len(WebhookTargets) == 0 {
		fmt.Println("⚠️ WARNING: WebhookTargets is empty. No notifications will be sent. Please check your WEBHOOK_TARGETS environment variable.")
	}
}

