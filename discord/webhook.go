package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/adamnasrudin03/pray-notify/config"
	"github.com/adamnasrudin03/pray-notify/models"
	"github.com/adamnasrudin03/pray-notify/reminders"
)

// WebhookResult menyimpan hasil pengiriman ke satu webhook
type WebhookResult struct {
	Name  string
	Error error
}

// SendWebhook mengirimkan pesan ke satu WebhookTarget
func SendWebhook(label string, embeds []models.DiscordEmbed, target models.WebhookTarget, categoryShalat string, testingMode bool) WebhookResult {
	result := WebhookResult{Name: target.Name}
	now := time.Now()

	// 🛑 VALIDASI 1: Cek Hari Libur/Weekend untuk Webhook "Divisi IT"
	if target.Name == config.DivisiIT && (now.Weekday() == time.Saturday || now.Weekday() == time.Sunday) {
		// Mengganti pesan error
		result.Error = fmt.Errorf("webhook DILEWATI: Divisi IT tidak aktif di akhir pekan (Sabtu/Minggu)")
		return result
	}
	isAdzanWork := map[string]bool{
		models.Dzuhur:  true,
		models.Ashar:   true,
		models.Maghrib: true,
	}

	message := "📢 @here \n"
	if !isAdzanWork[label] {
		message = "📢📢📢 \n"
	}
	message += reminders.GetReminderByPrefix(categoryShalat)
	message += " \n\n\n"

	mentions := target.Mentions
	if target.Name == config.DivisiIT && isAdzanWork[label] {
		randomMembers := models.GetRandomMemberDivIT(7)
		var m []string
		for _, member := range randomMembers {
			m = append(m, member.UserID)
		}
		mentions = strings.Join(m, " ")
	}

	if mentions != "" {
		message += "Cc: " + mentions
	}

	// Payload utama Discord Webhook (Field-field ini tetap 'username', 'content', dll. karena itu adalah kunci JSON Discord)
	payload := map[string]any{
		"content": message,
		"embeds":  embeds,
		"tts":     true,
	}

	payload["avatar_url"] = models.RandomAvatarURL[rand.Intn(len(models.RandomAvatarURL))]
	payload["username"] = "Reminder Adzan " + label

	if target.Name == config.DivisiIT {
		temp := models.GetRandomMemberDivIT(1)
		for _, v := range temp {
			payload["username"] = v.Name
			break
		}
	}
	if target.Name != config.DivisiIT && categoryShalat == models.SebelumAdzan {
		payload["username"] = config.MyUsername
	}

	if testingMode {
		payload["username"] = config.Test
		payload["avatar_url"] = models.TestAvaatar
	}

	body, err := json.Marshal(payload)
	if err != nil {
		result.Error = fmt.Errorf("proses konversi JSON (marshalling) gagal: %w", err)
		return result
	}

	resp, err := http.Post(target.URL, "application/json", bytes.NewBuffer(body))

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		result.Error = fmt.Errorf("permintaan HTTP post gagal: %w", err)
	} else {
		if resp.StatusCode >= 300 {
			result.Error = fmt.Errorf("kode status %d", resp.StatusCode)
			fmt.Printf("Gagal mengirim webhook ke %s: %+v\n", target.Name, bodyString)
		}
	}
	return result
}
