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

	// Definisikan offsets yang akan digunakan dalam mode testing, mirip dengan mode normal
	offsetsMain := []time.Duration{}
	for m := -5; m <= 30; m += 5 {
		offsetsMain = append(offsetsMain, time.Duration(m)*time.Minute)
	}

	offsetsCaseSpecific := []time.Duration{}
	for m := -5; m <= 5; m += 1 {
		offsetsCaseSpecific = append(offsetsCaseSpecific, time.Duration(m)*time.Minute)
	}

	fmt.Println("==================================================================================")
	fmt.Printf("!!! MODE UJI COBA AKTIF: Semua notifikasi akan dialihkan ke Webhook Testing !!!\n")
	fmt.Println("==================================================================================")

	// Hanya loop pada target yang sama dengan MyChannel untuk mendapatkan jadwal shalat
	for _, target := range config.WebhookTargets {
		if target.Name != config.MyChannel {
			continue
		}
		fmt.Printf("\n--- Memproses Target Webhook: %s (%s) ---\n", target.Name, target.City)

		// 2. Load Timezone spesifik target
		loc, err := time.LoadLocation(target.Tz)
		if err != nil {
			fmt.Printf("Error memuat zona waktu %s untuk %s: %v. Melewati target.\n", target.Tz, target.Name, err)
			continue
		}

		today := time.Now().In(loc)
		todayFormatted := today.Format("2006-01-02")

		// 3. Hitung Jadwal Shalat untuk Target ini
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

		// 4. Ambil waktu shalat dalam zona waktu target
		fajr := schedule.Fajr.In(loc)
		sunrise := schedule.Sunrise.In(loc)
		dhuhr := schedule.Zuhr.In(loc)
		asr := schedule.Asr.In(loc)
		maghrib := schedule.Maghrib.In(loc)
		isha := schedule.Isha.In(loc)

		// 🔥 HITUNG AWAL WAKTU DHUHA (18 menit setelah Sunrise)
		dhuhaStart := sunrise.Add(18 * time.Minute)

		prayerTimes := []struct {
			label string
			t     time.Time
		}{
			{models.Subuh, fajr},
			{models.Dhuha, dhuhaStart}, // Tambahkan Dhuha di sini untuk testing
			{models.Dzuhur, dhuhr},
			{models.Ashar, asr},
			{models.Maghrib, maghrib},
			{models.Isya, isha},
		}

		// 5. Loop di setiap waktu shalat dan offset
		for _, p := range prayerTimes {
			currentOffsets := offsetsMain
			// Gunakan offsetsCaseSpecific untuk Dhuha
			if p.label == models.Dhuha {
				currentOffsets = offsetsCaseSpecific
			}

			for _, offset := range currentOffsets {
				notificationTime := p.t.Add(offset)

				if offset.Minutes() != 0 {
					continue
				}

				// jika testingMode=true dan target adalah target yang sedang diproses.
				notification.CheckAndSend(notificationTime, p.label, p.t, offset, target, testingMode)
				time.Sleep(1 * time.Minute)
			}
		}
		break
	}
}
