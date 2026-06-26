package notification

import (
	"testing"
	"time"
)

func TestIsInNotificationWindow(t *testing.T) {
	now := time.Date(2026, 6, 26, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		notifAt  time.Time
		testing  bool
		expected bool
	}{
		{"inside window", now.Add(-3 * time.Minute), false, true},
		{"exactly now", now, false, true},
		{"five minutes ago boundary", now.Add(-5 * time.Minute), false, false},
		{"too old", now.Add(-6 * time.Minute), false, false},
		{"future", now.Add(1 * time.Minute), false, false},
		{"testing mode bypass", now.Add(-1 * time.Hour), true, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsInNotificationWindow(tc.notifAt, now, tc.testing)
			if got != tc.expected {
				t.Fatalf("got %v, want %v", got, tc.expected)
			}
		})
	}
}

func TestFilterActiveOffsets(t *testing.T) {
	loc := time.FixedZone("WIB", 7*3600)
	shalatTime := time.Date(2026, 6, 26, 12, 0, 0, 0, loc)
	now := shalatTime.Add(2 * time.Minute)

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
