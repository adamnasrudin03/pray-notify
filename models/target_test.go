package models

import (
	"testing"
	"time"
)

func boolPtr(v bool) *bool { return &v }

func TestShouldSkipWeekend(t *testing.T) {
	skip := true
	target := WebhookTarget{SkipWeekends: &skip}

	if !target.ShouldSkipWeekend(time.Saturday) {
		t.Fatal("expected Saturday to be skipped")
	}
	if !target.ShouldSkipWeekend(time.Sunday) {
		t.Fatal("expected Sunday to be skipped")
	}
	if target.ShouldSkipWeekend(time.Monday) {
		t.Fatal("expected Monday not to be skipped")
	}

	noSkip := WebhookTarget{}
	if noSkip.ShouldSkipWeekend(time.Saturday) {
		t.Fatal("expected no skip when SkipWeekends is nil")
	}
}

func TestContentPrefix(t *testing.T) {
	target := WebhookTarget{MentionMode: "here"}

	if prefix := target.ContentPrefix(Dzuhur); prefix != "📢 @here \n" {
		t.Fatalf("dzuhur prefix: got %q", prefix)
	}
	if prefix := target.ContentPrefix(Subuh); prefix != "📢📢📢 \n" {
		t.Fatalf("subuh prefix: got %q", prefix)
	}

	none := WebhookTarget{MentionMode: "none"}
	if prefix := none.ContentPrefix(Dzuhur); prefix != "📢📢📢 \n" {
		t.Fatalf("none mode: got %q", prefix)
	}

	custom := WebhookTarget{MentionMode: "custom", Mentions: "<@&123>"}
	if prefix := custom.ContentPrefix(Subuh); prefix != "📢 <@&123> \n" {
		t.Fatalf("custom mode: got %q", prefix)
	}
}

func TestEnableTextToSpeech(t *testing.T) {
	enabled := WebhookTarget{EnableTTS: boolPtr(true)}
	if !enabled.EnableTextToSpeech() {
		t.Fatal("expected TTS enabled")
	}

	disabled := WebhookTarget{}
	if disabled.EnableTextToSpeech() {
		t.Fatal("expected TTS disabled by default")
	}
}
