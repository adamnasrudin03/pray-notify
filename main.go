package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hablullah/go-prayer"
)

type WebhookTarget struct {
	Name string
	URL  string
}

type DiscordEmbed struct {
	Color       int                `json:"color"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Timestamp   string             `json:"timestamp"`
	Footer      *map[string]string `json:"footer,omitempty"`
}

type Quote struct {
	Text      string
	Reference string
}

const (
	DivisiIT  = "Divisi IT"
	MyChannel = "My Channel"

	Dhuha   = "Dhuha"
	Subuh   = "Subuh"
	Dzuhur  = "Dzuhur"
	Ashar   = "Ashar"
	Maghrib = "Maghrib"
	Isya    = "Isya"

	SebelumAdzan = "SebelumAdzan"
	AdzanSubuh   = "AdzanSubuh"
	AdzanDzuhur  = "AdzanDzuhur"
	AdzanAshar   = "AdzanAshar"
	AdzanMaghrib = "AdzanMaghrib"
	AdzanIsya    = "AdzanIsya"
	SetelahAdzan = "SetelahAdzan"
)

var (
	// Koordinat Lion Parcel Kedoya, Kebon Jeruk, Jakarta Barat
	lat      = -6.189541191231946
	lon      = 106.76490428724479
	timezone = "Asia/Jakarta"

	// Daftar Webhook Target
	webhookTargets = []WebhookTarget{
		{
			Name: DivisiIT,
			URL:  "https://discord.com/api/webhooks/1443879495324008488/z-jzPBXv6vBvE45-QFrV-wwcIp8Qz9GtMw0bSQeKh1ltjHXZjcUPmjN3pMpS7Jg4yZJ0",
		},
		{
			Name: MyChannel,
			URL:  "https://discordapp.com/api/webhooks/1443512213842296904/EHX00Lvi8Hmqxzc1WsgJK9pyDsJzAmXWhexbqSsJnez8KSzaAHaBy_8UQIAgtRQyQEix",
		},
	}
)

// Koleksi Quotes berdasarkan kategori waktu (Tidak ada perubahan)
var randomQuotes = map[string][]Quote{
	// 🔥 Quote Khusus Dhuha
	Dhuha: {
		{"Dua rakaat shalat Dhuha setara dengan 360 sedekah yang diwajibkan atas setiap persendian tubuhmu.", "HR. Muslim"},
		{"Barangsiapa shalat Dhuha dua belas rakaat, niscaya Allah membangunkan baginya istana di surga.", "HR. Tirmidzi & Ibnu Majah"},
		{"Shalat Dhuha adalah shalat awwabin (orang-orang yang kembali taat), ketika anak unta kepanasan.", "HR. Muslim"},
		{"Jadikan awal harimu penuh berkah dengan Shalat Dhuha. Rezeki menanti para hamba yang bertakwa.", "Nasehat Islam"},
	},
	// Quotes untuk notifikasi SEBELUM Adzan (Mendorong tepat waktu/persiapan)
	SebelumAdzan: {
		{"Jagalah shalat lima waktu, terutama shalat yang pertengahan (Ashar) dan berdirilah untuk Allah (dalam shalatmu) dengan khusyuk.", "QS. Al-Baqarah: 238"},
		{"Shalat di awal waktu itu adalah amalan yang paling dicintai oleh Allah.", "HR. Bukhari dan Muslim"},
		{"Bersegeralah mendirikan shalat sebelum sibuk dengan urusan dunia.", "Nasehat Islam"},
		{"Sesungguhnya shalat pada waktunya adalah kewajiban yang telah ditentukan atas orang-orang yang beriman.", "QS. An-Nisa: 103"},
	},
	AdzanSubuh: {
		{"Sesungguhnya shalat Subuh itu disaksikan (oleh malaikat-malaikat malam dan malaikat-malaikat siang).", "QS. Al-Isra: 78"},
		{"Barangsiapa shalat Shubuh berjama’ah, maka dia berada dalam jaminan Allah.", "HR. Muslim"},
	},
	AdzanDzuhur: {
		{"Jagalah (peliharalah) segala shalatmu dan (peliharalah) shalat wustha. Berdirilah untuk Allah (dalam shalatmu) dengan khusyuk.", "QS. Al-Baqarah: 238"},
		{"Apabila kamu telah menyelesaikan shalat, ingatlah Allah ketika kamu berdiri, ketika kamu duduk dan ketika kamu berbaring.", "QS. An-Nisa: 103"},
	},
	AdzanAshar: {
		{"Barangsiapa meninggalkan shalat Ashar, maka seolah-olah dia telah kehilangan keluarga dan hartanya.", "HR. Bukhari dan Muslim"},
		{"Shalat wustha adalah shalat Ashar. Siapa yang menjaganya, ia telah mendapat keberuntungan besar.", "HR. Tirmidzi"},
	},
	AdzanMaghrib: {
		{"Jika datang malam dari arah sini (timur) dan pergi siang dari arah sini (barat), serta terbenam matahari, maka sungguh orang yang puasa telah berbuka (waktunya Maghrib).", "HR. Bukhari dan Muslim"},
		{"Allah tidak akan menyia-nyiakan pahala orang-orang yang berbuat baik.", "QS. Ali Imran: 171"},
	},
	AdzanIsya: {
		{"Barangsiapa yang shalat Isya' berjama'ah, maka seolah-olah ia telah shalat setengah malam.", "HR. Muslim"},
		{"Shalat yang paling berat bagi orang munafik adalah shalat Isya' dan shalat Subuh.", "HR. Bukhari dan Muslim"},
	},
	// Quotes untuk notifikasi SETELAH Adzan (Mendorong segera menunaikan/mendirikan)
	SetelahAdzan: {
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

	randomIndex := rand.Intn(len(quotes))
	return quotes[randomIndex]
}

// getPrayerQuoteEmbed membuat embed untuk kutipan ayat/hadits berdasarkan status waktu notifikasi
func getPrayerQuoteEmbed(label string, isBefore bool, isAdzan bool, isAfter bool) (result DiscordEmbed, categoryShalat string) {
	var category string

	// 🔥 LOGIKA BARU UNTUK QUOTE SPESIFIK SAAT ADZAN
	if label == Dhuha {
		category = Dhuha
	} else if isAdzan {
		// Jika tepat waktu Adzan, gunakan kategori spesifik (e.g., "AdzanSubuh")
		category = "Adzan" + label
	} else if isBefore {
		category = SebelumAdzan
	} else if isAfter {
		category = SetelahAdzan
	} else {
		// Fallback: Untuk waktu Adzan yang tidak terdeteksi labelnya (seharusnya tidak terjadi)
		category = "Adzan" + label
	}

	selected := getRandomQuote(category)

	return DiscordEmbed{
		Color:       3447003, // Warna Biru (Blue)
		Title:       "📖 Kutipan",
		Description: fmt.Sprintf("> *%s*\n\n— %s", selected.Text, selected.Reference),
	}, category
}

func BuildShalatMessage(categoryShalat string) string {

	var content string

	switch categoryShalat {

	case SebelumAdzan:
		content = `
$ system-event trigger shalat_prepare
----------------------------------------
 event: persiapan_shalat
 status: reminder
 message: Waktu shalat hampir tiba
----------------------------------------
>> Persiapkan diri & hentikan aktivitas.
`

	case SetelahAdzan:
		content = `
[Pipeline] stage("Adzan Alert")
[INFO] Waktu shalat telah masuk...
[INFO] Triggering: adzan_notification()
[INFO] Actions:
        - alert_user("Shalat time!")
        - prepare_wudhu()
        - pause_non_critical_tasks()
[Pipeline] End stage("Adzan Alert")
`

	case Dhuha:
		content = `
[Pipeline] stage("Dhuha Alert")
[INFO] Waktu Shalat Dhuha telah tiba...
[INFO] Triggering: dhuha_notification()
[INFO] Actions:
        - alert_user("Shalat Dhuha time!")
        - prepare_wudhu()
        - pause_non_critical_tasks()
[Pipeline] End stage("Dhuha Alert")
`

	default:
		content = `
[Pipeline] stage("Adzan Alert")
[INFO] Waktu shalat telah masuk...
[INFO] Triggering: adzan_notification()
[INFO] Actions:
        - pause_tasks()
        - prepare_wudhu()
        - notify_team("Shalat time!")
[Pipeline] End stage("Adzan Alert")
`
	}

	// Bungkus ke Discord code block ```txt
	return fmt.Sprintf("```txt\n%s\n```", content)
}

