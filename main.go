package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hablullah/go-prayer"
)

var (
	// koordinat Kebon Jeruk, Jakarta Barat
	lat      = -6.183
	lon      = 106.784
	timezone = "Asia/Jakarta"

	webhookURL = "https://discordapp.com/api/webhooks/1443512213842296904/EHX00Lvi8Hmqxzc1WsgJK9pyDsJzAmXWhexbqSsJnez8KSzaAHaBy_8UQIAgtRQyQEix"
)

func sendWebhook(message string) error {
	payload := map[string]string{
		"text": message,
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned %d", resp.StatusCode)
	}
	return nil
}

func scheduleNotification(t time.Time, label string) {
	now := time.Now()
	diff := t.Sub(now)
	if diff <= 0 {
		fmt.Printf("Skipped %s — time already passed: %v\n", label, t)
		return
	}

	go func() {
		timer := time.NewTimer(diff)
		<-timer.C
		msg := fmt.Sprintf("Adzan %s akan segera (%s)", label, t.Format("15:04"))
		if err := sendWebhook(msg); err != nil {
			fmt.Printf("Error sending webhook: %v\n", err)
		} else {
			fmt.Printf("Sent notification for %s\n", label)
		}
	}()
}

func main() {
	loc, _ := time.LoadLocation(timezone)
	year := time.Now().Year()

	cfg := prayer.Config{
		Latitude:  lat,
		Longitude: lon,
		Timezone:  loc,
	}

	schedules, err := prayer.Calculate(cfg, year)
	if err != nil {
		panic(err)
	}

	// Cari jadwal hari ini
	today := time.Now().Format("2006-01-02")

	for _, s := range schedules {
		if s.Date != today {
			continue
		}

		// Waktu salat lokal
		fajr := s.Fajr.In(loc)
		dhuhr := s.Zuhr.In(loc)
		asr := s.Asr.In(loc)
		maghrib := s.Maghrib.In(loc)
		isha := s.Isha.In(loc)

		// Schedule 30, 15, 5 menit sebelum tiap salat
		for _, item := range []struct {
			label string
			t     time.Time
		}{
			{"Subuh", fajr},
			{"Dzuhur", dhuhr},
			{"Ashar", asr},
			{"Maghrib", maghrib},
			{"Isya", isha},
		} {
			scheduleNotification(item.t.Add(-30*time.Minute), item.label)
			scheduleNotification(item.t.Add(-15*time.Minute), item.label)
			scheduleNotification(item.t.Add(-5*time.Minute), item.label)
		}
	}

	select {} // keep program running sampai semua timer selesai
}
