package util

import (
	"regexp"
	"strings"
)

// ExtractFieldFromDetail extracts the field name from a detail string.
// The detail string should contain a pattern in the format "Key (field)=value".
// Returns the field name if found, otherwise returns an empty string.
func ExtractFieldFromDetail(detail string) string {
	detail = strings.TrimSpace(detail)

	// Compile the regular expression to find the pattern "Key (field)="
	re := regexp.MustCompile(`(?i)Key\s*\(\s*(.*?)\s*\)\s*=`)
	matches := re.FindStringSubmatch(detail)

	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
