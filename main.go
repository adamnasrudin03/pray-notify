package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	// 🔥 TAMBAH: Untuk fungsi random
	"math/rand"

	"github.com/hablullah/go-prayer"
)

// Struktur baru untuk menyimpan URL, Nama/Label Webhook, dan Mentions CC
type WebhookTarget struct {
	Name     string
	URL      string
	Mentions string // Field untuk CC (@user @user)
}

// Struktur untuk menampung data Embed yang akan dikirim
type DiscordEmbed struct {
	Color       int                `json:"color"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Timestamp   string             `json:"timestamp"`
	Footer      *map[string]string `json:"footer,omitempty"`
}

// Struktur baru untuk quote
type Quote struct {
	Text      string
	Reference string
}

var (
	// Koordinat Kebon Jeruk, Jakarta Barat
	lat      = -6.183
	lon      = 106.784
	timezone = "Asia/Jakarta"

	// Daftar Webhook Target
	webhookTargets = []WebhookTarget{
		{
			Name:     "Divisi IT",
			URL:      "https://discord.com/api/webhooks/1443879495324008488/z-jzPBXv6vBvE45-QFrV-wwcIp8Qz9GtMw0bSQeKh1ltjHXZjcUPmjN3pMpS7Jg4yZJ0",
			Mentions: "@abi_lp @tedja_putra_lp @ade_muhammad_lp @bian_kusumah_lp",
		},
		{
			Name:     "Apran Channel",
			URL:      "https://discordapp.com/api/webhooks/1443512213842296904/EHX00Lvi8Hmqxzc1WsgJK9pyDsJzAmXWhexbqSsJnez8KSzaAHaBy_8UQIAgtRQyQEix",
			Mentions: "",
		},
	}
)

// Koleksi Quotes berdasarkan kategori waktu
var randomQuotes = map[string][]Quote{
	// Quotes untuk notifikasi SEBELUM Adzan (Mendorong tepat waktu/persiapan)
	"SebelumAdzan": {
		{"Jagalah shalat lima waktu, terutama shalat yang pertengahan (Ashar) dan berdirilah untuk Allah (dalam shalatmu) dengan khusyuk.", "QS. Al-Baqarah: 238"},
		{"Shalat di awal waktu itu adalah amalan yang paling dicintai oleh Allah.", "HR. Bukhari dan Muslim"},
		{"Bersegeralah mendirikan shalat sebelum sibuk dengan urusan dunia.", "Nasehat Islam"},
		{"Sesungguhnya shalat pada waktunya adalah kewajiban yang telah ditentukan atas orang-orang yang beriman.", "QS. An-Nisa: 103"},
	},
	// Quotes untuk notifikasi SAAT Adzan (Spesifik per waktu salat)
	"SaatAdzan": {
		{"Sesungguhnya shalat Subuh itu disaksikan (oleh malaikat-malaikat malam dan malaikat-malaikat siang).", "QS. Al-Isra: 78 (Subuh)"},
		{"Jagalah (peliharalah) segala shalatmu dan (peliharalah) shalat wustha. Berdirilah untuk Allah (dalam shalatmu) dengan khusyuk.", "QS. Al-Baqarah: 238 (Dzuhur/Ashar)"},
		{"Barangsiapa meninggalkan shalat Ashar, maka seolah-olah dia telah kehilangan keluarga dan hartanya.", "HR. Bukhari dan Muslim (Ashar)"},
		{"Allah tidak akan menyia-nyiakan pahala orang-orang yang berbuat baik.", "QS. Ali Imran: 171 (Maghrib)"},
		{"Barangsiapa yang shalat Isya' berjama'ah, maka seolah-olah ia telah shalat setengah malam.", "HR. Muslim (Isya)"},
	},
	// Quotes untuk notifikasi SETELAH Adzan (Mendorong segera menunaikan/mendirikan)
	"SetelahAdzan": {
		{"Dirikanlah salat dan tunaikanlah zakat. Kebaikan apa pun yang kamu kerjakan untuk dirimu, kamu akan mendapatkannya (pahala) di sisi Allah.", "QS. Al-Baqarah: 110"},
		{"Sesungguhnya orang-orang yang beriman itu hanyalah orang-orang yang apabila disebutkan Allah gemetarlah hati mereka, dan apabila dibacakan kepada mereka ayat-ayat-Nya, bertambahlah iman mereka.", "QS. Al-Anfal: 2 (Mengingatkan kewajiban setelah mendengar Adzan)"},
		{"Dan mintalah pertolongan (kepada Allah) dengan sabar dan salat. Dan sesungguhnya yang demikian itu sungguh berat, kecuali bagi orang-orang yang khusyuk.", "QS. Al-Baqarah: 45"},
		{"Peliharalah semua salat(mu), dan (peliharalah) salat wusṭā. Berdirilah karena Allah (dalam salatmu) dengan khusyuk.", "QS. Al-Baqarah: 238 (Pengulangan penekanan kewajiban)"},
	},
}

// WebhookResult menyimpan hasil pengiriman ke satu webhook
type WebhookResult struct {
	Name  string
	Error error
}

// getRandomQuote memilih quote secara acak berdasarkan kategori waktu
func getRandomQuote(category string) Quote {
	quotes, ok := randomQuotes[category]
	if !ok || len(quotes) == 0 {
		return Quote{Text: "Dirikanlah salat dan tunaikanlah zakat.", Reference: "QS. Al-Baqarah: 110"} // Default
	}

	// Inisialisasi rand (hanya diperlukan sekali, tapi Go akan menangani seed dengan baik)
	// rand.Seed(time.Now().UnixNano()) // Dihapus karena Go modern menangani seed secara default

	randomIndex := rand.Intn(len(quotes))
	return quotes[randomIndex]
}

// getPrayerQuoteEmbed membuat embed untuk kutipan ayat/hadits berdasarkan status waktu notifikasi
func getPrayerQuoteEmbed(label string, isBefore bool, isAdzan bool, isAfter bool) DiscordEmbed {
	var category string
	if isAdzan {
		category = "SaatAdzan"
	} else if isBefore {
		category = "SebelumAdzan"
	} else if isAfter {
		category = "SetelahAdzan"
	} else {
		category = "SaatAdzan" // Default jika ada masalah logika
	}

	selected := getRandomQuote(category)

	return DiscordEmbed{
		Color:       3447003, // Warna Biru (Blue)
		Title:       "📖 Kutipan Ayat/Hadits",
		Description: fmt.Sprintf("> *%s*\n\n— %s", selected.Text, selected.Reference),
	}
}

// sendWebhook mengirimkan pesan ke satu WebhookTarget
func sendWebhook(embeds []DiscordEmbed, target WebhookTarget, shalatTimeLabel string) WebhookResult {

	result := WebhookResult{Name: target.Name}
	now := time.Now()

	// 🛑 VALIDASI 1: Cek Hari Libur/Weekend untuk Webhook "Divisi IT"
	if target.Name == "Divisi IT" && (now.Weekday() == time.Saturday || now.Weekday() == time.Sunday) {
		result.Error = fmt.Errorf("webhook SKIPPED: Divisi IT is not active on weekends (Saturday/Sunday)")
		return result
	}

	// Payload utama Discord Webhook
	payload := map[string]any{
		"embeds": embeds,
	}

	body, err := json.Marshal(payload)

	if err != nil {
		result.Error = fmt.Errorf("marshalling failed: %w", err)
		return result
	}

	resp, err := http.Post(target.URL, "application/json", bytes.NewBuffer(body))

	if err != nil {
		result.Error = fmt.Errorf("HTTP post failed: %w", err)
	} else {
		resp.Body.Close()
		if resp.StatusCode >= 300 {
			result.Error = fmt.Errorf("status code %d", resp.StatusCode)
		}
	}
	return result
}

// checkAndSend notifikasi
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

	// Tentukan apakah notifikasi tepat diwaktu adzan (offset 0 menit)
	isExactlyAdzan := absDuration.Minutes() < 1
	// Tentukan apakah notifikasi setelah adzan (offset positif)
	isAfter := !isBefore && !isExactlyAdzan

	if notificationTime.After(fiveMinutesAgo) && (notificationTime.Before(now) || notificationTime.Equal(now)) {

		var mainTitle string
		var mainDescription string
		var embedColor int // Warna dalam format desimal

		if isBefore {
			durationUntilShalat := shalatTime.Sub(now).Truncate(time.Second)
			embedColor = 16763904 // Warna Kuning (Pengingat)

			if durationUntilShalat > 0 {
				mainTitle = fmt.Sprintf("🕌 Adzan %s Akan Tiba!", label)
				mainDescription = fmt.Sprintf("Waktu Shalat akan berkumandang dalam **%s** (Pukul **%s**). Segera siapkan diri.",
					durationUntilShalat.Truncate(time.Minute),
					shalatTime.Format("15:04"),
				)
			} else {
				// Ini berarti notif sebelum adzan terlewat, tapi masih dalam 5 menit jendela
				mainTitle = fmt.Sprintf("🚨 ADZAN %s TELAH TIBA (Delayed)!", label)
				mainDescription = fmt.Sprintf("Adzan telah berkumandang (Pukul **%s**). Mohon maaf, pengingat **%v** sebelum Adzan tertunda. Segera Tunaikan Shalat.",
					shalatTime.Format("15:04"),
					absDuration,
				)
			}

		} else {
			embedColor = 5814783 // Warna Hijau (Waktu Salat Berlangsung)

			// Logika untuk Saat Adzan dan Setelah Adzan
			if isExactlyAdzan {
				mainTitle = fmt.Sprintf("📢 WAKTU ADZAN %s", label)
				mainDescription = fmt.Sprintf("Adzan %s baru saja berkumandang (Pukul **%s**). Saatnya untuk berwudhu dan bersiap menunaikan Shalat.",
					label,
					shalatTime.Format("15:04"),
				)
			} else { // isAfter
				mainTitle = fmt.Sprintf("✅ WAKTU SALAT %s SEDANG BERLANGSUNG", label)
				mainDescription = fmt.Sprintf("Shalat %s sudah dimulai (Pukul **%s**). Ini adalah notifikasi **%v** setelah Adzan. Saatnya Tunaikan Shalat.",
					label,
					shalatTime.Format("15:04"),
					absDuration,
				)
			}
		}

		wibLoc := time.FixedZone("WIB", 7*3600)

		// 1. EMBED UTAMA (Waktu Salat)
		mainEmbed := DiscordEmbed{
			Color:       embedColor,
			Title:       mainTitle,
			Description: mainDescription,
			Timestamp:   time.Now().Format(time.RFC3339),
			Footer: &map[string]string{
				"text": fmt.Sprintf("Lokasi: Kebon Jeruk, Jakarta Barat | Dihitung: %s WIB", time.Now().In(wibLoc).Format("02/01/2006 15:04")),
			},
		}

		// 2. EMBED KUTIPAN (Ayat/Hadits)
		// 🔥 MENGGUNAKAN LOGIKA PEMILIHAN RANDOM BARU
		quoteEmbed := getPrayerQuoteEmbed(label, isBefore, isExactlyAdzan, isAfter)

		// 4. EMBED DISCLAIMER
		disclaimerEmbed := DiscordEmbed{
			Color:       10070709, // Warna Ungu (Purple/Blurple)
			Title:       "✨ Catatan Santun",
			Description: "Pesan ini adalah pengingat otomatis. Kami menghormati keyakinan dan praktik ibadah Anda. Notifikasi ini ditujukan khusus bagi **Rekan Muslim yang menunaikan salat**.",
		}

		// Loop dan Kirim ke Semua Webhook
		for _, target := range webhookTargets {

			// Siapkan array embed (min 3)
			embedsToSend := []DiscordEmbed{mainEmbed, quoteEmbed}

			// 🛑 VALIDASI 2: Cek waktu Salat untuk Mentions (Hanya Dzuhur, Ashar, Maghrib)
			shouldMention := false
			if target.Mentions != "" {
				switch label {
				case "Dzuhur", "Ashar", "Maghrib":
					shouldMention = true
				}
			}

			// 3. EMBED KOORDINASI (Conditional Mentions)
			if shouldMention {

				// Gabungkan Mentions di Description Embed
				coordinationDescription := fmt.Sprintf(
					"🔊 **Kepada Bapak/Ibu Rekan Kerja %s**, mohon izin untuk saling mengingatkan dan mengkoordinasikan agar ibadah salat dapat ditunaikan dengan tepat waktu. Terima kasih banyak atas perhatian dan kerjasamanya.",
					target.Mentions,
				)

				coordinationEmbed := DiscordEmbed{
					Color:       15158332,                    // Warna Merah
					Title:       "📢 Koordinasi Waktu Ibadah", // Judul baru
					Description: coordinationDescription,
				}
				embedsToSend = append(embedsToSend, coordinationEmbed)
			}

			// Tambahkan Disclaimer sebagai embed terakhir
			embedsToSend = append(embedsToSend, disclaimerEmbed)

			// Panggil sendWebhook

			// Hapus mentions dari payload content jika tidak seharusnya di mention (Subuh/Isya)
			currentTarget := target
			if !shouldMention {
				currentTarget.Mentions = ""
			}

			result := sendWebhook(embedsToSend, currentTarget, label)

			// Log hasil pengiriman
			logMessage := fmt.Sprintf("%s at %s. Type: %s. Webhook: %s",
				label,
				notificationTime.Format("15:04"),
				func() string {
					if isBefore {
						return "SEBELUM"
					} else if isExactlyAdzan {
						return "SAAT ADZAN"
					} else {
						return "SETELAH"
					}
				}(),
				result.Name, // Nama Webhook
			)

			if result.Error != nil {
				// Cek jika errornya adalah SKIPPED WEEKEND
				if strings.Contains(result.Error.Error(), "SKIPPED: Divisi IT") {
					fmt.Printf("⏩ SKIPPED %s (WEEKEND) for Webhook: %s\n", label, result.Name)
				} else {
					fmt.Printf("❌ Error sending multi-EMBED notification for %s. Error: %v\n", logMessage, result.Error)
				}
			} else {
				// Tambahkan indikator mentions di log
				mentionsInfo := ""
				if shouldMention && target.Mentions != "" {
					mentionsInfo = " (with CC)"
				}
				fmt.Printf("✅ Sent multi-EMBED notification for %s. Embeds: %d%s.\n", logMessage, len(embedsToSend), mentionsInfo)
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
	// PENTING: Seed Randomizer hanya sekali di main
	rand.Seed(time.Now().UnixNano())

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

	offsets := []time.Duration{}
	for m := -20; m <= 10; m += 2 {
		offsets = append(offsets, time.Duration(m)*time.Minute)
	}

	for _, prayer := range prayerTimes {
		for _, offset := range offsets {
			notificationTime := prayer.t.Add(offset)
			checkAndSend(notificationTime, prayer.label, prayer.t, offset)
		}
	}
}
