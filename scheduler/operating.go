package scheduler

import (
	"time"

	"github.com/adamnasrudin03/pray-notify/config"
)

// Window operasi selaras dengan cron-job.org (timezone target).
const (
	OperatingStartHour   = 3
	OperatingStartMinute = 30
	OperatingEndHour     = 19
	OperatingEndMinute   = 59
)

// IsWithinOperatingWindow memeriksa apakah waktu sekarang dalam window cron 03:30–19:59.
func IsWithinOperatingWindow(now time.Time, loc *time.Location) bool {
	t := now.In(loc)
	minutes := t.Hour()*60 + t.Minute()
	start := OperatingStartHour*60 + OperatingStartMinute
	end := OperatingEndHour*60 + OperatingEndMinute
	return minutes >= start && minutes <= end
}

// ShouldRunToday memeriksa apakah ada target yang perlu diproses (window operasi + bukan disable).
func ShouldRunToday(testingMode bool) bool {
	if testingMode {
		return true
	}
	for _, target := range config.WebhookTargets {
		if target.IsDisable {
			continue
		}
		loc, err := time.LoadLocation(target.Tz)
		if err != nil {
			continue
		}
		if IsWithinOperatingWindow(time.Now(), loc) {
			return true
		}
	}
	return false
}
