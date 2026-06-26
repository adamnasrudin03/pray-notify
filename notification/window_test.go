package notification

import (
	"testing"
	"time"
)

func TestIsInNotificationWindow(t *testing.T) {
	notif := time.Date(2026, 6, 26, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		now      time.Time
		testing  bool
		expected bool
	}{
		{"exact time", notif, false, true},
		{"30s after", notif.Add(30 * time.Second), false, true},
		{"89s after", notif.Add(89 * time.Second), false, true},
		{"90s boundary", notif.Add(90 * time.Second), false, false},
		{"too old", notif.Add(2 * time.Minute), false, false},
		{"future", notif.Add(-1 * time.Minute), false, false},
		{"testing mode bypass", notif.Add(-1 * time.Hour), true, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsInNotificationWindow(notif, tc.now, tc.testing)
			if got != tc.expected {
				t.Fatalf("got %v, want %v", got, tc.expected)
			}
		})
	}
}

func TestFilterActiveOffsets(t *testing.T) {
	loc := time.FixedZone("WIB", 7*3600)
	shalatTime := time.Date(2026, 6, 26, 12, 0, 0, 0, loc)
	now := shalatTime.Add(45 * time.Second)

	offsets := []time.Duration{
		-30 * time.Minute,
		-5 * time.Minute,
		0,
		5 * time.Minute,
	}

	active := FilterActiveOffsets(shalatTime, now, offsets, false)
	if len(active) != 1 {
		t.Fatalf("expected 1 active offset, got %d: %v", len(active), active)
	}
	if active[0] != 0 {
		t.Fatalf("expected offset 0, got %v", active[0])
	}
}

func TestIsInNotificationWindow_noDuplicateNextMinute(t *testing.T) {
	notif := time.Date(2026, 6, 26, 12, 0, 0, 0, time.UTC)

	if !IsInNotificationWindow(notif, notif.Add(30*time.Second), false) {
		t.Fatal("expected delivery within first minute")
	}
	if IsInNotificationWindow(notif, notif.Add(90*time.Second), false) {
		t.Fatal("should not deliver again on next cron cycle")
	}
}
