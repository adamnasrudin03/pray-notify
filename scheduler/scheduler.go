package scheduler

import (
	"fmt"
	"time"

	"github.com/adamnasrudin03/pray-notify/config"
	"github.com/adamnasrudin03/pray-notify/models"
	"github.com/adamnasrudin03/pray-notify/notification"
	"github.com/hablullah/go-prayer"
)

var offsetsMain = MainOffsets

type prayerSlot struct {
	label string
	t     time.Time
}

// Run menjalankan scheduler utama untuk notifikasi shalat
func Run(testingMode bool) {
	if !ShouldRunToday(testingMode) {
		fmt.Println("⏩ DILEWATI: di luar window operasi cron (03:30–19:59 WIB)")
		return
	}

	if !testingMode {
		if err := config.ValidateTargets(config.WebhookTargets); err != nil {
			fmt.Printf("❌ Konfigurasi webhook: %v\n", err)
			return
		}
	}

	for _, target := range config.WebhookTargets {
		if target.IsDisable {
			fmt.Printf("\n--- Target Webhook: %s (%s) DILEWATI (IsDisable: true) ---\n", target.Name, target.City)
			continue
		}

		loc, err := time.LoadLocation(target.Tz)
		if err != nil {
			fmt.Printf("Error memuat zona waktu %s untuk %s: %v. Melewati target.\n", target.Tz, target.Name, err)
			continue
		}

		now := time.Now().In(loc)
		if !testingMode && !IsWithinOperatingWindow(now, loc) {
			fmt.Printf("⏩ DILEWATI %s: di luar window operasi cron (03:30–19:59 %s)\n", target.Name, target.Tz)
			continue
		}

		if target.ShouldSkipWeekend(now.Weekday()) {
			fmt.Printf("webhook DILEWATI: %s tidak aktif di akhir pekan (Sabtu/Minggu)\n", target.Name)
			continue
		}

		fmt.Printf("\n--- Memproses Target Webhook: %s (%s) ---\n", target.Name, target.City)

		todayFormatted := now.Format("2006-01-02")

		cfg := prayer.Config{
			Latitude:           target.Lat,
			Longitude:          target.Lon,
			Timezone:           loc,
			TwilightConvention: prayer.Kemenag(),
			AsrConvention:      prayer.Shafii,
			PreciseToSeconds:   true,
		}

		schedules, err := prayer.Calculate(cfg, now.Year())
		if err != nil {
			fmt.Printf("Error menghitung jadwal shalat untuk %s: %v. Melewati target.\n", target.Name, err)
			continue
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
			fmt.Printf("Error: Jadwal hari ini tidak ditemukan untuk %s.\n", target.Name)
			continue
		}

		fajr := schedule.Fajr.In(loc)
		sunrise := schedule.Sunrise.In(loc)
		dhuhr := schedule.Zuhr.In(loc)
		asr := schedule.Asr.In(loc)
		maghrib := schedule.Maghrib.In(loc)
		isha := schedule.Isha.In(loc)

		dhuhaStart := sunrise.Add(18 * time.Minute)

		prayerTimes := []prayerSlot{
			{models.Subuh, fajr},
			{models.Dhuha, dhuhaStart},
			{models.Dzuhur, dhuhr},
			{models.Ashar, asr},
			{models.Maghrib, maghrib},
			{models.Isya, isha},
		}

		if !testingMode && !hasPendingNotifications(prayerTimes, now, testingMode) {
			fmt.Printf("⏩ DILEWATI %s: tidak ada notifikasi dalam window aktif\n", target.Name)
			continue
		}

		for _, p := range prayerTimes {
			if p.label == models.Dhuha {
				notification.CheckAndSend(p.t, p.label, p.t, 0, target, testingMode, false)
				continue
			}

			activeOffsets := notification.FilterActiveOffsets(p.t, now, offsetsMain, testingMode)
			for _, offset := range activeOffsets {
				notificationTime := p.t.Add(offset)
				notification.CheckAndSend(notificationTime, p.label, p.t, offset, target, testingMode, false)
			}
		}
	}
}

func hasPendingNotifications(prayerTimes []prayerSlot, now time.Time, testingMode bool) bool {
	for _, p := range prayerTimes {
		if p.label == models.Dhuha {
			if notification.IsInNotificationWindow(p.t, now, testingMode) {
				return true
			}
			continue
		}
		if len(notification.FilterActiveOffsets(p.t, now, offsetsMain, testingMode)) > 0 {
			return true
		}
	}
	return false
}
