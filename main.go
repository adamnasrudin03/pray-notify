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

	// Pastikan ini adalah URL webhook Discord yang valid
	webhookURL = "https://discordapp.com/api/webhooks/1443512213842296904/EHX00Lvi8Hmqxzc1WsgJK9pyDsJzAmXWhexbqSsJnez8KSzaAHaBy_8UQIAgtRQyQEix"
)

// sendWebhook mengirimkan pesan ke Discord Webhook
func sendWebhook(message string) error {
	// Discord menggunakan kunci "content", bukan "text"
	payload := map[string]string{
		"content": message,
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

// checkAndSend notifikasi hanya jika waktunya berada dalam jendela eksekusi 5 menit terakhir.
// offsetDuration: Negatif untuk sebelum Adzan, Positif untuk setelah Adzan.
func checkAndSend(notificationTime time.Time, label string, salatTime time.Time, offsetDuration time.Duration) {
	now := time.Now()
	// Asumsi GA dijalankan setiap 5 menit
	fiveMinutesAgo := now.Add(-5 * time.Minute)

	// Tentukan apakah notifikasi ini sebelum atau sesudah Adzan
	isBefore := offsetDuration < 0

	// Gunakan time.Abs() untuk mendapatkan durasi absolut (misal: 5m, 15m)
	absDuration := offsetDuration
	if absDuration < 0 {
		absDuration = -absDuration
	}
	absDuration = absDuration.Truncate(time.Minute) // Pastikan hanya menit

	// Cek apakah waktu notifikasi berada dalam jendela 5 menit terakhir
	if notificationTime.After(fiveMinutesAgo) && (notificationTime.Before(now) || notificationTime.Equal(now)) {

		var msg string

		if isBefore {
			// --- Notifikasi Sebelum Adzan (30, 15, 5 menit) ---
			durationUntilSalat := salatTime.Sub(now).Truncate(time.Second)

			if durationUntilSalat > 0 {
				// Pesan saat masih ada waktu tersisa
				msg = fmt.Sprintf("🔔 **Adzan %s** akan tiba dalam %s (Pukul %s).",
					label,
					durationUntilSalat.Truncate(time.Minute), // Sisa waktu
					salatTime.Format("15:04"),
				)
			} else {
				// Pesan darurat jika GA terlambat dan waktu Adzan sudah lewat
				msg = fmt.Sprintf("🚨 **Adzan %s** telah tiba (Pukul %s). Cepat! (Pengingat %v sebelum Adzan tertunda).",
					label,
					salatTime.Format("15:04"),
					absDuration, // Gunakan absDuration yang sudah di-Truncate
				)
			}

		} else {
			// --- Notifikasi Setelah Adzan (5, 10 menit) ---
			msg = fmt.Sprintf("📢 **Waktu salat %s** sudah dimulai (Pukul %s). Ini adalah notifikasi %v setelah Adzan. Saatnya Tunaikan Salat.",
				label,
				salatTime.Format("15:04"),
				absDuration, // Durasi 5 atau 10 menit
			)
		}

		if err := sendWebhook(msg); err != nil {
			fmt.Printf("❌ Error sending webhook for %s: %v\n", label, err)
		} else {
			fmt.Printf("✅ Sent notification for %s at %s. Type: %s.\n",
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
		// --- Custom Debug Log (Skipped) ---
		// Pesan yang lebih informatif (sesuai permintaan sebelumnya)
		fmt.Printf(
			"⏩ SKIPPED %s (%v %s Adzan). Waktu Adzan: %s. Waktu Notif Seharusnya: %s. Waktu Saat Ini: %s.\n",
			label,
			absDuration, // Durasi Notifikasi (5m, 10m, dst)
			func() string {
				if isBefore {
					return "SEBELUM"
				} else {
					return "SETELAH"
				}
			}(), // Label Sebelum/Setelah
			salatTime.Format("15:04"),           // Waktu Salat
			notificationTime.Format("15:04:05"), // Waktu Notifikasi yang dijadwalkan
			now.Format("15:04:05"),              // Waktu Eksekusi GA saat ini
		)
	}
}

func main() {
	// sendWebhook("🚀 Prayer Notifier started.")
	// Memuat lokasi zona waktu
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		fmt.Printf("Error loading timezone: %v\n", err)
		os.Exit(1)
	}

	// Tanggal dan waktu saat ini di zona lokal
	today := time.Now().In(loc)
	todayFormatted := today.Format("2006-01-02")

	cfg := prayer.Config{
		Latitude:  lat,
		Longitude: lon,
		Timezone:  loc,
	}

	// FIX: Gunakan Calculate dan Filter (karena CalculateByDay tidak tersedia)
	schedules, err := prayer.Calculate(cfg, today.Year())
	if err != nil {
		fmt.Printf("Error calculating prayer schedule: %v\n", err)
		os.Exit(1)
	}

	var schedule prayer.Schedule
	found := false

	// Cari jadwal hari ini
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
	// END FIX

	// Waktu salat lokal (menggunakan jadwal yang telah difilter)
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

	// Jangka waktu notifikasi. Negatif (-) untuk SEBELUM, Positif (+) untuk SETELAH.
	offsets := []time.Duration{
		-30 * time.Minute, // 30 menit sebelum
		-15 * time.Minute, // 15 menit sebelum
		-5 * time.Minute,  // 5 menit sebelum
		5 * time.Minute,   // 5 menit setelah (BARU)
		10 * time.Minute,  // 10 menit setelah (BARU)
	}

	// Jadwalkan dan kirim notifikasi jika waktunya tepat
	for _, prayer := range prayerTimes {
		for _, offset := range offsets {
			// notificationTime = WaktuSalat + offset (bisa negatif atau positif)
			notificationTime := prayer.t.Add(offset)

			// checkAndSend menerima offset Duration (termasuk tanda +/-)
			checkAndSend(notificationTime, prayer.label, prayer.t, offset)
		}
	}
}
