package eflag

import (
	"strings"
	"unicode"
)

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

// MixedCapsToScreamingSnake converts a mixedCaps string to a MIXED_CAPS string
func MixedCapsToScreamingSnake(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 && !unicode.IsUpper(rune(s[i-1])) {
			result = append(result, '_', r)
		} else {
			result = append(result, unicode.ToUpper(r))
		}
	}
	return string(result)
}
