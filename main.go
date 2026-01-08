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

	"github.com/adamnasrudin03/pray-notify/models"
	"github.com/hablullah/go-prayer"
)

const (
	MyUsername = "Adam Nasrudin (@adamnasrudin)"
	Test       = "Testing"
	DivisiIT   = "Divisi IT"
	MyChannel  = "My Channel"
)

func isAdzanInProgress(c string) bool {
	return strings.Contains(c, "Adzan") &&
		!strings.Contains(c, "Sebelum") &&
		!strings.Contains(c, "Setelah")
}

var (
	testingMode bool
	// Koordinat Lion Parcel Kedoya, Kebon Jeruk, Jakarta Barat

	// Daftar Webhook Target dengan lokasi spesifik (tidak diubah)
	webhookTargets = []models.WebhookTarget{
		{
			Name: DivisiIT,
			URL:  "https://discord.com/api/webhooks/1443879495324008488/z-jzPBXv6vBvE45-QFrV-wwcIp8Qz9GtMw0bSQeKh1ltjHXZjcUPmjN3pMpS7Jg4yZJ0",
			Lat:  models.DefaultLat,
			Lon:  models.DefaultLon,
			Tz:   models.DefaultTz,
			City: models.DefaultCity,
		},
		{
			Name: MyChannel,
			URL:  "https://discordapp.com/api/webhooks/1443512213842296904/EHX00Lvi8Hmqxzc1WsgJK9pyDsJzAmXWhexbqSsJnez8KSzaAHaBy_8UQIAgtRQyQEix",
			Lat:  models.DefaultLat,
			Lon:  models.DefaultLon,
			Tz:   models.DefaultTz,
			City: models.DefaultCity,
		},
		{
			Name:     "Yusuf",
			URL:      "https://discord.com/api/webhooks/1445697235231838208/UaQeKj2rS9tXdoDdB-SIMna_yfd2MKlPBlyOPWfUMIQlMvn0RM9CHjAOJ_XiB6dcygQx",
			Lat:      models.DefaultLat,
			Lon:      models.DefaultLon,
			Tz:       models.DefaultTz,
			City:     models.DefaultCity,
			Mentions: "<@1211900014595407883>",
		},
	}
)

