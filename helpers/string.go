package helpers

import (
	"math/rand"
	"strings"
)

func GetRandomString(length int, records []string) string {
	if length <= 0 {
		return ""
	}

	var result []string
	for range length {
		result = append(result, records[rand.Intn(len(records))])
	}
	return strings.Join(result, " ")
}
