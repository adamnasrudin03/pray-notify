package scheduler

import "testing"

func TestMainOffsets(t *testing.T) {
	if len(MainOffsets) != 5 {
		t.Fatalf("expected 5 offsets, got %d", len(MainOffsets))
	}
	expected := []int{-10, -5, 0, 5, 10}
	for i, off := range MainOffsets {
		got := int(off.Minutes())
		if got != expected[i] {
			t.Fatalf("offset[%d]: got %d, want %d", i, got, expected[i])
		}
	}
}