// sendWebhook mengirimkan pesan ke satu WebhookTarget
func sendWebhook(embeds []DiscordEmbed, target WebhookTarget, categoryShalat string) WebhookResult {

	result := WebhookResult{Name: target.Name}
	now := time.Now()

	// 🛑 VALIDASI 1: Cek Hari Libur/Weekend untuk Webhook "Divisi IT"
	if target.Name == DivisiIT && (now.Weekday() == time.Saturday || now.Weekday() == time.Sunday) {
		result.Error = fmt.Errorf("webhook SKIPPED: Divisi IT is not active on weekends (Saturday/Sunday)")
		return result
	}
	message := BuildShalatMessage(categoryShalat)

	message += " \n\n\n"
	// Payload utama Discord Webhook
	payload := map[string]any{
		"content": message,
		"embeds":  embeds,
	}

	body, err := json.Marshal(payload)

	if err != nil {
		result.Error = fmt.Errorf("marshalling failed: %w", err)
		return result
	}

	// Kirim HANYA SATU request HTTP POST
	resp, err := http.Post(target.URL, "application/json", bytes.NewBuffer(body))

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		result.Error = fmt.Errorf("HTTP post failed: %w", err)
	} else {
		if resp.StatusCode >= 300 {
			result.Error = fmt.Errorf("status code %d", resp.StatusCode)
			fmt.Printf("Error sending webhook to %s: %+v\n", target.Name, bodyString)
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

	// Tentukan apakah notifikasi tepat diwaktu "Adzan" (offset 0 menit atau untuk Dhuha)
	isExactlyTime := absDuration.Minutes() < 1
	// Tentukan apakah notifikasi setelah "Adzan" (offset positif)
	isAfter := !isBefore && !isExactlyTime

	if notificationTime.After(fiveMinutesAgo) && (notificationTime.Before(now) || notificationTime.Equal(now)) {

		var mainTitle string
		var mainDescription string
		var embedColor int // Warna dalam format desimal

		// 🔥 LOGIKA KHUSUS UNTUK DHUHA
		if label == "Dhuha" {
			embedColor = 16763904 // Warna Kuning/Emas untuk Sunnah
			mainTitle = "✨ WAKTU SHALAT DHUHA DIMULAI!"
			mainDescription = fmt.Sprintf("📢 @here Shalat Dhuha sudah dapat ditunaikan mulai Pukul **%s**. Mari raih keutamaan sunnah untuk kelancaran rezeki dan aktivitas harian.", shalatTime.Format("15:04"))

		} else if isBefore {
			durationUntilShalat := shalatTime.Sub(now).Truncate(time.Second)
			minutes := int(durationUntilShalat.Minutes())
			embedColor = 16763904 // Warna Kuning (Pengingat)

			if durationUntilShalat > 0 {
				mainTitle = fmt.Sprintf("🕌 Adzan %s Akan Tiba!", label)
				mainDescription = fmt.Sprintf("📢 @here Waktu Shalat akan berkumandang dalam **%d menit** (Pukul **%s**). Segera siapkan diri.",
					minutes,
					shalatTime.Format("15:04"),
				)
			} else {
				mainTitle = fmt.Sprintf("🚨 ADZAN %s TELAH TIBA (Delayed)!", label)
				mainDescription = fmt.Sprintf("📢 @here Adzan telah berkumandang (Pukul **%s**). Mohon maaf, pengingat **%v** sebelum Adzan tertunda. Segera Tunaikan Shalat.",
					shalatTime.Format("15:04"),
					absDuration,
				)
			}

		} else {
			embedColor = 5814783 // Warna Hijau (Waktu Salat Berlangsung)

			if isExactlyTime {
				mainTitle = fmt.Sprintf("📢 WAKTU ADZAN %s", label)
				mainDescription = fmt.Sprintf("🔊 @here Adzan %s baru saja berkumandang (Pukul **%s**). Saatnya untuk berwudhu dan bersiap menunaikan Shalat.",
					label,
					shalatTime.Format("15:04"),
				)
			} else { // isAfter
				mainTitle = fmt.Sprintf("✅ WAKTU SHALAT %s SEDANG BERLANGSUNG", label)
				mainDescription = fmt.Sprintf("📢 @here Shalat %s sudah dimulai (Pukul **%s**). Ini adalah notifikasi **%v** setelah Adzan. Saatnya Tunaikan Shalat.",
					label,
					shalatTime.Format("15:04"),
					absDuration,
				)
			}
		}

		wibLoc := time.FixedZone("WIB", 7*3600)

		mainEmbed := DiscordEmbed{
			Color:       embedColor,
			Title:       mainTitle,
			Description: mainDescription,
			Timestamp:   time.Now().Format(time.RFC3339),
			Footer: &map[string]string{
				"text": fmt.Sprintf("Lokasi: Kebon Jeruk, Jakarta Barat | Dihitung: %s WIB", time.Now().In(wibLoc).Format("02/01/2006 15:04")),
			},
		}

		isAdzanForQuote := isExactlyTime && label != "Dhuha"
		quoteEmbed, categoryShalat := getPrayerQuoteEmbed(label, isBefore, isAdzanForQuote, isAfter)

		disclaimerEmbed := DiscordEmbed{
			Color:       10070709, // Warna Ungu (Purple/Blurple)
			Title:       "✨ Catatan",
			Description: "Pesan ini adalah pengingat otomatis. Kami menghormati keyakinan dan praktik ibadah Anda. Notifikasi ini ditujukan khusus bagi **Rekan Muslim yang menunaikan ibadah shalat**.",
		}

		// Loop dan Kirim ke Semua Webhook
		for _, target := range webhookTargets {

			// Siapkan array embed (min 3)
			embedsToSend := []DiscordEmbed{mainEmbed, quoteEmbed}

			currentTarget := target

			coordinationDescription := "Mohon izin untuk saling mengingatkan dan mengkoordinasikan agar ibadah shalat dapat ditunaikan dengan tepat waktu. Terima kasih banyak atas perhatian dan kerjasamanya."

			coordinationEmbed := DiscordEmbed{
				Color:       15158332, // Warna Merah
				Title:       "📢 Koordinasi Waktu Ibadah",
				Description: coordinationDescription,
			}
			embedsToSend = append(embedsToSend, coordinationEmbed)

			// Tambahkan Disclaimer sebagai embed terakhir
			embedsToSend = append(embedsToSend, disclaimerEmbed)

			// Panggil sendWebhook
			// Perubahan sudah di sendWebhook untuk menggabungkan mention dan embed di satu pesan
			result := sendWebhook(embedsToSend, currentTarget, categoryShalat)

			// Log hasil pengiriman
			logMessage := fmt.Sprintf("%s at %s. Type: %s. Webhook: %s",
				label,
				notificationTime.Format("15:04"),
				func() string {
					if isBefore {
						return "SEBELUM"
					} else if isExactlyTime {
						return "SAAT MULAI"
					} else {
						return "SETELAH"
					}
				}(),
				result.Name, // Nama Webhook
			)

			if result.Error != nil {
				if strings.Contains(result.Error.Error(), "SKIPPED: Divisi IT") {
					fmt.Printf("⏩ SKIPPED %s (WEEKEND) for Webhook: %s\n", label, result.Name)
				} else {
					fmt.Printf("❌ Error sending message for %s. Error: %v\n", logMessage, result.Error)
				}
			} else {
				fmt.Printf("✅ Successfully Sent message for %+v. Embeds: %+v.\n", logMessage, len(embedsToSend))
			}
		}

	} else {
		// Logika SKIPPED
		fmt.Printf(
			"⏩ SKIPPED %s (%v %s Waktu). Waktu Mulai: %s. Waktu Notif Seharusnya: %s. Waktu Saat Ini: %s WIB.\n",
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
	sunrise := schedule.Sunrise.In(loc) // 🔥 Ambil waktu Sunrise
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
		{Subuh, fajr},
		{Dhuha, dhuhaStart}, // 🔥 Tambahkan Dhuha
		{Dzuhur, dhuhr},
		{Ashar, asr},
		{Maghrib, maghrib},
		{Isya, isha},
	}

	offsetsMain := []time.Duration{}
	for m := -30; m <= 30; m += 5 {
		offsetsMain = append(offsetsMain, time.Duration(m)*time.Minute)
	}

	offsetsCaseSpecific := []time.Duration{}
	for m := -30; m <= 30; m += 1 {
		offsetsCaseSpecific = append(offsetsCaseSpecific, time.Duration(m)*time.Minute)
	}

	for _, prayer := range prayerTimes {
		currentOffsets := offsetsMain
		if prayer.label == Dhuha {
			currentOffsets = offsetsCaseSpecific
		}

		for _, offset := range currentOffsets {
			notificationTime := prayer.t.Add(offset)

			// 🔥 KHUSUS Dhuha: Hanya kirim notifikasi saat offset adalah 0 menit
			if prayer.label == Dhuha && offset.Minutes() != 0 {
				continue // Lewati reminder sebelum dan sesudah Dhuha
			}

			checkAndSend(notificationTime, prayer.label, prayer.t, offset)
		}
	}
}
