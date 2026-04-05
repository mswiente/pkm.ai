package slug

import (
	"regexp"
	"strings"
	"unicode"
)

const maxSlugLen = 64

var multiHyphen = regexp.MustCompile(`-{2,}`)

// umlauts maps German special characters to ASCII equivalents.
var umlauts = map[rune]string{
	'ä': "ae", 'ö': "oe", 'ü': "ue",
	'Ä': "ae", 'Ö': "oe", 'Ü': "ue",
	'ß': "ss",
}

// FromTitle converts a title string to a filename-safe slug.
//
// Steps:
//  1. Replace German umlauts with ASCII equivalents
//  2. Lowercase all ASCII uppercase letters
//  3. Replace spaces, underscores, dots, slashes with hyphens
//  4. Drop all characters that are not [a-z0-9-]
//  5. Collapse consecutive hyphens
//  6. Trim leading/trailing hyphens
//  7. Truncate to maxSlugLen at a hyphen boundary if possible
//
// Returns "untitled" for an empty or all-stripped input.
func FromTitle(title string) string {
	var b strings.Builder
	for _, r := range title {
		if sub, ok := umlauts[r]; ok {
			b.WriteString(sub)
			continue
		}
		if unicode.IsUpper(r) && r < 128 {
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		switch {
		case r == ' ' || r == '_' || r == '.' || r == '/':
			b.WriteRune('-')
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-':
			b.WriteRune(r)
		}
	}

	s := multiHyphen.ReplaceAllString(b.String(), "-")
	s = strings.Trim(s, "-")

	if s == "" {
		return "untitled"
	}

	if len(s) > maxSlugLen {
		s = truncateAtBoundary(s, maxSlugLen)
	}

	return s
}

// truncateAtBoundary truncates s to at most max bytes, preferring to cut at
// the last hyphen before max.
func truncateAtBoundary(s string, max int) string {
	s = s[:max]
	if idx := strings.LastIndex(s, "-"); idx > 0 {
		s = s[:idx]
	}
	return strings.TrimRight(s, "-")
}
