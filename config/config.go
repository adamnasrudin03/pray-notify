package config

import "github.com/adamnasrudin03/pray-notify/models"

const (
	MyUsername = "Adam Nasrudin (@adamnasrudin)"
	Test       = "Testing"
	DivisiIT   = "Divisi IT"
	MyChannel  = "My Channel"
)

// WebhookTargets adalah daftar target webhook dengan lokasi spesifik
var WebhookTargets = []models.WebhookTarget{
	{
		Name: DivisiIT,
		URL:  "https://discord.com/api/webhooks/1443879495324008488/z-jzPBXv6vBvE45-QFrV-wwcIp8Qz9GtMw0bSQeKh1ltjHXZjcUPmjN3pMpS7Jg4yZJ0",
		Lat:  models.DefaultLat,
		Lon:  models.DefaultLon,
		Tz:   models.DefaultTz,
		City: models.DefaultCity,
	},
	{
		Name:      MyChannel,
		URL:       "https://discordapp.com/api/webhooks/1443512213842296904/EHX00Lvi8Hmqxzc1WsgJK9pyDsJzAmXWhexbqSsJnez8KSzaAHaBy_8UQIAgtRQyQEix",
		Lat:       models.DefaultLat,
		Lon:       models.DefaultLon,
		Tz:        models.DefaultTz,
		City:      models.DefaultCity,
		IsDisable: true,
	},
	{
		Name:     "Yusuf",
		URL:      "https://discord.com/api/webhooks/1445697235231838208/UaQeKj2rS9tXdoDdB-SIMna_yfd2MKlPBlyOPWfUMIQlMvn0RM9CHjAOJ_XiB6dcygQx",
		Lat:      models.DefaultLat,
		Lon:      models.DefaultLon,
		Tz:       models.DefaultTz,
		City:     models.DefaultCity,
		Mentions: "<@1211900014595407883>",
	},
}
