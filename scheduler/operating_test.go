package scheduler

import (
	"testing"
	"time"
)

func TestIsWithinOperatingWindow(t *testing.T) {
	loc := time.FixedZone("WIB", 7*3600)

	tests := []struct {
		name     string
		at       time.Time
		expected bool
	}{
		{"before window", time.Date(2026, 6, 26, 3, 0, 0, 0, loc), false},
		{"window start", time.Date(2026, 6, 26, 3, 30, 0, 0, loc), true},
		{"midday", time.Date(2026, 6, 26, 12, 0, 0, 0, loc), true},
		{"window end", time.Date(2026, 6, 26, 19, 59, 0, 0, loc), true},
		{"after window", time.Date(2026, 6, 26, 20, 0, 0, 0, loc), false},
		{"late isya offset", time.Date(2026, 6, 26, 19, 45, 0, 0, loc), true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsWithinOperatingWindow(tc.at, loc)
			if got != tc.expected {
				t.Fatalf("got %v, want %v", got, tc.expected)
			}
		})
	}
}
