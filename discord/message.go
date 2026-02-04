package discord

import (
	"fmt"

	"github.com/adamnasrudin03/pray-notify/models"
)

// BuildShalatMessage membuat pesan shalat berdasarkan kategori
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
