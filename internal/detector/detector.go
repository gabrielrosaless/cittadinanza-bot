package detector

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// HasKeyword returns true if the given text contains any of the keywords.
// Comparison is case-insensitive and accent-insensitive (e.g. "ciudadanía" == "ciudadania").
func HasKeyword(text string, keywords []string) bool {
	normalized := normalize(text)
	for _, kw := range keywords {
		if strings.Contains(normalized, normalize(kw)) {
			return true
		}
	}
	return false
}

// normalize converts a string to lowercase and strips diacritical marks (accents).
// e.g. "Ciudadanía" → "ciudadania", "Habilitación" → "habilitacion"
func normalize(s string) string {
	// NFD decomposes accented characters into base character + combining mark.
	// The runes.Remove filter then strips all combining marks (category Mn).
	// NFC recomposes back to a canonical form.
	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)
	result, _, err := transform.String(t, s)
	if err != nil {
		// Fallback: return lowercased original if transformation fails
		return strings.ToLower(s)
	}
	return strings.ToLower(result)
}
