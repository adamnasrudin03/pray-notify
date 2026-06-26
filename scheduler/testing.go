package scheduler

import (
	"fmt"
	"time"

	"github.com/adamnasrudin03/pray-notify/config"
	"github.com/adamnasrudin03/pray-notify/models"
	"github.com/adamnasrudin03/pray-notify/notification"
	"github.com/hablullah/go-prayer"
)

// RunTesting menjalankan mode testing untuk notifikasi
func RunTesting(testingMode bool) {
	if !testingMode {
		return
	}

	offsetsMainTest := []time.Duration{}
	for m := -5; m <= 30; m += 5 {
		offsetsMainTest = append(offsetsMainTest, time.Duration(m)*time.Minute)
	}

	fmt.Println("==================================================================================")
	fmt.Printf("!!! MODE UJI COBA AKTIF: Semua notifikasi akan dialihkan ke Webhook Testing !!!\n")
	fmt.Println("==================================================================================")

	for _, target := range config.WebhookTargets {
		if target.Name != config.MyChannel {
			continue
		}
		fmt.Printf("\n--- Memproses Target Webhook: %s (%s) ---\n", target.Name, target.City)

		loc, err := time.LoadLocation(target.Tz)
		if err != nil {
			fmt.Printf("Error memuat zona waktu %s untuk %s: %v. Melewati target.\n", target.Tz, target.Name, err)
			continue
		}

		today := time.Now().In(loc)
		todayFormatted := today.Format("2006-01-02")

		cfg := prayer.Config{
			Latitude:           target.Lat,
			Longitude:          target.Lon,
			Timezone:           loc,
			TwilightConvention: prayer.Kemenag(),
			AsrConvention:      prayer.Shafii,
		}

		schedules, err := prayer.Calculate(cfg, today.Year())
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

		prayerTimes := []struct {
			label string
			t     time.Time
		}{
			{models.Subuh, fajr},
			{models.Dhuha, dhuhaStart},
			{models.Dzuhur, dhuhr},
			{models.Ashar, asr},
			{models.Maghrib, maghrib},
			{models.Isya, isha},
		}

		for _, p := range prayerTimes {
			if p.label == models.Dhuha {
				notification.CheckAndSend(p.t, p.label, p.t, 0, target, testingMode)
				time.Sleep(1 * time.Minute)
				continue
			}

			for _, offset := range offsetsMainTest {
				if offset.Minutes() != 0 {
					continue
				}
				notificationTime := p.t.Add(offset)
				notification.CheckAndSend(notificationTime, p.label, p.t, offset, target, testingMode)
				time.Sleep(1 * time.Minute)
			}
		}
		break
	}
}
