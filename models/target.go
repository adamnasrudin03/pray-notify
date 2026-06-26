package models

import "time"

var workAdzanLabels = map[string]bool{
	Dzuhur:  true,
	Ashar:   true,
	Maghrib: true,
}

// ShouldSkipWeekend mengembalikan true jika target tidak aktif di akhir pekan.
func (t WebhookTarget) ShouldSkipWeekend(weekday time.Weekday) bool {
	if t.SkipWeekends != nil && *t.SkipWeekends {
		return weekday == time.Saturday || weekday == time.Sunday
	}
	return false
}

// EnableTextToSpeech mengembalikan apakah TTS Discord diaktifkan.
func (t WebhookTarget) EnableTextToSpeech() bool {
	return t.EnableTTS != nil && *t.EnableTTS
}

// ContentPrefix mengembalikan baris pembuka pesan webhook (mention @here/custom/none).
func (t WebhookTarget) ContentPrefix(label string) string {
	switch t.MentionMode {
	case "none":
		return "📢📢📢 \n"
	case "custom":
		if t.Mentions != "" {
			return "📢 " + t.Mentions + " \n"
		}
		return "📢📢📢 \n"
	default:
		if workAdzanLabels[label] {
			return "📢 @here \n"
		}
		return "📢📢📢 \n"
	}
}
