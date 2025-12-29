package models

import (
	"slices"

	"github.com/adamnasrudin03/pray-notify/helpers"
)

// Manager
var membersManagerDivIT = []DiscordMember{
	{UserID: "<@1211554519981236248>", Name: "Abi"},
	{UserID: "<@1210162463782215713>", Name: "Ade Muhammad"},
}

// SRE
var membersSreDivIT = []DiscordMember{
	{UserID: "<@1341004183192473624>", Name: "Abdul Zailani"},
	{UserID: "<@820218039106273290>", Name: "Fauzi Musyafa"},
	{UserID: "<@1313051235225702421>", Name: "Kresno Hadinata"},
}

// FORCE
var membersForceDivIT = []DiscordMember{
	{UserID: "<@1209829712290123806>", Name: "Tedja Adhi Putra"},
	{UserID: "<@873083047501787156>", Name: "Adam"},
	{UserID: "<@1403317129163440241>", Name: "Achmad H"},
	{UserID: "<@1450763384424370227>", Name: "Fuji Nugraha"},
	// {UserID: "<@571235433271394304>", Name: "Andrean"},
}

// SDET
var membersSdetDivIT = []DiscordMember{
	{UserID: "<@812158172608593940>", Name: "Arifin"},
	{UserID: "<@235784508291547136>", Name: "Hanif Afianto DN"},
	{UserID: "<@1389437159496220694>", Name: "Adit Triyadi"},
	{UserID: "<@369433975304880131>", Name: "Septian Ega"},
	{UserID: "<@1290943031653961731>", Name: "Andri Sukma"},
	{UserID: "<@1301798033805283341>", Name: "Ayu Anriyani"},
	{UserID: "<@395822450983567379>", Name: "Hendra R Gunawan"},
}

// PE
var membersPeDivIT = []DiscordMember{
	{UserID: "<@1210574321886429246>", Name: "Bian Syahid Kusumah"},
	{UserID: "<@1210167304533581835>", Name: "Ahmad Safar"},
	{UserID: "<@1385455877359665234>", Name: "Anggie Dimasta"},
	{UserID: "<@1325673268111343631>", Name: "Ari Firdaus"},
	// {UserID: "<@1290595061318553721>", Name: "Bambang Harianto Sianturi"},
	{UserID: "<@1376420636796260382>", Name: "Bayu Muhammad Iqbal"},
	{UserID: "<@1211558916857798727>", Name: "Hasan Ikhtiarta"},
	{UserID: "<@1313024881989779526>", Name: "Ihsan Fajar Ramadhan"},
	{UserID: "<@1199235749485224026>", Name: "Rohmat"},
	{UserID: "<@978122157454405672>", Name: "Jafar Rohmadi"},
}

func GetDiscordUserIDManagerDivIT() []string {
	var userID []string
	for _, member := range membersManagerDivIT {
		userID = append(userID, member.UserID)
	}
	return userID
}

func GetDiscordUserIDSreDivIT() []string {
	var userID []string
	for _, member := range membersSreDivIT {
		userID = append(userID, member.UserID)
	}
	return userID
}

func GetDiscordUserIDForceDivIT() []string {
	var userID []string
	for _, member := range membersForceDivIT {
		userID = append(userID, member.UserID)
	}
	return userID
}

func GetDiscordUserIDSdetDivIT() []string {
	var userID []string
	for _, member := range membersSdetDivIT {
		userID = append(userID, member.UserID)
	}
	return userID
}

func GetDiscordUserIDPeDivIT() []string {
	var userID []string
	for _, member := range membersPeDivIT {
		userID = append(userID, member.UserID)
	}
	return userID
}

func GetDiscordUserIDDivIT() []string {
	var userID []string
	userID = append(userID, GetDiscordUserIDManagerDivIT()...)
	userID = append(userID, GetDiscordUserIDSreDivIT()...)
	userID = append(userID, GetDiscordUserIDForceDivIT()...)
	userID = append(userID, GetDiscordUserIDSdetDivIT()...)
	userID = append(userID, GetDiscordUserIDPeDivIT()...)
	return userID
}

func GetRandomMemberDivIT(total int) string {
	if total < 5 {
		return helpers.GetRandomString(total, GetDiscordUserIDDivIT())
	}
	result := ""

	// Collect already selected user IDs
	var alreadySelected []string

	selectedManager := helpers.GetRandomString(1, GetDiscordUserIDManagerDivIT())
	result += selectedManager + " "
	alreadySelected = append(alreadySelected, selectedManager)

	selectedSre := helpers.GetRandomString(1, GetDiscordUserIDSreDivIT())
	result += selectedSre + " "
	alreadySelected = append(alreadySelected, selectedSre)

	selectedSdet := helpers.GetRandomString(1, GetDiscordUserIDSdetDivIT())
	result += selectedSdet + " "
	alreadySelected = append(alreadySelected, selectedSdet)

	selectedPe := helpers.GetRandomString(1, GetDiscordUserIDPeDivIT())
	result += selectedPe + " "
	alreadySelected = append(alreadySelected, selectedPe)

	selectedForce := helpers.GetRandomString(1, GetDiscordUserIDForceDivIT())
	result += selectedForce + " "
	alreadySelected = append(alreadySelected, selectedForce)

	// Filter out already selected IDs from the full list
	allUserIDs := GetDiscordUserIDDivIT()
	var availableUserIDs []string
	for _, userID := range allUserIDs {
		isAlreadySelected := slices.Contains(alreadySelected, userID)
		if !isAlreadySelected {
			availableUserIDs = append(availableUserIDs, userID)
		}
	}

	result += helpers.GetRandomString(total-5, availableUserIDs)
	return result
}
