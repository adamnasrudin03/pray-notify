package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/adamnasrudin03/pray-notify/config"
	"github.com/adamnasrudin03/pray-notify/models"
	"github.com/adamnasrudin03/pray-notify/reminders"
)

// WebhookResult menyimpan hasil pengiriman ke satu webhook
type WebhookResult struct {
	Name  string
	Error error
}

// BuildWebhookContent menyusun konten teks pesan webhook.
func BuildWebhookContent(target models.WebhookTarget, label string, categoryShalat string) string {
	message := target.ContentPrefix(label)
	message += reminders.GetReminderByPrefix(categoryShalat)
	message += " \n\n\n"
	return message
}

// ResolveUsername menentukan nama bot yang ditampilkan di webhook.
func ResolveUsername(target models.WebhookTarget, label string, categoryShalat string, testingMode bool) string {
	if testingMode {
		return config.Test
	}
	if target.Username != "" {
		return target.Username
	}
	if target.UsernameBeforeAdzan != "" && categoryShalat == models.SebelumAdzan {
		return target.UsernameBeforeAdzan
	}
	return "Reminder Adzan " + label
}

// SendWebhook mengirimkan pesan ke satu WebhookTarget
func SendWebhook(label string, embeds []models.DiscordEmbed, target models.WebhookTarget, categoryShalat string, testingMode bool) WebhookResult {
	result := WebhookResult{Name: target.Name}

	message := BuildWebhookContent(target, label, categoryShalat)

	payload := map[string]any{
		"content": message,
		"embeds":  embeds,
		"tts":     target.EnableTextToSpeech(),
	}

	if testingMode {
		payload["avatar_url"] = models.TestAvaatar
	} else {
		payload["avatar_url"] = models.RandomAvatarURL[rand.Intn(len(models.RandomAvatarURL))]
	}
	payload["username"] = ResolveUsername(target, label, categoryShalat, testingMode)

	body, err := json.Marshal(payload)
	if err != nil {
		result.Error = fmt.Errorf("proses konversi JSON (marshalling) gagal: %w", err)
		return result
	}

	resp, err := http.Post(target.URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		result.Error = fmt.Errorf("permintaan HTTP post gagal: %w", err)
		return result
	}
	defer resp.Body.Close()

	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		result.Error = fmt.Errorf("gagal membaca respons HTTP: %w", readErr)
		return result
	}

	if resp.StatusCode >= 300 {
		result.Error = fmt.Errorf("kode status %d", resp.StatusCode)
		fmt.Printf("Gagal mengirim webhook ke %s: %s\n", target.Name, string(bodyBytes))
	}
	return result
}
