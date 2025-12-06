package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/hablullah/go-prayer"
)

type WebhookTarget struct {
	Name string
	URL  string
	// Data Lokasi Spesifik
	Lat  float64
	Lon  float64
	Tz   string // Timezone string, misal "Asia/Jakarta"
	City string // Nama kota/lokasi untuk ditampilkan di footer
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
	MyUsername = "Adam Nasrudin (@adamnasrudin)"
	MyAvatar   = "https://media.discordapp.net/attachments/936166044748369930/1446540243993104546/DALLE_2024-10-28_14.40.34_-_A_fun_and_energetic_Discord_profile_picture_featuring_a_colorful_and_animated_character_with_a_playful_expression._The_character_has_vibrant_colors_b.png?ex=69345b15&is=69330995&hm=e3ee118fbc4faa7188af080d7e35b99031892b0f44ebe85205121d6836ba0972&=&format=webp&quality=lossless&width=810&height=810"

	Test        = "Testing"
	TestAvaatar = "https://media.discordapp.net/attachments/936166044748369930/1446540244781367441/DALLE_2024-10-28_14.42.11_-_A_fun_and_colorful_Discord_banner_to_match_an_energetic_profile_picture_featuring_a_playful_and_vibrant_scene._The_banner_includes_a_gradient_backgro.png?ex=69345b15&is=69330995&hm=b45fe7ca9ada4123096cadf404bb03ad8380d0ad212c3e5c2410bff233509917&=&format=webp&quality=lossless&width=1416&height=809"

	DivisiIT  = "Divisi IT"
	MyChannel = "My Channel"

	Dhuha   = "dhuha"
	Subuh   = "subuh"
	Dzuhur  = "dzuhur"
	Ashar   = "ashar"
	Maghrib = "maghrib"
	Isya    = "isya"

	SebelumAdzan = "SebelumAdzan"
	AdzanSubuh   = "AdzanSubuh"
	AdzanDzuhur  = "AdzanDzuhur"
	AdzanAshar   = "AdzanAshar"
	AdzanMaghrib = "AdzanMaghrib"
	AdzanIsya    = "AdzanIsya"
	SetelahAdzan = "SetelahAdzan"
)

var (
	testingMode bool
	// Koordinat Lion Parcel Kedoya, Kebon Jeruk, Jakarta Barat
	latDefault    = -6.189541191231946
	lonDefault    = 106.76490428724479
	tzAsiaJakarta = "Asia/Jakarta"
	cityDefault   = "Kebon Jeruk, Jakarta Barat" // ** BARU **

	// Daftar Webhook Target dengan lokasi spesifik (tidak diubah)
	webhookTargets = []WebhookTarget{
		{
			Name: DivisiIT,
			URL:  "https://discord.com/api/webhooks/1443879495324008488/z-jzPBXv6vBvE45-QFrV-wwcIp8Qz9GtMw0bSQeKh1ltjHXZjcUPmjN3pMpS7Jg4yZJ0",
			Lat:  latDefault,
			Lon:  lonDefault,
			Tz:   tzAsiaJakarta,
			City: cityDefault,
		},
		{
			Name: MyChannel,
			URL:  "https://discordapp.com/api/webhooks/1443512213842296904/EHX00Lvi8Hmqxzc1WsgJK9pyDsJzAmXWhexbqSsJnez8KSzaAHaBy_8UQIAgtRQyQEix",
			Lat:  latDefault,
			Lon:  lonDefault,
			Tz:   tzAsiaJakarta,
			City: cityDefault,
		},
		{
			Name: "Yusuf",
			URL:  "https://discord.com/api/webhooks/1445697235231838208/UaQeKj2rS9tXdoDdB-SIMna_yfd2MKlPBlyOPWfUMIQlMvn0RM9CHjAOJ_XiB6dcygQx",
			Lat:  latDefault,
			Lon:  lonDefault,
			Tz:   tzAsiaJakarta,
			City: cityDefault,
		},
	}
)

