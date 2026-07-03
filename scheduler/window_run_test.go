package scheduler

import "testing"

func TestIsValidPrayerLabel(t *testing.T) {
	if !IsValidPrayerLabel("subuh") {
		t.Fatal("subuh should be valid")
	}
	if IsValidPrayerLabel("invalid") {
		t.Fatal("invalid should not be valid")
	}
}

func TestOffsetsForPrayer(t *testing.T) {
	if len(offsetsForPrayer("dhuha")) != 1 {
		t.Fatal("dhuha should have 1 offset")
	}
	if len(offsetsForPrayer("ashar")) != len(MainOffsets) {
		t.Fatal("ashar should use main offsets")
	}
}
