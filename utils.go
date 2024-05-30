package eflag

import "strings"

// SplitWithComma splits a string into a slice of strings using a comma as the separator.
// Each part is trimmed of leading and trailing whitespace.
func SplitWithComma(s string) []string {
	return SplitWith(s, ",")
}

// SplitWith splits a string into a slice of strings using the specified separator.
// Each part is trimmed of leading and trailing whitespace.
func SplitWith(s, sep string) []string {
	parts := strings.Split(s, sep)
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}
