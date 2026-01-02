package models

import (
	"math/rand"
	"slices"
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

func GetMemberManagerDivIT() []DiscordMember {
	return membersManagerDivIT
}

func GetMemberSreDivIT() []DiscordMember {
	return membersSreDivIT
}

func GetMemberForceDivIT() []DiscordMember {
	return membersForceDivIT
}

func GetMemberSdetDivIT() []DiscordMember {
	return membersSdetDivIT
}

func GetMemberPeDivIT() []DiscordMember {
	return membersPeDivIT
}

func GetMemberDivIT() []DiscordMember {
	var members []DiscordMember
	members = append(members, GetMemberManagerDivIT()...)
	members = append(members, GetMemberSreDivIT()...)
	members = append(members, GetMemberForceDivIT()...)
	members = append(members, GetMemberSdetDivIT()...)
	members = append(members, GetMemberPeDivIT()...)
	return members
}

func extractUserIDs(members []DiscordMember) []string {
	var userID []string
	for _, member := range members {
		userID = append(userID, member.UserID)
	}
	return userID
}

func getRandomMembers(length int, records []DiscordMember) []DiscordMember {
	if length <= 0 || len(records) == 0 {
		return nil
	}

	var result []DiscordMember
	for range length {
		result = append(result, records[rand.Intn(len(records))])
	}
	return result
}

func GetRandomMemberDivIT(total int) []DiscordMember {
	allMembers := GetMemberDivIT()
	if total < 5 {
		return getRandomMembers(total, allMembers)
	}

	var result []DiscordMember
	var alreadySelectedIDs []string

	selectRandom := func(source []DiscordMember) {
		if len(source) == 0 {
			return
		}
		selected := source[rand.Intn(len(source))]
		result = append(result, selected)
		alreadySelectedIDs = append(alreadySelectedIDs, selected.UserID)
	}

	selectRandom(GetMemberManagerDivIT())
	selectRandom(GetMemberSreDivIT())
	selectRandom(GetMemberSdetDivIT())
	selectRandom(GetMemberPeDivIT())
	selectRandom(GetMemberForceDivIT())

	// Filter out already selected members from the full list
	var availableMembers []DiscordMember
	for _, member := range allMembers {
		if !slices.Contains(alreadySelectedIDs, member.UserID) {
			availableMembers = append(availableMembers, member)
		}
	}

	remaining := total - 5
	if remaining > 0 {
		result = append(result, getRandomMembers(remaining, availableMembers)...)
	}

	return result
}
