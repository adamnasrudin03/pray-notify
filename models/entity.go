package models

type WebhookTarget struct {
	Name     string `json:"Name"`
	URL      string `json:"URL"`
	Mentions string `json:"Mentions,omitempty"`
	// Data Lokasi Spesifik
	Lat  float64 `json:"Lat"`
	Lon  float64 `json:"Lon"`
	Tz   string  `json:"Tz"`   // Timezone string, misal "Asia/Jakarta"
	City string  `json:"City"` // Nama kota/lokasi untuk ditampilkan di footer
	// SkipWeekends: jika true, notifikasi tidak dikirim pada Sabtu/Minggu
	SkipWeekends *bool `json:"SkipWeekends,omitempty"`
	// MentionMode: "here" (default), "none", atau "custom" (gunakan field Mentions)
	MentionMode string `json:"MentionMode,omitempty"`
	// EnableTTS: aktifkan text-to-speech di Discord (default: false)
	EnableTTS *bool `json:"EnableTTS,omitempty"`
	// Username: override nama bot untuk semua notifikasi
	Username string `json:"Username,omitempty"`
	// UsernameBeforeAdzan: override nama bot khusus notifikasi sebelum adzan
	UsernameBeforeAdzan string `json:"UsernameBeforeAdzan,omitempty"`
	IsDisable           bool   `json:"IsDisable,omitempty"`
}

type DiscordEmbed struct {
	Color       int                `json:"color"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Timestamp   string             `json:"timestamp"`
	Footer      *map[string]string `json:"footer,omitempty"`
}

type Quote struct {
	Text      string
	Reference string
}
