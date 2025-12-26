package helpers

import (
	"strings"
	"unicode"
)

func Slugify(s string, id int64) string {
	s = strings.ToLower(s)

	var b strings.Builder
	lastWasDash := false

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(r)
			lastWasDash = false
			continue
		}

		if !lastWasDash {
			b.WriteRune('-')
			lastWasDash = true
		}
	}

	result := b.String()
	result = strings.Trim(result, "-")

	return result
}
