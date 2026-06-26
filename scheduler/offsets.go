package scheduler

import "time"

// MainOffsets: notifikasi -10, -5, 0, +5, +10 menit dari waktu shalat.
var MainOffsets = []time.Duration{
	-10 * time.Minute,
	-5 * time.Minute,
	0,
	5 * time.Minute,
	10 * time.Minute,
}
