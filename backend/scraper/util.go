package scraper

import (
	"regexp"
	"strings"
)

func extractTelephone(text string) string {
	re := regexp.MustCompile(`\+\d{1,4}[-.\s]?\(?\d{1,3}\)?[-.\s]?\d{1,4}[-.\s]?\d{1,4}[-.\s]?\d{1,9}`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 0 {
		return strings.TrimSpace(matches[0])
	}
	return ""
}

func extractEmail(text string) string {
	re := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 0 {
		return strings.TrimSpace(matches[0])
	}
	return ""
}
