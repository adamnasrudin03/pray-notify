package notification

import "time"

// maxDeliveryDelay: cron-job.org jalan tiap menit; window sempit mencegah kirim ulang notifikasi sama.
// 90 detik = 1 siklus cron + toleransi startup GitHub Actions.
const maxDeliveryDelay = 90 * time.Second

// IsInNotificationWindow memeriksa apakah waktu notifikasi siap dikirim pada run ini.
func IsInNotificationWindow(notificationTime, now time.Time, testingMode bool) bool {
	if testingMode {
		return true
	}
	if now.Before(notificationTime) {
		return false
	}
	return now.Sub(notificationTime) < maxDeliveryDelay
}

// FilterActiveOffsets memfilter offset yang waktu notifikasinya masuk dalam window aktif.
func FilterActiveOffsets(shalatTime time.Time, now time.Time, offsets []time.Duration, testingMode bool) []time.Duration {
	if testingMode {
		return offsets
	}
	var active []time.Duration
	for _, offset := range offsets {
		if IsInNotificationWindow(shalatTime.Add(offset), now, false) {
			active = append(active, offset)
		}
	}
	return active
}