var randomQuotes = map[string][]Quote{
	Dhuha: {
		{"Dua rakaat shalat Dhuha setara dengan 360 sedekah yang diwajibkan atas setiap persendian tubuhmu.", "HR. Muslim"},
		{"Barangsiapa shalat Dhuha dua belas rakaat, niscaya Allah membangunkan baginya istana di surga.", "HR. Tirmidzi & Ibnu Majah"},
		{"Shalat Dhuha adalah shalat awwabin (orang-orang yang kembali taat), ketika anak unta kepanasan.", "HR. Muslim"},
		{"Jadikan awal harimu penuh berkah dengan Shalat Dhuha. Rezeki menanti para hamba yang bertakwa.", "Nasehat Islam"},
	},
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

// getRandomQuote memilih quote secara acak berdasarkan kategori waktu (Tidak ada perubahan)
func getRandomQuote(category string) Quote {
	quotes, ok := randomQuotes[category]
	if !ok || len(quotes) == 0 {
		return Quote{Text: "Dirikanlah salat dan tunaikanlah zakat.", Reference: "QS. Al-Baqarah: 110"} // Default
	}

	randomIndex := rand.Intn(len(quotes))
	return quotes[randomIndex]
}

// getPrayerQuoteEmbed membuat embed untuk kutipan ayat/hadits berdasarkan status waktu notifikasi (Tidak ada perubahan)
func getPrayerQuoteEmbed(label string, isBefore bool, isAdzan bool, isAfter bool) (result DiscordEmbed, categoryShalat string) {
	var category string

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

var src = rand.NewSource(time.Now().UnixNano())
var r = rand.New(src)

func GetReminderByPrefix(prefix string) string {
	reminders := map[string][]string{
		Dhuha: {
			"Waktu Dhuha telah masuk. Luangkan sedikit waktu untuk shalat Dhuha agar rezekimu hari ini diberi kelapangan dan keberkahan.",
			"Sudah masuk waktu Dhuha. Di tengah aktivitasmu, sisihkan sejenak waktu untuk shalat Dhuha sebagai bentuk syukur atas nikmat Allah.",
		},
		Subuh: {
			"Selamat pagi, semoga Allah memudahkan langkahmu. Jangan lupa shalat Subuh agar harimu diawali dengan keberkahan.",
			"Shalat Subuh telah tiba. Mulailah harimu dengan sujud yang menenangkan jiwa dan mengundang keberkahan.",
		},
		Dzuhur: {
			"Waktu Dzuhur sudah masuk. Istirahatlah sejenak dari kesibukan dan tenangkan hatimu dengan sujud kepada-Nya.",
			"Dzuhur telah tiba. Di tengah rutinitas harian, jangan lupa memenuhi panggilan Allah.",
		},
		Ashar: {
			"Ashar telah masuk. Saat yang tepat untuk menenangkan pikiran dan menyempurnakan kewajiban di tengah aktivitasmu.",
			"Sudah masuk waktu Ashar. Segarkan kembali hati dan pikiranmu dengan shalat.",
		},
		Maghrib: {
			"Maghrib telah tiba. Di waktu pergantian siang dan malam ini, jangan lewatkan kesempatan untuk mendekat kepada Allah.",
			"Waktu Maghrib telah masuk. Semoga setiap sujudmu membawa ketenangan hati.",
		},
		Isya: {
			"Isya telah masuk. Sempurnakan harimu dengan shalat sebelum beristirahat.",
			"Sudah masuk waktu Isya. Tenangkan hatimu dan mintalah perlindungan untuk malam ini.",
		},
		"sebelum": {
			"Waktu shalat sebentar lagi akan tiba. Bersiaplah untuk menunaikan ibadah dan sucikan dirimu.",
			"Adzan akan berkumandang. Mari hentikan sejenak aktivitasmu dan persiapkan diri untuk shalat.",
		},
		"setelah": {
			"Alhamdulillah. Bagi yang sudah shalat, semoga ibadahnya diterima. Bagi yang belum, mari bersegera tunaikan kewajiban. Kita doakan bersama, semoga Allah senantiasa membimbing hati kita menuju ketaatan.",
			"Waktu shalat sedang berlangsung. Mari manfaatkan kesempatan ini; bagi yang sudah, syukuri nikmatnya. Bagi yang belum, segerakan dan semoga hidayah-Nya selalu menyertai langkah kita.",
			"Istirahat sejenak untuk menunaikan shalat akan membawa ketenangan. Bagi yang telah selesai, jangan lupa berzikir. Bagi yang belum, mari segerakan. Semoga kita semua diringankan untuk memenuhi panggilan-Nya.",
		},
	}

	general := []string{
		"Sudah masuk waktu shalat. Sisihkan sebentar waktumu untuk menenangkan hati dan mengingat Allah.",
		"Ini waktu yang baik untuk menghentikan aktivitas sejenak dan menunaikan shalat.",
		"Jangan lupa menyempatkan diri untuk shalat. Semoga hatimu selalu dijaga dan dilapangkan.",
	}

	p := strings.ToLower(prefix)

	var key string

	switch {
	case strings.Contains(p, "sebelum"):
		key = "sebelum"
	case strings.Contains(p, "setelah"):
		key = "setelah"
	case strings.HasPrefix(p, Dhuha):
		key = Dhuha
	case strings.HasPrefix(p, Subuh):
		key = Subuh
	case strings.HasPrefix(p, Dzuhur):
		key = Dzuhur
	case strings.HasPrefix(p, Ashar):
		key = Ashar
	case strings.HasPrefix(p, Maghrib):
		key = Maghrib
	case strings.HasPrefix(p, Isya):
		key = Isya
	default:
		return general[r.Intn(len(general))]
	}

	if list, ok := reminders[key]; ok {
		return list[r.Intn(len(list))]
	}

	return general[r.Intn(len(general))]
}

func BuildShalatMessage(categoryShalat string) string {
	var content string

	switch categoryShalat {

	case SebelumAdzan:
		content = `
$ sistem-acara picu persiapan_shalat
----------------------------------------
 acara: persiapan_shalat
 status: pengingat_segera
 pesan: Waktu shalat hampir tiba
----------------------------------------
>> Tindakan: **Ambil wudhu dan sucikan hati.**
>> Status: Tugas yang tidak mendesak dapat ditunda sementara.
`

	case SetelahAdzan:
		content = `
[Alur Kerja] tahap("Tindak Lanjut Shalat")
[INFO] Notifikasi: Shalat sedang berlangsung
[INFO] Status: +10 menit setelah adzan
[INFO] Tindakan:
    - **Cek Status:** Pastikan kewajiban shalat sudah tertunaikan.
    - **Prioritas:** Menunaikan shalat adalah prioritas tertinggi saat ini.
    - **Koordinasi:** Mohon saling mengingatkan sesama rekan kerja.
[Alur Kerja] Selesai tahap("Tindak Lanjut Shalat")
`

	case Dhuha:
		content = `
[Alur Kerja] tahap("Peringatan Dhuha - Sunnah")
[INFO] Notifikasi: Waktu Shalat Dhuha telah tiba (Sunnah)
[INFO] Pemicu: notifikasi_dhuha()
[INFO] Tindakan:
    - **Tujuan:** Meraih keutamaan dan kelancaran rezeki harian.
    - **Durasi:** Bisa ditunaikan sampai 15 menit sebelum Dzuhur.
    - **Prioritas:** Luangkan waktu 5-10 menit.
[Alur Kerja] Selesai tahap("Peringatan Dhuha - Sunnah")
`

	default: // Digunakan untuk waktu Adzan (Subuh, Dzuhur, Ashar, Maghrib, Isya)
		content = `
[Alur Kerja] tahap("Peringatan Adzan - Wajib")
[INFO] Waktu shalat telah masuk (Fardhu)
[INFO] Pemicu: notifikasi_adzan()
[INFO] Tindakan:
    - **Prioritas:** Hentikan tugas dan penuhi panggilan-Nya.
    - **Langkah:** Bersiap, berwudhu, menuju tempat shalat.
    - **Penekanan:** Shalat di awal waktu lebih utama.
[Alur Kerja] Selesai tahap("Peringatan Adzan - Wajib")
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
		// Mengganti pesan error
		result.Error = fmt.Errorf("webhook DILEWATI: Divisi IT tidak aktif di akhir pekan (Sabtu/Minggu)")
		return result
	}

	message := GetReminderByPrefix(categoryShalat)
	message += " \n"
	message += BuildShalatMessage(categoryShalat)
	message += " \n\n\n"

	// Payload utama Discord Webhook (Field-field ini tetap 'username', 'content', dll. karena itu adalah kunci JSON Discord)
	payload := map[string]any{
		"content": message,
		"embeds":  embeds,
	}

	if strings.HasPrefix(categoryShalat, "adzan") {
		payload["tts"] = true
	}

	if categoryShalat == SebelumAdzan {
		payload["username"] = MyUsername
		payload["avatar_url"] = MyAvatar
	}

	if testingMode {
		payload["username"] = Test
		payload["avatar_url"] = TestAvaatar
	}

	body, err := json.Marshal(payload)

	if err != nil {
		result.Error = fmt.Errorf("proses konversi JSON (marshalling) gagal: %w", err)
		return result
	}

	resp, err := http.Post(target.URL, "application/json", bytes.NewBuffer(body))

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if err != nil {
		result.Error = fmt.Errorf("permintaan HTTP post gagal: %w", err)
	} else {
		if resp.StatusCode >= 300 {
			result.Error = fmt.Errorf("kode status %d", resp.StatusCode)
			fmt.Printf("Gagal mengirim webhook ke %s: %+v\n", target.Name, bodyString)
		}
	}
	return result
}

func checkAndSend(notificationTime time.Time, label string, shalatTime time.Time, offsetDuration time.Duration, target WebhookTarget) {

	// ** Gunakan Timezone dari target Webhook **
	loc, err := time.LoadLocation(target.Tz)
	if err != nil {
		fmt.Printf("Error memuat zona waktu target %s: %v\n", target.Tz, err)
		return
	}
	now := time.Now().In(loc) // now sekarang adalah waktu lokal target

	fiveMinutesAgo := now.Add(-5 * time.Minute)

	isBefore := offsetDuration < 0

	absDuration := offsetDuration
	if absDuration < 0 {
		absDuration = -absDuration
	}
	// Pastikan durasi dibulatkan ke menit terdekat sebelum diformat
	absDuration = absDuration.Truncate(time.Minute)

	// Tentukan apakah notifikasi tepat diwaktu "Adzan" (offset 0 menit atau untuk Dhuha)
	isExactlyTime := absDuration.Minutes() < 1
	// Tentukan apakah notifikasi setelah "Adzan" (offset positif)
	isAfter := !isBefore && !isExactlyTime

	// ** LOGIKA PENGIRIMAN: HANYA BERLAKU JIKA WAKTU NOTIFIKASI DI ANTARA 5 MENIT LALU HINGGA SEKARANG **
	if testingMode || (notificationTime.After(fiveMinutesAgo) && (notificationTime.Before(now) || notificationTime.Equal(now))) {

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
				mainTitle = fmt.Sprintf("🚨 ADZAN %s TELAH TIBA (Tertunda)!", label)
				// FORMAT BARU UNTUK ABSDURATION DI ISBEFORE
				formattedDurationBefore := fmt.Sprintf("%d menit", int(absDuration.Minutes()))
				mainDescription = fmt.Sprintf("📢 @here Adzan telah berkumandang (Pukul **%s**). Mohon maaf, pengingat **%s** sebelum Adzan tertunda. Segera Tunaikan Shalat.",
					shalatTime.Format("15:04"),
					formattedDurationBefore,
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
				// Ubah absDuration (time.Duration) menjadi string "X menit"
				formattedDurationAfter := fmt.Sprintf("%d menit", int(absDuration.Minutes()))
				mainDescription = fmt.Sprintf("📢 @here Shalat %s sudah dimulai (Pukul **%s**). Ini adalah notifikasi **%s** setelah Adzan. Saatnya Tunaikan Shalat.",
					label,
					shalatTime.Format("15:04"),
					formattedDurationAfter, // Menggunakan format "X menit"
				)
			}
		}

		// Hitung zona waktu untuk footer (misal: "WIB")
		zoneName, _ := now.Zone()

		mainEmbed := DiscordEmbed{
			Color:       embedColor,
			Title:       mainTitle,
			Description: mainDescription,
			Timestamp:   time.Now().Format(time.RFC3339),
			Footer: &map[string]string{
				"text": fmt.Sprintf("Lokasi: %s | Dihitung: %s %s", target.City, time.Now().In(loc).Format("02/01/2006 15:04"), zoneName),
			},
		}

		isAdzanForQuote := isExactlyTime && label != "Dhuha"
		quoteEmbed, categoryShalat := getPrayerQuoteEmbed(label, isBefore, isAdzanForQuote, isAfter)

		disclaimerEmbed := DiscordEmbed{
			Color:       10070709, // Warna Ungu (Purple/Blurple)
			Title:       "✨ Catatan",
			Description: "Pesan ini adalah pengingat otomatis. Kami menghormati keyakinan dan praktik ibadah Anda. Notifikasi ini ditujukan khusus bagi **Rekan Muslim yang menunaikan ibadah shalat**.",
		}

		// Siapkan array embed (min 3)
		embedsToSend := []DiscordEmbed{mainEmbed, quoteEmbed}

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
		result := sendWebhook(embedsToSend, target, categoryShalat)

		// Log hasil pengiriman
		logMessage := fmt.Sprintf("%s pada %s (%s). Tipe: %s. Webhook: %s",
			label,
			notificationTime.Format("15:04"),
			zoneName, // Tambahkan ZONA WAKTU di log
			func() string {
				if isBefore {
					return "SEBELUM"
				} else if isExactlyTime {
					return "SAAT MULAI"
				} else {
					return "SETELAH"
				}
			}(),
			result.Name, // Nama Webhook (akan berisi [TESTING] jika diaktifkan)
		)

		if result.Error != nil {
			if strings.Contains(result.Error.Error(), "DILEWATI: Divisi IT") {
				fmt.Printf("⏩ DILEWATI %s (AKHIR PEKAN) untuk Webhook: %s\n", label, result.Name)
			} else {
				fmt.Printf("❌ Gagal mengirim pesan untuk %s. Error: %v\n", logMessage, result.Error)
			}
		} else {
			fmt.Printf("✅ Berhasil Mengirim pesan untuk %+v. Embeds: %+v.\n", logMessage, len(embedsToSend))
		}

	} else {
		zoneName, _ := now.Zone()
		// Logika DILEWATI
		fmt.Printf(
			"⏩ DILEWATI %s (%v Waktu %s). Waktu Mulai: %s. Waktu Notif Seharusnya: %s. Waktu Saat Ini: %s %s.\n",
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
			zoneName, // Tampilkan zona waktu target di log
		)
	}
}

// ---------------------------------------------------------------------------------------
func runTesting() {
	if !testingMode {
		return
	}

	// Definisikan offsets yang akan digunakan dalam mode testing, mirip dengan mode normal
	offsetsMain := []time.Duration{}
	for m := -5; m <= 350; m += 5 {
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
	for _, target := range webhookTargets {
		if target.Name != MyChannel {
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
			Latitude:  target.Lat,
			Longitude: target.Lon,
			Timezone:  loc,
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
			{Subuh, fajr},
			{Dhuha, dhuhaStart}, // Tambahkan Dhuha di sini untuk testing
			{Dzuhur, dhuhr},
			{Ashar, asr},
			{Maghrib, maghrib},
			{Isya, isha},
		}

		// 5. Loop di setiap waktu shalat dan offset
		for _, prayer := range prayerTimes {
			currentOffsets := offsetsMain
			// Gunakan offsetsCaseSpecific untuk Dhuha
			if prayer.label == Dhuha {
				currentOffsets = offsetsCaseSpecific
			}

			for _, offset := range currentOffsets {
				notificationTime := prayer.t.Add(offset)

				if offset.Minutes() != 0 {
					continue
				}

				// jika testingMode=true dan target adalah target yang sedang diproses.
				checkAndSend(notificationTime, prayer.label, prayer.t, offset, target)
				time.Sleep(1 * time.Minute)
			}
		}
		break
	}

}

func main() {
	// Pendaftaran flag: -test
	flag.BoolVar(&testingMode, "test", false, "Aktifkan mode uji coba, kirim semua notifikasi ke testWebhookURL.")
	flag.Parse()
	runTesting()

	// Definisikan offsets yang akan digunakan dalam mode testing, mirip dengan mode normal
	offsetsMain := []time.Duration{}
	for m := -30; m <= 30; m += 5 {
		offsetsMain = append(offsetsMain, time.Duration(m)*time.Minute)
	}

	offsetsCaseSpecific := []time.Duration{}
	for m := -30; m <= 30; m += 1 {
		offsetsCaseSpecific = append(offsetsCaseSpecific, time.Duration(m)*time.Minute)
	}

	// 1. Loop di setiap Webhook Target
	for _, target := range webhookTargets {
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
			Latitude:  target.Lat,
			Longitude: target.Lon,
			Timezone:  loc,
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
			{Subuh, fajr},
			{Dhuha, dhuhaStart},
			{Dzuhur, dhuhr},
			{Ashar, asr},
			{Maghrib, maghrib},
			{Isya, isha},
		}

		// 5. Loop di setiap waktu shalat dan offset
		for _, prayer := range prayerTimes {
			currentOffsets := offsetsMain
			if prayer.label == Dhuha {
				currentOffsets = offsetsCaseSpecific
			}

			for _, offset := range currentOffsets {
				notificationTime := prayer.t.Add(offset)

				// KHUSUS Dhuha: Hanya kirim notifikasi saat offset adalah 0 menit
				if prayer.label == Dhuha && offset.Minutes() != 0 {
					continue // Lewati pengingat sebelum dan sesudah Dhuha
				}

				// ** Panggil checkAndSend dengan target Webhook saat ini **
				checkAndSend(notificationTime, prayer.label, prayer.t, offset, target)
			}
		}
	}
}
