package scheduler

import (
	"fmt"
	"time"

	"github.com/adamnasrudin03/pray-notify/config"
	"github.com/adamnasrudin03/pray-notify/models"
	"github.com/adamnasrudin03/pray-notify/notification"
	goprayer "github.com/hablullah/go-prayer"
)

// ValidPrayerLabels daftar label shalat yang didukung mode window.
var ValidPrayerLabels = []string{
	models.Subuh, models.Dhuha, models.Dzuhur, models.Ashar, models.Maghrib, models.Isya,
}

// IsValidPrayerLabel memeriksa apakah label shalat dikenali.
func IsValidPrayerLabel(label string) bool {
	for _, p := range ValidPrayerLabels {
		if label == p {
			return true
		}
	}
	return false
}

type targetPlan struct {
	target     models.WebhookTarget
	loc        *time.Location
	shalatTime time.Time
}

// RunPrayerWindow menjalankan notifikasi satu window shalat: sleep sampai tiap offset, lalu kirim.
// Dipicu 1× per shalat dari cron-job.org → GitHub repository_dispatch.
func RunPrayerWindow(prayer string, testingMode bool) error {
	if !IsValidPrayerLabel(prayer) {
		return fmt.Errorf("prayer tidak valid %q (gunakan: %v)", prayer, ValidPrayerLabels)
	}

	offsets := offsetsForPrayer(prayer)
	plans, err := collectTargetPlans(prayer, testingMode)
	if err != nil {
		return err
	}
	if len(plans) == 0 {
		fmt.Println("⏩ DILEWATI: tidak ada target aktif untuk window ini")
		return nil
	}

	fmt.Printf("\n=== Window shalat: %s (%d target) ===\n", prayer, len(plans))
	for _, p := range plans {
		fmt.Printf("  • %s (%s) adzan %s\n", p.target.Name, p.target.City, p.shalatTime.Format("15:04"))
	}

	for _, offset := range offsets {
		earliest := earliestNotificationTime(plans, offset)
		if !testingMode && time.Now().After(earliest.Add(notification.MaxDeliveryDelay())) {
			fmt.Printf("⏩ DILEWATI offset %v: sudah lewat window pengiriman\n", offset)
			continue
		}

		sleepUntil(earliest, testingMode)

		for _, plan := range plans {
			notificationTime := plan.shalatTime.Add(offset)
			if !testingMode && time.Now().After(notificationTime.Add(notification.MaxDeliveryDelay())) {
				continue
			}
			notification.CheckAndSend(notificationTime, prayer, plan.shalatTime, offset, plan.target, testingMode, true)
		}
	}

	fmt.Printf("=== Selesai window %s ===\n", prayer)
	return nil
}

func offsetsForPrayer(prayer string) []time.Duration {
	if prayer == models.Dhuha {
		return []time.Duration{0}
	}
	return MainOffsets
}

func collectTargetPlans(prayer string, testingMode bool) ([]targetPlan, error) {
	var plans []targetPlan

	for _, target := range config.WebhookTargets {
		if target.IsDisable {
			fmt.Printf("⏩ DILEWATI %s (IsDisable)\n", target.Name)
			continue
		}

		loc, err := time.LoadLocation(target.Tz)
		if err != nil {
			fmt.Printf("Error timezone %s untuk %s: %v\n", target.Tz, target.Name, err)
			continue
		}

		now := time.Now().In(loc)
		if !testingMode && !IsWithinOperatingWindow(now, loc) {
			fmt.Printf("⏩ DILEWATI %s: di luar window operasi\n", target.Name)
			continue
		}
		if target.ShouldSkipWeekend(now.Weekday()) {
			fmt.Printf("⏩ DILEWATI %s: akhir pekan\n", target.Name)
			continue
		}

		shalatTime, err := prayerTimeForTarget(target, prayer, now, loc)
		if err != nil {
			fmt.Printf("Error jadwal %s untuk %s: %v\n", prayer, target.Name, err)
			continue
		}

		plans = append(plans, targetPlan{target: target, loc: loc, shalatTime: shalatTime})
	}

	return plans, nil
}

func prayerTimeForTarget(target models.WebhookTarget, label string, now time.Time, loc *time.Location) (time.Time, error) {
	cfg := goprayer.Config{
		Latitude:           target.Lat,
		Longitude:          target.Lon,
		Timezone:           loc,
		TwilightConvention: goprayer.Kemenag(),
		AsrConvention:      goprayer.Shafii,
		PreciseToSeconds:   true,
	}

	schedules, err := goprayer.Calculate(cfg, now.Year())
	if err != nil {
		return time.Time{}, err
	}

	today := now.Format("2006-01-02")
	for _, s := range schedules {
		if s.Date != today {
			continue
		}
		switch label {
		case models.Subuh:
			return s.Fajr.In(loc), nil
		case models.Dhuha:
			return s.Sunrise.In(loc).Add(18 * time.Minute), nil
		case models.Dzuhur:
			return s.Zuhr.In(loc), nil
		case models.Ashar:
			return s.Asr.In(loc), nil
		case models.Maghrib:
			return s.Maghrib.In(loc), nil
		case models.Isya:
			return s.Isha.In(loc), nil
		}
	}
	return time.Time{}, fmt.Errorf("jadwal hari ini tidak ditemukan")
}

func earliestNotificationTime(plans []targetPlan, offset time.Duration) time.Time {
	var earliest time.Time
	for _, plan := range plans {
		t := plan.shalatTime.Add(offset)
		if earliest.IsZero() || t.Before(earliest) {
			earliest = t
		}
	}
	return earliest
}

func sleepUntil(until time.Time, testingMode bool) {
	if testingMode {
		return
	}
	if d := time.Until(until); d > 0 {
		fmt.Printf("⏳ Menunggu sampai %s (%s)\n", until.Format("15:04:05"), d.Truncate(time.Second))
		time.Sleep(d)
	}
}
