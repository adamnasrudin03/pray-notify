package notification

import "time"

const notificationWindow = 5 * time.Minute

// IsInNotificationWindow memeriksa apakah waktu notifikasi berada dalam window 5 menit terakhir.
func IsInNotificationWindow(notificationTime, now time.Time, testingMode bool) bool {
	if testingMode {
		return true
	}
	fiveMinutesAgo := now.Add(-notificationWindow)
	return notificationTime.After(fiveMinutesAgo) && !notificationTime.After(now)
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
