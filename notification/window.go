package notification

import "time"

// maxDeliveryDelay dipakai mode legacy Run() (polling). Mode window memakai sleep + forceSend.
const maxDeliveryDelay = 90 * time.Second

// MaxDeliveryDelay diekspor untuk scheduler window mode.
func MaxDeliveryDelay() time.Duration {
	return maxDeliveryDelay
}

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
