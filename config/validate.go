package config

import (
	"fmt"
	"strings"

	"github.com/adamnasrudin03/pray-notify/models"
)

// ValidateWebhookURL memeriksa format URL Discord webhook.
func ValidateWebhookURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL kosong")
	}
	if strings.Contains(url, "REPLACE_ME") {
		return fmt.Errorf("URL masih placeholder REPLACE_ME — ganti dengan webhook Discord asli di deploy/targets.json")
	}
	const (
		prefix1 = "https://discord.com/api/webhooks/"
		prefix2 = "https://discordapp.com/api/webhooks/"
	)
	if !strings.HasPrefix(url, prefix1) && !strings.HasPrefix(url, prefix2) {
		return fmt.Errorf("URL bukan Discord webhook (harus diawali %s...)", prefix1)
	}
	parts := strings.Split(strings.TrimPrefix(strings.TrimPrefix(url, prefix1), prefix2), "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("URL webhook tidak lengkap (format: .../webhooks/ID/TOKEN)")
	}
	return nil
}

// ValidateTargets memvalidasi semua target sebelum kirim.
func ValidateTargets(targets []models.WebhookTarget) error {
	for _, t := range targets {
		if t.IsDisable {
			continue
		}
		if err := ValidateWebhookURL(t.URL); err != nil {
			return fmt.Errorf("%s: %w", t.Name, err)
		}
	}
	return nil
}