var randomQuotes = map[string][]models.Quote{
	models.Dhuha: {
		{"Dua rakaat shalat Dhuha setara dengan 360 sedekah yang diwajibkan atas setiap persendian tubuhmu.", "HR. Muslim"},
		{"Barangsiapa shalat Dhuha dua belas rakaat, niscaya Allah SWT membangunkan baginya istana di surga.", "HR. Tirmidzi & Ibnu Majah"},
		{"Shalat Dhuha adalah shalat awwabin (orang-orang yang kembali taat), ketika anak unta kepanasan.", "HR. Muslim"},
		{"Jadikan awal harimu penuh berkah dengan Shalat Dhuha. Rezeki menanti para hamba yang bertakwa.", "Nasehat Islam"},
	},
	models.Jumat: {
		{"Wahai orang-orang yang beriman! Apabila telah diseru untuk melaksanakan salat pada hari Jumat, maka segeralah kamu mengingat Allah dan tinggalkanlah jual beli.", "QS. Al-Jumu'ah: 9"},
		{"Sebaik-baik hari yang matahari terbit padanya adalah hari Jumat; padanya Adam diciptakan, dimasukkan ke surga dan dikeluarkan darinya.", "HR. Muslim"},
		{"Barangsiapa mandi pada hari Jumat, memakai pakaian terbaik, dan memakai wangi-wangian, lalu pergi shalat Jumat, maka diampuni dosa-dosanya antara dua Jumat.", "HR. Abu Daud"},
	},
	models.SebelumAdzan: {
		{"Jagalah shalat lima waktu, terutama shalat yang pertengahan (Ashar) dan berdirilah untuk Allah SWT (dalam shalatmu) dengan khusyuk.", "QS. Al-Baqarah: 238"},
		{"Shalat di awal waktu itu adalah amalan yang paling dicintai oleh Allah SWT.", "HR. Bukhari dan Muslim"},
		{"Bersegeralah mendirikan shalat sebelum sibuk dengan urusan dunia.", "Nasehat Islam"},
		{"Sesungguhnya shalat pada waktunya adalah kewajiban yang telah ditentukan atas orang-orang yang beriman.", "QS. An-Nisa: 103"},
	},
	models.AdzanSubuh: {
		{"Sesungguhnya shalat Subuh itu disaksikan (oleh malaikat-malaikat malam dan malaikat-malaikat siang).", "QS. Al-Isra: 78"},
		{"Barangsiapa shalat Shubuh berjama’ah, maka dia berada dalam jaminan Allah SWT.", "HR. Muslim"},
	},
	models.AdzanDzuhur: {
		{"Jagalah (peliharalah) segala shalatmu dan (peliharalah) shalat wustha. Berdirilah untuk Allah SWT(dalam shalatmu) dengan khusyuk.", "QS. Al-Baqarah: 238"},
		{"Apabila kamu telah menyelesaikan shalat, ingatlah Allah SWT ketika kamu berdiri, ketika kamu duduk dan ketika kamu berbaring.", "QS. An-Nisa: 103"},
	},
	models.AdzanAshar: {
		{"Barangsiapa meninggalkan shalat Ashar, maka seolah-olah dia telah kehilangan keluarga dan hartanya.", "HR. Bukhari dan Muslim"},
		{"Shalat wustha adalah shalat Ashar. Siapa yang menjaganya, ia telah mendapat keberuntungan besar.", "HR. Tirmidzi"},
	},
	models.AdzanMaghrib: {
		{"Jika datang malam dari arah sini (timur) dan pergi siang dari arah sini (barat), serta terbenam matahari, maka sungguh orang yang puasa telah berbuka (waktunya Maghrib).", "HR. Bukhari dan Muslim"},
		{"Allah SWT tidak akan menyia-nyiakan pahala orang-orang yang berbuat baik.", "QS. Ali Imran: 171"},
	},
	models.AdzanIsya: {
		{"Barangsiapa yang shalat Isya' berjama'ah, maka seolah-olah ia telah shalat setengah malam.", "HR. Muslim"},
		{"Shalat yang paling berat bagi orang munafik adalah shalat Isya' dan shalat Subuh.", "HR. Bukhari dan Muslim"},
	},
	models.SetelahAdzan: {
		{"Dirikanlah salat dan tunaikanlah zakat. Kebaikan apa pun yang kamu kerjakan untuk dirimu, kamu akan mendapatkannya (pahala) di sisi Allah SWT.", "QS. Al-Baqarah: 110"},
		{"Sesungguhnya orang-orang yang beriman itu hanyalah orang-orang yang apabila disebutkan Allah SWT gemetarlah hati mereka, dan apabila dibacakan kepada mereka ayat-ayat-Nya, bertambahlah iman mereka.", "QS. Al-Anfal: 2 (Mengingatkan kewajiban setelah mendengar Adzan)"},
		{"Dan mintalah pertolongan (kepada Allah SWT) dengan sabar dan salat. Dan sesungguhnya yang demikian itu sungguh berat, kecuali bagi orang-orang yang khusyuk.", "QS. Al-Baqarah: 45"},
		{"Peliharalah semua salat(mu), dan (peliharalah) salat wusṭā. Berdirilah karena Allah SWT (dalam salatmu) dengan khusyuk.", "QS. Al-Baqarah: 238 (Pengulangan penekanan kewajiban)"},
	},
}

// WebhookResult menyimpan hasil pengiriman ke satu webhook
type WebhookResult struct {
	Name  string
	Error error
}

// getRandomQuote memilih quote secara acak berdasarkan kategori waktu (Tidak ada perubahan)
func getRandomQuote(category string) models.Quote {
	quotes, ok := randomQuotes[category]
	if !ok || len(quotes) == 0 {
		return models.Quote{Text: "Dirikanlah salat dan tunaikanlah zakat.", Reference: "QS. Al-Baqarah: 110"} // Default
	}

	randomIndex := rand.Intn(len(quotes))
	return quotes[randomIndex]
}

