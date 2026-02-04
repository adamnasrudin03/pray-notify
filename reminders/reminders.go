package reminders

import (
	"math/rand"
	"strings"
	"time"

	"github.com/adamnasrudin03/pray-notify/models"
)

var (
	src = rand.NewSource(time.Now().UnixNano())
	r   = rand.New(src)
)

var reminders = map[string][]string{
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

var general = []string{
	"Sudah masuk waktu shalat. Sisihkan sebentar waktumu untuk menenangkan hati dan mengingat Allah SWT.",
	"Ini waktu yang baik untuk menghentikan aktivitas sejenak dan menunaikan shalat.",
	"Jangan lupa menyempatkan diri untuk shalat. Semoga hatimu selalu dijaga dan dilapangkan.",
	"Ingatlah Allah SWT di sela kesibukanmu. Segera tunaikan shalat agar aktivitasmu lebih berkah.",
	"Mari raih ketenangan dengan bersujud. Waktu shalat adalah waktu terbaik untuk mengisi energi spiritual.",
}

// GetReminderByPrefix mengembalikan pesan pengingat berdasarkan prefix waktu shalat
func GetReminderByPrefix(prefix string) string {
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
