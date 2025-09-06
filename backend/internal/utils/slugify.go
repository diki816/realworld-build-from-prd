package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Slugify converts a string to a URL-friendly slug
func Slugify(s string) string {
	if s == "" {
		return ""
	}

	// Normalize unicode characters
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	slug, _, _ := transform.String(t, s)

	// Convert to lowercase
	slug = strings.ToLower(slug)

	// Replace non-alphanumeric characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	// Limit length to 100 characters
	if len(slug) > 100 {
		slug = slug[:100]
		slug = strings.Trim(slug, "-")
	}

	return slug
}

// GenerateUniqueSlug creates a unique slug by appending a timestamp if needed
func GenerateUniqueSlug(title string, checkExists func(string) bool) string {
	baseSlug := Slugify(title)
	if baseSlug == "" {
		baseSlug = "article"
	}

	slug := baseSlug
	
	// Check if slug exists and modify if necessary
	if checkExists(slug) {
		// Append timestamp to make it unique
		timestamp := time.Now().Unix()
		slug = fmt.Sprintf("%s-%d", baseSlug, timestamp)
		
		// If still exists (very unlikely), append random number
		if checkExists(slug) {
			slug = fmt.Sprintf("%s-%d-%d", baseSlug, timestamp, time.Now().Nanosecond()%1000)
		}
	}

	return slug
}