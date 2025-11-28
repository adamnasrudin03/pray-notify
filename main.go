package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hablullah/go-prayer"
)

// Struktur baru untuk menyimpan URL, Nama/Label Webhook, dan Mentions CC
type WebhookTarget struct {
	Name     string
	URL      string
	Mentions string // Field baru untuk CC (@user @user)
}

var (
	// Koordinat Kebon Jeruk, Jakarta Barat
	lat      = -6.183
	lon      = 106.784
	timezone = "Asia/Jakarta"

	// Daftar Webhook Target yang Diperbarui dengan Mentions
	webhookTargets = []WebhookTarget{
		{
			Name: "Divisi IT",
			URL:  "https://discord.com/api/webhooks/1443879495324008488/z-jzPBXv6vBvE45-QFrV-wwcIp8Qz9GtMw0bSQeKh1ltjHXZjcUPmjN3pMpS7Jg4yZJ0",
			// Mentions: "@abi_lp @tedja_putra_lp @ade_muhammad_lp @bian_kusumah_lp", // CC hanya untuk Divisi IT
		},
		{
			Name:     "Apran Channel",
			URL:      "https://discordapp.com/api/webhooks/1443512213842296904/EHX00Lvi8Hmqxzc1WsgJK9pyDsJzAmXWhexbqSsJnez8KSzaAHaBy_8UQIAgtRQyQEix",
			Mentions: "", // Tidak ada CC untuk channel ini
		},
		// Tambahkan webhook lain di sini
	}
)

// WebhookResult menyimpan hasil pengiriman ke satu webhook
type WebhookResult struct {
	Name  string
	Error error
}

// sendWebhook mengirimkan pesan ke semua WebhookTarget
// Mengembalikan slice dari hasil pengiriman (berhasil atau gagal)
func sendWebhook(title, description string, embedColor int) []WebhookResult {

	results := []WebhookResult{}
	wibLoc := time.FixedZone("WIB", 7*3600)

	// Siapkan data Embed (Payload ini akan sama untuk semua target)
	embed := map[string]interface{}{
		"color":       embedColor,
		"title":       title,
		"description": description,
		"timestamp":   time.Now().Format(time.RFC3339),
		"footer": map[string]string{
			"text": fmt.Sprintf("Lokasi: Kebon Jeruk, Jakarta Barat | Dihitung: %s WIB", time.Now().In(wibLoc).Format("02/01/2006 15:04")),
		},
	}

	// Loop melalui semua target webhook
	for _, target := range webhookTargets {

		// 🔥 MODIFIKASI: Buat payload baru di dalam loop untuk menyesuaikan content
		payload := map[string]interface{}{
			// Masukkan Mentions di sini. Jika kosong, content akan kosong.
			"content": target.Mentions,
			"embeds":  []map[string]interface{}{embed},
		}

		body, err := json.Marshal(payload)

		result := WebhookResult{Name: target.Name}

		if err != nil {
			result.Error = fmt.Errorf("marshalling failed: %w", err)
			results = append(results, result)
			continue // Lanjut ke target berikutnya jika marshalling gagal
		}

		resp, err := http.Post(target.URL, "application/json", bytes.NewBuffer(body))

		if err != nil {
			result.Error = fmt.Errorf("HTTP post failed: %w", err)
		} else {
			// Perhatian: defer resp.Body.Close() tidak bisa digunakan di dalam loop,
			// atau harus ditangani secara eksplisit. Untuk kesederhanaan, kita akan tetap
			// memanggilnya langsung karena Go akan membersihkan sumber daya.
			resp.Body.Close()
			if resp.StatusCode >= 300 {
				result.Error = fmt.Errorf("status code %d", resp.StatusCode)
			}
		}
		results = append(results, result)
	}

	return results
}

