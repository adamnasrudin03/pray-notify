package main

import (
	"math/rand"
	"strings"
	"time"
)

// Inisialisasi sumber angka acak global (lebih modern dari rand.Seed)
// Source disarankan diinisialisasi sekali per aplikasi.
var src = rand.NewSource(time.Now().UnixNano())
var r = rand.New(src)

// GetReminderByPrefix: Menentukan jenis shalat dari prefix dan mengembalikan reminder acak.
// Input contoh "AdzanDhuha" → output reminder dhuha
func GetReminderByPrefix(prefix string) string {
	// Daftar reminder per shalat
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
		// Tambahan reminder untuk kasus Sebelum/Setelah Adzan (sesuaikan dengan kebutuhan spesifik Anda)
		// Karena tidak ada data reminder baru yang diberikan, kita buat kunci baru
		"sebelum": {
			"Waktu shalat sebentar lagi akan tiba. Bersiaplah untuk menunaikan ibadah dan sucikan dirimu.",
			"Adzan akan berkumandang. Mari hentikan sejenak aktivitasmu dan persiapkan diri untuk shalat.",
		},
		"setelah": {
			"Setelah shalat, luangkan waktu sebentar untuk berzikir dan berdoa. Semoga ibadah kita diterima.",
			"Alhamdulillah sudah menunaikan shalat. Jangan lupa beristighfar dan memohon ampunan-Nya.",
		},
	}

	// Default general reminders jika prefix tidak cocok
	general := []string{
		"Sudah masuk waktu shalat. Sisihkan sebentar waktumu untuk menenangkan hati dan mengingat Allah.",
		"Ini waktu yang baik untuk menghentikan aktivitas sejenak dan menunaikan shalat.",
		"Jangan lupa menyempatkan diri untuk shalat. Semoga hatimu selalu dijaga dan dilapangkan.",
	}

	// Normalisasi prefix: semua huruf kecil
	p := strings.ToLower(prefix)

	var key string

	// Switch–case baru menggunakan logika Contains atau HasPrefix
	// Saya memilih HasPrefix karena mencocokkan kata kunci shalat dari awal input yang sudah dinormalisasi.
	switch {
	// 1. Kasus Khusus: Sebelum/Setelah Adzan
	case strings.Contains(p, "sebelum"):
		key = "sebelum"
	case strings.Contains(p, "setelah"):
		key = "setelah"
	// 2. Kasus Waktu Shalat Utama (atau saat Adzan)
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
		// Jika tidak cocok → gunakan general reminder
		return general[r.Intn(len(general))]
	}

	// Ambil reminder acak berdasarkan key
	if list, ok := reminders[key]; ok {
		return list[r.Intn(len(list))]
	}

	// Safety fallback
	return general[r.Intn(len(general))]
}
