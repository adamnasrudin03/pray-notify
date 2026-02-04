package discord

import (
	"fmt"

	"github.com/adamnasrudin03/pray-notify/models"
	"github.com/adamnasrudin03/pray-notify/quotes"
)

// GetPrayerQuoteEmbed membuat embed untuk kutipan ayat/hadits berdasarkan status waktu notifikasi
func GetPrayerQuoteEmbed(label string, isBefore bool, isAdzan bool, isAfter bool) (result models.DiscordEmbed, categoryShalat string) {
	var category string

	if label == models.Dhuha {
		category = models.Dhuha
	} else if isAdzan {
		// Jika tepat waktu Adzan, gunakan kategori spesifik (e.g., "AdzanSubuh")
		category = "Adzan" + label
	} else if isBefore {
		category = models.SebelumAdzan
	} else if isAfter {
		category = models.SetelahAdzan
	} else {
		// Fallback: Untuk waktu Adzan yang tidak terdeteksi labelnya (seharusnya tidak terjadi)
		category = "Adzan" + label
	}

	selected := quotes.GetRandomQuote(category)

	return models.DiscordEmbed{
		Color:       3447003, // Warna Biru (Blue)
		Title:       "📖 Kutipan",
		Description: fmt.Sprintf("> *%s*\n\n— %s", selected.Text, selected.Reference),
	}, category
}
