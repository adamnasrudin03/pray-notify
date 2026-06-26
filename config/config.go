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
	targets, err := LoadWebhookTargets()
	if err != nil {
		fmt.Printf("❌ ERROR: gagal memuat webhook targets: %v\n", err)
		return
	}
	if len(targets) > 0 {
		WebhookTargets = targets
		return
	}
	fmt.Println("⚠️ WARNING: WebhookTargets is empty. No notifications will be sent. Set WEBHOOK_TARGETS atau WEBHOOK_TARGETS_FILE.")
}

// LoadWebhookTargets memuat target dari WEBHOOK_TARGETS (env) atau WEBHOOK_TARGETS_FILE.
// Env WEBHOOK_TARGETS diutamakan jika keduanya diset.
func LoadWebhookTargets() ([]models.WebhookTarget, error) {
	raw, source, err := readTargetsRaw()
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, nil
	}

	var targets []models.WebhookTarget
	if err := json.Unmarshal(raw, &targets); err != nil {
		return nil, fmt.Errorf("parse %s: %w", source, err)
	}
	if len(targets) == 0 {
		return nil, nil
	}
	return normalizeTargets(targets), nil
}

func readTargetsRaw() ([]byte, string, error) {
	if env := os.Getenv("WEBHOOK_TARGETS"); env != "" {
		return []byte(env), "WEBHOOK_TARGETS", nil
	}
	if file := os.Getenv("WEBHOOK_TARGETS_FILE"); file != "" {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, file, fmt.Errorf("baca file: %w", err)
		}
		return data, file, nil
	}
	return nil, "", nil
}

func normalizeTargets(targets []models.WebhookTarget) []models.WebhookTarget {
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
		if targets[i].MentionMode == "" {
			targets[i].MentionMode = "here"
		}
		if targets[i].SkipWeekends == nil && targets[i].Name == DivisiIT {
			skip := true
			targets[i].SkipWeekends = &skip
		}
		if targets[i].Name != DivisiIT && targets[i].UsernameBeforeAdzan == "" {
			targets[i].UsernameBeforeAdzan = MyUsername
		}
	}
	return targets
}
