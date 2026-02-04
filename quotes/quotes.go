package quotes

import (
	"math/rand"

	"github.com/adamnasrudin03/pray-notify/models"
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
		{"Barangsiapa shalat Shubuh berjama'ah, maka dia berada dalam jaminan Allah SWT.", "HR. Muslim"},
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

// GetRandomQuote memilih quote secara acak berdasarkan kategori waktu
func GetRandomQuote(category string) models.Quote {
	quotes, ok := randomQuotes[category]
	if !ok || len(quotes) == 0 {
		return models.Quote{Text: "Dirikanlah salat dan tunaikanlah zakat.", Reference: "QS. Al-Baqarah: 110"} // Default
	}

	randomIndex := rand.Intn(len(quotes))
	return quotes[randomIndex]
}
