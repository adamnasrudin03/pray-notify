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

var (
	// Koordinat Kebon Jeruk, Jakarta Barat
	lat      = -6.183
	lon      = 106.784
	timezone = "Asia/Jakarta"

	webhookURL = "https://discordapp.com/api/webhooks/1443512213842296904/EHX00Lvi8Hmqxzc1WsgJK9pyDsJzAmXWhexbqSsJnez8KSzaAHaBy_8UQIAgtRQyQEix"
)

// sendWebhook mengirimkan pesan ke Discord Webhook menggunakan Embeds
func sendWebhook(title, description string, embedColor int) error {

	// Konfigurasi zona WIB untuk Footer
	wibLoc := time.FixedZone("WIB", 7*3600)

	// Siapkan data Embed
	embed := map[string]interface{}{
		"color":       embedColor,
		"title":       title,
		"description": description,
		"timestamp":   time.Now().Format(time.RFC3339),

		// 🚨 INFO LOKASI DITAMBAHKAN DI FOOTER
		"footer": map[string]string{
			"text": fmt.Sprintf("Lokasi: Kebon Jeruk, Jakarta Barat | Dihitung: %s WIB", time.Now().In(wibLoc).Format("02/01/2006 15:04")),
		},
	}

	// Payload utama Discord Webhook
	payload := map[string]interface{}{
		"content": "", // Teks utama
		"embeds":  []map[string]interface{}{embed},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling payload: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-success status: %d", resp.StatusCode)
	}
	return nil
}

// checkAndSend notifikasi
func checkAndSend(notificationTime time.Time, label string, salatTime time.Time, offsetDuration time.Duration) {

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
			durationUntilSalat := salatTime.Sub(now).Truncate(time.Second)
			embedColor = 16763904 // Warna Kuning

			if durationUntilSalat > 0 {
				title = fmt.Sprintf("🕌 Adzan %s Akan Tiba!", label)
				description = fmt.Sprintf("Waktu salat akan berkumandang dalam **%s** (Pukul **%s**). Segera siapkan diri.",
					durationUntilSalat.Truncate(time.Minute),
					salatTime.Format("15:04"),
				)
			} else {
				title = fmt.Sprintf("🚨 ADZAN %s TELAH TIBA!", label)
				description = fmt.Sprintf("@here Adzan telah berkumandang (Pukul **%s**). Mohon maaf, pengingat **%v** sebelum Adzan tertunda. Segera Tunaikan Salat.",
					salatTime.Format("15:04"),
					absDuration,
				)
			}

		} else {
			embedColor = 5814783 // Warna Hijau

			title = fmt.Sprintf("✅ WAKTU SALAT %s SEDANG BERLANGSUNG", label)
			description = fmt.Sprintf("Salat %s sudah dimulai (Pukul **%s**). Ini adalah notifikasi **%v** setelah Adzan. Saatnya Tunaikan Salat.",
				label,
				salatTime.Format("15:04"),
				absDuration,
			)
		}

		if err := sendWebhook(title, description, embedColor); err != nil {
			fmt.Printf("❌ Error sending webhook for %s: %v\n", label, err)
		} else {
			fmt.Printf("✅ Sent EMBED notification for %s at %s. Type: %s.\n",
				label,
				notificationTime.Format("15:04"),
				func() string {
					if isBefore {
						return "SEBELUM"
					} else {
						return "SETELAH"
					}
				}(),
			)
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
			salatTime.Format("15:04"),
			notificationTime.Format("15:04:05"),
			now.Format("15:04:05"),
		)
	}
}

// ---------------------------------------------------------------------------------------

func main() {
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
		-15 * time.Minute,
		-5 * time.Minute,
		5 * time.Minute,
		10 * time.Minute,
		15 * time.Minute,
	}

	for _, prayer := range prayerTimes {
		for _, offset := range offsets {
			notificationTime := prayer.t.Add(offset)
			checkAndSend(notificationTime, prayer.label, prayer.t, offset)
		}
	}
}
