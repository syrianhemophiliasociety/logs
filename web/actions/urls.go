package actions

import (
	"strings"

	"github.com/mozillazg/go-unidecode"
)

func Slugify(content string) string {
	cleanContent := unidecode.Unidecode(strings.ToLower(content))

	var slugBuilder strings.Builder

	for _, ch := range cleanContent {
		if ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' || ch == ' ' || ch == '-' {
			slugBuilder.WriteRune(ch)
		} else {
			slugBuilder.WriteRune('-')
		}
	}

	slug := strings.ReplaceAll(slugBuilder.String(), " ", "-")
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	slug = strings.Trim(slug, "-")

	return slug
}
