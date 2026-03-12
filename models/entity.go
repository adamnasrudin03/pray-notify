package models

type WebhookTarget struct {
	Name     string
	URL      string
	Mentions string
	// Data Lokasi Spesifik
	Lat       float64
	Lon       float64
	Tz        string // Timezone string, misal "Asia/Jakarta"
	City      string // Nama kota/lokasi untuk ditampilkan di footer
	IsDisable bool
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

type DiscordMember struct {
	UserID string
	Name   string
}