// getPrayerQuoteEmbed membuat embed untuk kutipan ayat/hadits berdasarkan status waktu notifikasi (Tidak ada perubahan)
func getPrayerQuoteEmbed(label string, isBefore bool, isAdzan bool, isAfter bool) (result models.DiscordEmbed, categoryShalat string) {
	var category string

	if label == models.Dhuha {
		category = models.Dhuha
	} else if isAdzan {
		// Jika tepat waktu Adzan, gunakan kategori spesifik (e.g., "AdzanSubuh")
		category = "Adzan" + label
	} else if isBefore {
		category = models.SebelumAdzan
	} else if isAfter {
		category = models.SetelahAdzan
	} else {
		// Fallback: Untuk waktu Adzan yang tidak terdeteksi labelnya (seharusnya tidak terjadi)
		category = "Adzan" + label
	}

	selected := getRandomQuote(category)

	return models.DiscordEmbed{
		Color:       3447003, // Warna Biru (Blue)
		Title:       "📖 Kutipan",
		Description: fmt.Sprintf("> *%s*\n\n— %s", selected.Text, selected.Reference),
	}, category
}

var (
	src = rand.NewSource(time.Now().UnixNano())
	r   = rand.New(src)
)

func GetReminderByPrefix(prefix string) string {
	reminders := map[string][]string{
		models.Dhuha: {
			"Waktu Dhuha telah masuk. Luangkan sedikit waktu untuk shalat Dhuha agar rezekimu hari ini diberi kelapangan dan keberkahan.",
			"Sudah masuk waktu Dhuha. Di tengah aktivitasmu, sisihkan sejenak waktu untuk shalat Dhuha sebagai bentuk syukur atas nikmat Allah SWT.",
			"Mari awali kesibukan dengan shalat Dhuha. Semoga Allah SWT memberkahi setiap usaha dan langkahmu hari ini.",
			"Waktu Dhuha adalah waktu yang berkah. Sempatkan dua rakaat untuk ketenangan jiwa dan kelancaran rezeki.",
			"Shalat Dhuha telah tiba. Jadikan momen ini untuk memohon petunjuk dan perlindungan-Nya di sepanjang hari ini.",
		},
		models.Subuh: {
			"Selamat pagi, semoga Allah SWT memudahkan langkahmu. Jangan lupa shalat Subuh agar harimu diawali dengan keberkahan.",
			"Shalat Subuh telah tiba. Mulailah harimu dengan sujud yang menenangkan jiwa dan mengundang keberkahan.",
			"Awali pagimu dengan ketaatan. Shalat Subuh adalah kunci pembuka pintu rezeki dan pertolongan Allah SWT.",
			"Udara pagi yang sejuk dan sujud Subuh yang khusyuk adalah kombinasi terbaik untuk memulai hari dengan energi positif.",
			"Subuh telah memanggil. Mari raih keutamaan shalat di awal waktu untuk hari yang lebih bermakna.",
			"Adzan Subuh sedang berkumandang. Mari segera bersiap menuju masjid atau tunaikan shalat di awal waktu.",
			"Dengarkan panggilan-Nya di waktu fajar. Mari raih berkah Subuh dengan shalat berjamaah.",
			"Panggilan fajar telah tiba. Hentikan tidur sejenak, shalat itu lebih baik daripada tidur.",
			"Adzan Subuh adalah awal kemenangan hari. Mari jawab panggilannya dengan penuh semangat.",
			"Keberkahan hari dimulai dari respon kita terhadap Adzan Subuh. Mari segera ke masjid.",
		},
		models.Dzuhur: {
			"Waktu Dzuhur sudah masuk. Istirahatlah sejenak dari kesibukan dan tenangkan hatimu dengan sujud kepada-Nya.",
			"Dzuhur telah tiba. Di tengah rutinitas harian, jangan lupa memenuhi panggilan Allah SWT.",
			"Segarkan kembali semangatmu dengan shalat Dzuhur. Biarkan hati sejenak beristirahat dalam komunikasi dengan Sang Pencipta.",
			"Panggilan Dzuhur adalah saat yang tepat untuk mengevaluasi niat dan memohon kekuatan untuk melanjutkan aktivitas.",
			"Mari sejenak hentikan pekerjaan dan tunaikan shalat Dzuhur. Keberkahan waktu lebih utama dari sekedar mengejar target.",
			"Adzan Dzuhur telah bergema. Mari hentikan segala aktivitas duniawi dan penuhi panggilan-Nya.",
			"Waktu istirahat terbaik adalah saat Adzan Dzuhur berkumandang. Mari tunaikan shalat tepat waktu.",
			"Dengarkan suara Adzan, itu adalah undangan langsung dari Allah SWT. Mari segera bersiap.",
			"Panggilan Dzuhur sudah tiba. Mari sucikan diri dan bersujud bersama rekan-rekan.",
			"Adzan Dzuhur mengingatkan kita bahwa ada yang lebih utama dari pekerjaan. Ayo shalat.",
		},
		models.Ashar: {
			"Ashar telah masuk. Saat yang tepat untuk menenangkan pikiran dan menyempurnakan kewajiban di tengah aktivitasmu.",
			"Sudah masuk waktu Ashar. Segarkan kembali hati dan pikiranmu dengan shalat.",
			"Jangan biarkan kesibukan sore menghalangimu dari shalat Ashar. Mari tunaikan kewajiban dengan penuh keikhlasan.",
			"Waktu Ashar adalah waktu yang mulia. Sempatkan diri untuk bersujud sebelum hari berganti menuju senja.",
			"Mari selesaikan sisa hari ini dengan keberkahan shalat Ashar yang tepat waktu.",
			"Adzan Ashar sedang berkumandang. Mari selesaikan tugas sejenak dan jawab panggilan suci ini.",
			"Jangan abaikan Adzan Ashar di tengah kesibukan soremu. Mari raih keberkahan sore dengan shalat.",
			"Panggilan Ashar telah tiba. Mari segarkan wudhu and bersiap menghadap-Nya.",
			"Adzan Ashar adalah pengingat untuk tetap terhubung dengan Sang Pencipta di sisa hari.",
			"Mari jawab seruan Adzan Ashar dengan segera menunaikan kewajiban.",
		},
		models.Maghrib: {
			"Maghrib telah tiba. Di waktu pergantian siang dan malam ini, jangan lewatkan kesempatan untuk mendekat kepada Allah SWT.",
			"Waktu Maghrib telah masuk. Semoga setiap sujudmu membawa ketenangan hati.",
			"Syukuri akhir harimu dengan shalat Maghrib. Biarkan hati tenang dalam naungan rahmat-Nya.",
			"Maghrib adalah momen refleksi. Mari tunaikan kewajiban sebelum beristirahat bersama keluarga.",
			"Waktu Maghrib sangat singkat. Ayo segera bersiap dan tunaikan shalat di awal waktu.",
			"Adzan Maghrib telah berkumandang. Mari segera berbuka bagi yang berpuasa dan segerakan shalat.",
			"Waktu pergantian hari telah tiba dengan Adzan Maghrib. Ayo bersujud syukur atas nikmat hari ini.",
			"Dengarkan Adzan Maghrib, mari segera tunaikan shalat di awal waktu yang singkat ini.",
			"Panggilan Maghrib adalah undangan untuk ketenangan malam. Ayo segera berwudhu.",
			"Mari jemput keberkahan malam hari dengan menjawab Adzan Maghrib tepat waktu.",
		},
		models.Isya: {
			"Isya telah masuk. Sempurnakan harimu dengan shalat sebelum beristirahat.",
			"Sudah masuk waktu Isya. Tenangkan hatimu dan mintalah perlindungan untuk malam ini.",
			"Tutup harimu dengan indah melalui shalat Isya. Semoga istirahatmu bernilai ibadah.",
			"Shalat Isya adalah penutup aktivitas yang menenangkan. Mari bersujud dan memohon ampunan atas kekilafan hari ini.",
			"Waktu Isya telah tiba. Luangkan waktu untuk menghadap Sang Khalik sebelum memejamkan mata.",
			"Adzan Isya sedang berkumandang. Mari tunaikan kewajiban terakhir hari ini dengan penuh khusyuk.",
			"Panggilan Isya telah tiba. Mari tutup hari ini dengan sujud yang sempurna.",
			"Dengarkan Adzan Isya, mari segera bersiap untuk shalat berjamaah di masjid.",
			"Waktu Isya telah masuk. Ayo jawab panggilannya sebelum kita beristirahat malam.",
			"Adzan Isya mengajak kita untuk melepaskan segala beban dunia sebelum tidur dalam doa.",
		},
		models.Jumat: {
			"Hari Jumat adalah sayyidul ayyam (penghulu hari). Mari persiapkan diri untuk Shalat Jumat berjamaah.",
			"Jumatan sebentar lagi. Bagi kaum laki-laki, segera persiapkan diri dan berangkat ke masjid lebih awal.",
			"Tinggalkan sejenak urusan duniamu, penuhi panggilan-Nya di hari Jumat yang mulia ini.",
			"Jangan lupa sunnah hari Jumat: Mandi, potong kuku, wangi-wangian, dan datang awal ke masjid.",
			"Shalat Jumat adalah kewajiban. Mari melangkah ke masjid dengan hati yang ikhlas mengharap ridha-Nya.",
		},
		"sebelum": {
			"Waktu shalat sebentar lagi akan tiba. Bersiaplah untuk menunaikan ibadah dan sucikan dirimu.",
			"Adzan akan berkumandang. Mari hentikan sejenak aktivitasmu dan persiapkan diri untuk shalat.",
			"Beberapa menit menuju waktu shalat. Mari mulai merapikan pekerjaan dan bersiap berwudhu.",
			"Persiapkan dirimu, waktu shalat hampir tiba. Semoga ibadah kita hari ini diterima Allah SWT.",
			"Mari bersiap menyambut panggilan Allah SWT. Segera bersihkan diri dan siapkan hati.",
			"Waktu shalat segera tiba. Mari bersiap, berwudhu, dan menenangkan hati sejenak dari aktivitas dunia.",
			"Panggilan Allah SWT hampir berkumandang. Mari hentikan sejenak urusan pekerjaan untuk menyambut waktu shalat.",
			"Hampir waktu shalat. Mari manfaatkan beberapa menit ini untuk penyucian diri dan persiapan ibadah.",
			"Segera bersiap, adzan akan segera memanggil. Mari raih keutamaan shalat di awal waktu.",
			"Waktu ibadah sudah di depan mata. Mari bergegas mempersiapkan diri demi kekhusyukan shalat.",
		},
		"setelah": {
			"Alhamdulillah. Bagi yang sudah shalat, semoga ibadahnya diterima. Bagi yang belum, mari bersegera tunaikan kewajiban. Kita doakan bersama, semoga Allah senantiasa membimbing hati kita menuju ketaatan.",
			"Waktu shalat sedang berlangsung. Mari manfaatkan kesempatan ini; bagi yang sudah, syukuri nikmatnya. Bagi yang belum, segerakan dan semoga hidayah-Nya selalu menyertai langkah kita.",
			"Istirahat sejenak untuk menunaikan shalat akan membawa ketenangan. Bagi yang telah selesai, jangan lupa berzikir. Bagi yang belum, mari segerakan. Semoga kita semua diringankan untuk memenuhi panggilan-Nya.",
			"Kewajiban shalat adalah prioritas. Manfaatkan waktu yang tersisa untuk segera menghadap Sang Khalik.",
			"Jangan tunda lagi, waktu shalat terus berjalan. Semoga Allah SWT memberikan kelapangan hati untuk menyegerakan ketaatan.",
			"Shalat sedang berlangsung. Mari manfaatkan waktu ini untuk bermunajat dan memohon bimbingan-Nya.",
			"Waktu shalat masih tersedia. Bagi yang belum, mari manfaatkan kesempatan ini untuk bersujud.",
			"Semoga setiap sujud kita menjadi penggugur dosa dan pengangkat derajat di sisi Allah SWT.",
			"Jangan biarkan waktu shalat berlalu begitu saja. Mari segerakan bagi yang masih memiliki keluangan.",
			"Ketenangan sejati ada dalam sujud. Mari tunaikan kewajiban selagi waktu masih ada.",
		},
	}

	general := []string{
		"Sudah masuk waktu shalat. Sisihkan sebentar waktumu untuk menenangkan hati dan mengingat Allah SWT.",
		"Ini waktu yang baik untuk menghentikan aktivitas sejenak dan menunaikan shalat.",
		"Jangan lupa menyempatkan diri untuk shalat. Semoga hatimu selalu dijaga dan dilapangkan.",
		"Ingatlah Allah SWT di sela kesibukanmu. Segera tunaikan shalat agar aktivitasmu lebih berkah.",
		"Mari raih ketenangan dengan bersujud. Waktu shalat adalah waktu terbaik untuk mengisi energi spiritual.",
	}

	p := strings.ToLower(prefix)

	var key string

	switch {
	case strings.Contains(p, "sebelum"):
		key = "sebelum"
	case strings.Contains(p, "setelah"):
		key = "setelah"
	case strings.HasPrefix(p, models.Dhuha):
		key = models.Dhuha
	case strings.HasPrefix(p, models.Subuh):
		key = models.Subuh
	case strings.HasPrefix(p, models.Dzuhur):
		key = models.Dzuhur
	case strings.HasPrefix(p, models.Ashar):
		key = models.Ashar
	case strings.HasPrefix(p, models.Maghrib):
		key = models.Maghrib
	case strings.HasPrefix(p, models.Isya):
		key = models.Isya
	case strings.HasPrefix(p, models.Jumat):
		key = models.Jumat
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

	case models.SebelumAdzan:
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

	case models.SetelahAdzan:
		content = `
[Alur Kerja] tahap("Tindak Lanjut Shalat")
[INFO] Notifikasi: Shalat sedang berlangsung
[INFO] Status: Pengingat tindak lanjut setelah Adzan berkumandang
[INFO] Tindakan:
    - **Cek Status:** Pastikan kewajiban shalat sudah tertunaikan.
    - **Prioritas:** Menunaikan shalat adalah prioritas tertinggi saat ini.
    - **Koordinasi:** Mohon saling mengingatkan sesama rekan kerja.
[Alur Kerja] Selesai tahap("Tindak Lanjut Shalat")
`

	case models.Dhuha:
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

	case models.Jumat:
		content = `
[Alur Kerja] tahap("Shalat Jumat - Wajib")
[INFO] Waktu Shalat Jumat telah tiba
[INFO] Pemicu: notifikasi_jumatan()
[INFO] Tindakan:
    - **Prioritas Utama:** Segera tinggalkan pekerjaan.
    - **Langkah:** Menuju masjid lebih awal untuk mendengarkan khutbah.
    - **Peringatan:** Jual beli/pekerjaan haram dilakukan saat adzan Jumat berkumandang.
[Alur Kerja] Selesai tahap("Shalat Jumat - Wajib")
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
func sendWebhook(label string, embeds []models.DiscordEmbed, target models.WebhookTarget, categoryShalat string) WebhookResult {
	result := WebhookResult{Name: target.Name}
	now := time.Now()

	// 🛑 VALIDASI 1: Cek Hari Libur/Weekend untuk Webhook "Divisi IT"
	if target.Name == DivisiIT && (now.Weekday() == time.Saturday || now.Weekday() == time.Sunday) {
		// Mengganti pesan error
		result.Error = fmt.Errorf("webhook DILEWATI: Divisi IT tidak aktif di akhir pekan (Sabtu/Minggu)")
		return result
	}
	isAdzanWork := map[string]bool{
		models.Dzuhur:  true,
		models.Ashar:   true,
		models.Maghrib: true,
	}

	message := "📢 @here \n"
	if !isAdzanWork[label] {
		message = "📢📢📢 \n"
	}
	message += GetReminderByPrefix(categoryShalat)
	message += " \n"
	message += BuildShalatMessage(categoryShalat)
	message += " \n\n\n"

	mentions := target.Mentions
	if target.Name == DivisiIT && isAdzanWork[label] {
		randomMembers := models.GetRandomMemberDivIT(7)
		var m []string
		for _, member := range randomMembers {
			m = append(m, member.UserID)
		}
		mentions = strings.Join(m, " ")
	}

	if mentions != "" {
		message += "Cc: " + mentions
	}

	// Payload utama Discord Webhook (Field-field ini tetap 'username', 'content', dll. karena itu adalah kunci JSON Discord)
	payload := map[string]any{
		"content": message,
		"embeds":  embeds,
		"tts":     true,
	}

	payload["avatar_url"] = models.RandomAvatarURL[rand.Intn(len(models.RandomAvatarURL))]
	payload["username"] = "Reminder Adzan " + label

	if target.Name == DivisiIT {
		temp := models.GetRandomMemberDivIT(1)
		for _, v := range temp {
			payload["username"] = v.Name
			break
		}
	}
	if target.Name != DivisiIT && categoryShalat == models.SebelumAdzan {
		payload["username"] = MyUsername
	}

	if testingMode {
		payload["username"] = Test
		payload["avatar_url"] = models.TestAvaatar
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

func checkAndSend(notificationTime time.Time, label string, shalatTime time.Time, offsetDuration time.Duration, target models.WebhookTarget) {
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
		quoteEmbed, categoryShalat := getPrayerQuoteEmbed(label, isBefore, isAdzanForQuote, isAfter)

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

		result := sendWebhook(label, embedsToSend, target, categoryShalat)

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

// ---------------------------------------------------------------------------------------
func runTesting() {
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
		for _, prayer := range prayerTimes {
			currentOffsets := offsetsMain
			// Gunakan offsetsCaseSpecific untuk Dhuha
			if prayer.label == models.Dhuha {
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
		// 🛑 VALIDASI 1: Cek Hari Libur/Weekend untuk Webhook "Divisi IT"
		if target.Name == DivisiIT && (today.Weekday() == time.Saturday || today.Weekday() == time.Sunday) {
			fmt.Println("webhook DILEWATI: Divisi IT tidak aktif di akhir pekan (Sabtu/Minggu)")
			continue
		}

		// 3. Hitung Jadwal Shalat untuk Target ini
		cfg := prayer.Config{
			Latitude:           target.Lat,
			Longitude:          target.Lon,
			Timezone:           loc,
			TwilightConvention: prayer.Kemenag(),
			AsrConvention:      prayer.Shafii,
			PreciseToSeconds:   true,
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
			{models.Dhuha, dhuhaStart},
			{models.Dzuhur, dhuhr},
			{models.Ashar, asr},
			{models.Maghrib, maghrib},
			{models.Isya, isha},
		}

		// 5. Loop di setiap waktu shalat dan offset
		for _, prayer := range prayerTimes {
			currentOffsets := offsetsMain
			if prayer.label == models.Dhuha {
				currentOffsets = offsetsCaseSpecific
			}

			for _, offset := range currentOffsets {
				notificationTime := prayer.t.Add(offset)

				// KHUSUS Dhuha: Hanya kirim notifikasi saat offset adalah 0 menit
				if prayer.label == models.Dhuha && offset.Minutes() != 0 {
					continue // Lewati pengingat sebelum dan sesudah Dhuha
				}

				// ** Panggil checkAndSend dengan target Webhook saat ini **
				checkAndSend(notificationTime, prayer.label, prayer.t, offset, target)
			}
		}
	}
}