// checkAndSend notifikasi (Tidak ada perubahan pada logika ini)
func checkAndSend(notificationTime time.Time, label string, shalatTime time.Time, offsetDuration time.Duration) {

	loc, _ := time.LoadLocation(timezone)
	now := time.Now().In(loc) // now sekarang adalah waktu WIB

	fiveMinutesAgo := now.Add(-5 * time.Minute)

	isBefore := offsetDuration < 0

	absDuration := offsetDuration
	if absDuration < 0 {
		absDuration = -absDuration
	}
	absDuration = absDuration.Truncate(time.Minute)

	if notificationTime.After(fiveMinutesAgo) && (notificationTime.Before(now) || notificationTime.Equal(now)) {

		var title string
		var description string
		var embedColor int // Warna dalam format desimal

		if isBefore {
			durationUntilShalat := shalatTime.Sub(now).Truncate(time.Second)
			embedColor = 16763904 // Warna Kuning

			if durationUntilShalat > 0 {
				title = fmt.Sprintf("🕌 Adzan %s Akan Tiba!", label)
				description = fmt.Sprintf("Waktu Shalat akan berkumandang dalam **%s** (Pukul **%s**). Segera siapkan diri.",
					durationUntilShalat.Truncate(time.Minute),
					shalatTime.Format("15:04"),
				)
			} else {
				title = fmt.Sprintf("🚨 ADZAN %s TELAH TIBA!", label)
				description = fmt.Sprintf("Adzan telah berkumandang (Pukul **%s**). Mohon maaf, pengingat **%v** sebelum Adzan tertunda. Segera Tunaikan Shalat.",
					shalatTime.Format("15:04"),
					absDuration,
				)
			}

		} else {
			embedColor = 5814783 // Warna Hijau

			title = fmt.Sprintf("✅ WAKTU SALAT %s SEDANG BERLANGSUNG", label)
			description = fmt.Sprintf("Shalat %s sudah dimulai (Pukul **%s**). Ini adalah notifikasi **%v** setelah Adzan. Saatnya Tunaikan Shalat.",
				label,
				shalatTime.Format("15:04"),
				absDuration,
			)
		}

		// Panggil sendWebhook dan terima hasilnya
		results := sendWebhook(title, description, embedColor)

		// Log hasil pengiriman untuk setiap webhook
		for _, result := range results {
			logMessage := fmt.Sprintf("%s at %s. Type: %s. Webhook: %s",
				label,
				notificationTime.Format("15:04"),
				func() string {
					if isBefore {
						return "SEBELUM"
					} else {
						return "SETELAH"
					}
				}(),
				result.Name, // Nama Webhook
			)

			if result.Error != nil {
				fmt.Printf("❌ Error sending EMBED notification for %s. Error: %v\n", logMessage, result.Error)
			} else {
				fmt.Printf("✅ Sent EMBED notification for %s.\n", logMessage)
			}
		}

	} else {
		// Logika SKIPPED
		fmt.Printf(
			"⏩ SKIPPED %s (%v %s Adzan). Waktu Adzan: %s. Waktu Notif Seharusnya: %s. Waktu Saat Ini: %s WIB.\n",
			label,
			absDuration,
			func() string {
				if isBefore {
					return "SEBELUM"
				} else {
					return "SETELAH"
				}
			}(),
			shalatTime.Format("15:04"),
			notificationTime.Format("15:04:05"),
			now.Format("15:04:05"),
		)
	}
}

// ---------------------------------------------------------------------------------------

func main() {
	// ... (Sisa fungsi main tetap sama) ...
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		fmt.Printf("Error loading timezone: %v\n", err)
		os.Exit(1)
	}

	today := time.Now().In(loc)
	todayFormatted := today.Format("2006-01-02")

	cfg := prayer.Config{
		Latitude:  lat,
		Longitude: lon,
		Timezone:  loc,
	}

	schedules, err := prayer.Calculate(cfg, today.Year())
	if err != nil {
		fmt.Printf("Error calculating prayer schedule: %v\n", err)
		os.Exit(1)
	}

	var schedule prayer.Schedule
	found := false

	for _, s := range schedules {
		if s.Date == todayFormatted {
			schedule = s
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Error: Today's schedule not found.")
		os.Exit(1)
	}

	fajr := schedule.Fajr.In(loc)
	dhuhr := schedule.Zuhr.In(loc)
	asr := schedule.Asr.In(loc)
	maghrib := schedule.Maghrib.In(loc)
	isha := schedule.Isha.In(loc)

	prayerTimes := []struct {
		label string
		t     time.Time
	}{
		{"Subuh", fajr},
		{"Dzuhur", dhuhr},
		{"Ashar", asr},
		{"Maghrib", maghrib},
		{"Isya", isha},
	}

	offsets := []time.Duration{
		-30 * time.Minute,
		-20 * time.Minute,
		-15 * time.Minute,
		-5 * time.Minute,
		5 * time.Minute,
		10 * time.Minute,
		15 * time.Minute,
		20 * time.Minute,
		30 * time.Minute,
	}

	for _, prayer := range prayerTimes {
		for _, offset := range offsets {
			notificationTime := prayer.t.Add(offset)
			checkAndSend(notificationTime, prayer.label, prayer.t, offset)
		}
	}
}
