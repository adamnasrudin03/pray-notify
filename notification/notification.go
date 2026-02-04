package notification

import (
	"fmt"
	"strings"
	"time"

	"github.com/adamnasrudin03/pray-notify/discord"
	"github.com/adamnasrudin03/pray-notify/models"
)

// CheckAndSend memeriksa waktu dan mengirim notifikasi jika dalam window yang valid
func CheckAndSend(notificationTime time.Time, label string, shalatTime time.Time, offsetDuration time.Duration, target models.WebhookTarget, testingMode bool) {
	// Gunakan Timezone dari target Webhook
	loc, err := time.LoadLocation(target.Tz)
	if err != nil {
		fmt.Printf("Error memuat zona waktu target %s: %v\n", target.Tz, err)
		return
	}

	now := time.Now().In(loc)
	fiveMinutesAgo := now.Add(-5 * time.Minute)

	isBefore := offsetDuration < 0

	absDuration := offsetDuration
	if absDuration < 0 {
		absDuration = -absDuration
	}
	absDuration = absDuration.Truncate(time.Minute)

	isExactlyTime := absDuration.Minutes() < 1
	isAfter := !isBefore && !isExactlyTime

	// Kirim hanya jika dalam window waktu yang valid
	if testingMode || (notificationTime.After(fiveMinutesAgo) && (notificationTime.Before(now) || notificationTime.Equal(now))) {

		var mainTitle string
		var mainDescription string
		var embedColor int

		// 🔥 KHUSUS DHUHA
		if label == "Dhuha" {
			embedColor = 16763904 // Kuning
			mainTitle = "✨ Waktu Shalat Dhuha Telah Dimulai"
			mainDescription = fmt.Sprintf(
				"📢 @here Shalat Dhuha dapat ditunaikan mulai **pukul %s WIB**. Semoga membawa keberkahan dan kelancaran dalam aktivitas hari ini.",
				shalatTime.Format("15:04"),
			)

		} else if label == models.Dzuhur && now.Weekday() == time.Friday && (isBefore || isExactlyTime) {
			// 🔥 KHUSUS JUMATAN (Laki-laki) - Range -30 s.d 0 menit
			if isExactlyTime {
				embedColor = 5814783 // Hijau
				mainTitle = "🕌 Waktu Shalat Jumat Telah Tiba"
				mainDescription = fmt.Sprintf(
					"📢 @here Adzan Jumat telah berkumandang pada **pukul %s WIB**. Bagi rekan laki-laki, mari segera menuju masjid untuk menunaikan Shalat Jumat.",
					shalatTime.Format("15:04"),
				)
			} else {
				// Reminder Sebelum Jumatan
				embedColor = 16763904 // Kuning
				durationUntilShalat := shalatTime.Sub(now).Truncate(time.Second)
				minutes := int(durationUntilShalat.Minutes())

				if durationUntilShalat > 0 {
					mainTitle = "🕌 Persiapan Shalat Jumat"
					mainDescription = fmt.Sprintf(
						"📢 @here Waktu Shalat Jumat akan tiba dalam **%d menit** (pukul **%s WIB**). Bagi rekan laki-laki, mari persiapkan diri dan segera menuju masjid.",
						minutes,
						shalatTime.Format("15:04"),
					)
				} else {
					// Fallback jika terlambat
					formattedDurationBefore := fmt.Sprintf("%d menit", int(absDuration.Minutes()))
					mainTitle = "⚠️ Waktu Jumatan Telah Tiba"
					mainDescription = fmt.Sprintf(
						"📢 @here Waktu Jumatan telah masuk pada **pukul %s WIB**. Mohon maaf, pengingat **%s** sebelum waktu terkirim terlambat.",
						shalatTime.Format("15:04"),
						formattedDurationBefore,
					)
				}
			}
		} else if isBefore {
			durationUntilShalat := shalatTime.Sub(now).Truncate(time.Second)
			minutes := int(durationUntilShalat.Minutes())
			embedColor = 16763904 // Kuning

			if durationUntilShalat > 0 {
				mainTitle = fmt.Sprintf("🕌 Adzan %s Akan Segera Tiba", label)
				mainDescription = fmt.Sprintf(
					"📢 @here Adzan %s akan berkumandang dalam **%d menit** (pukul **%s WIB**). Mari bersiap untuk menunaikan shalat.",
					label,
					minutes,
					shalatTime.Format("15:04"),
				)
			} else {
				formattedDurationBefore := fmt.Sprintf("%d menit", int(absDuration.Minutes()))
				mainTitle = fmt.Sprintf("⚠️ Adzan %s Telah Berlangsung", label)
				mainDescription = fmt.Sprintf(
					"📢 @here Adzan %s telah berkumandang pada **pukul %s WIB**. Mohon maaf, pengingat **%s** sebelum adzan terkirim terlambat. Silakan segera menunaikan shalat.",
					label,
					shalatTime.Format("15:04"),
					formattedDurationBefore,
				)
			}

		} else {
			embedColor = 5814783 // Hijau

			if isExactlyTime {
				mainTitle = fmt.Sprintf("🔊 Waktu Adzan %s", label)
				mainDescription = fmt.Sprintf(
					"📢 @here Adzan %s telah berkumandang pada **pukul %s WIB**. Mari bersiap untuk menunaikan shalat.",
					label,
					shalatTime.Format("15:04"),
				)
			} else {
				formattedDurationAfter := fmt.Sprintf("%d menit", int(absDuration.Minutes()))
				mainTitle = fmt.Sprintf("✅ Waktu Shalat %s Sedang Berlangsung", label)
				mainDescription = fmt.Sprintf(
					"📢 @here Shalat %s telah dimulai pada **pukul %s WIB**. Notifikasi ini dikirim **%s setelah adzan** sebagai pengingat untuk menunaikan shalat.",
					label,
					shalatTime.Format("15:04"),
					formattedDurationAfter,
				)
			}
		}

		zoneName, _ := now.Zone()

		mainEmbed := models.DiscordEmbed{
			Color:       embedColor,
			Title:       mainTitle,
			Description: mainDescription,
			Timestamp:   time.Now().Format(time.RFC3339),
			Footer: &map[string]string{
				"text": fmt.Sprintf(
					"📍 %s | %s %s",
					target.City,
					time.Now().In(loc).Format("02/01/2006 15:04"),
					zoneName,
				),
			},
		}

		isAdzanForQuote := isExactlyTime && label != "Dhuha"
		quoteEmbed, categoryShalat := discord.GetPrayerQuoteEmbed(label, isBefore, isAdzanForQuote, isAfter)

		if label == models.Dzuhur && now.Weekday() == time.Friday && (isBefore || isExactlyTime) {
			categoryShalat = models.Jumat
		}

		disclaimerEmbed := models.DiscordEmbed{
			Color: 10070709,
			Title: "✨ Catatan",
			Description: "Pesan ini merupakan pengingat otomatis. Kami menghormati setiap keyakinan dan praktik ibadah. " +
				"Notifikasi ini ditujukan bagi **Rekan Muslim yang menunaikan ibadah shalat**.",
		}

		coordinationEmbed := models.DiscordEmbed{
			Color: 15158332,
			Title: "📢 Koordinasi Waktu Ibadah",
			Description: "Mohon izin untuk saling mengingatkan dan menyesuaikan aktivitas agar ibadah shalat dapat ditunaikan tepat waktu. " +
				"Terima kasih atas perhatian dan kerja samanya.",
		}

		embedsToSend := []models.DiscordEmbed{
			mainEmbed,
			quoteEmbed,
			coordinationEmbed,
			disclaimerEmbed,
		}

		result := discord.SendWebhook(label, embedsToSend, target, categoryShalat, testingMode)

		logMessage := fmt.Sprintf(
			"%s pada %s (%s). Tipe: %s. Webhook: %s",
			label,
			notificationTime.Format("15:04"),
			zoneName,
			func() string {
				if isBefore {
					return "SEBELUM"
				} else if isExactlyTime {
					return "SAAT MULAI"
				}
				return "SETELAH"
			}(),
			result.Name,
		)

		if result.Error != nil {
			if strings.Contains(result.Error.Error(), "DILEWATI: Divisi IT") {
				fmt.Printf("⏩ DILEWATI %s (AKHIR PEKAN) untuk Webhook: %s\n", label, result.Name)
			} else {
				fmt.Printf("❌ Gagal mengirim pesan untuk %s. Error: %v\n", logMessage, result.Error)
			}
		} else {
			fmt.Printf("✅ Berhasil mengirim pesan untuk %s. Total embeds: %d\n", logMessage, len(embedsToSend))
		}

	} else {
		zoneName, _ := now.Zone()
		fmt.Printf(
			"⏩ DILEWATI %s (%v Waktu %s). Waktu Shalat: %s. Waktu Notif: %s. Sekarang: %s %s\n",
			label,
			absDuration,
			func() string {
				if isBefore {
					return "SEBELUM"
				}
				return "SETELAH"
			}(),
			shalatTime.Format("15:04"),
			notificationTime.Format("15:04:05"),
			now.Format("15:04:05"),
			zoneName,
		)
	}
}
