package valid

import (
	"unicode"
	"unicode/utf8"
)

// Name returns true if s is a string consisting of cyrillic
// characters, ' ' or '-'.
func Name(s string) bool {
	allowed := map[rune]bool{
		' ': true,
		'-': true,
	}
	if !utf8.ValidString(s) {
		return false
	}
	if utf8.RuneCountInString(s) > 40 {
		return false
	}
	for _, r := range s {
		if !unicode.Is(unicode.Cyrillic, r) && !allowed[r] {
			return false
		}
	}
	return true
}
