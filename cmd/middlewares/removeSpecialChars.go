package middlewares

import (
	"regexp"
	"strings"
)

func RemoveSpecialChars(text string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")

	return strings.ReplaceAll(reg.ReplaceAllString(text, ""), "-", "")
}
