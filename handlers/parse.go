package handlers

import (
	"regexp"
	"strings"
)

func NewRegexp(pattern string) *regexp.Regexp {
	return regexp.MustCompile("(?i)" + pattern)
}

func HasPrefix(s string, ps ...string) bool {
	for _, p := range ps {
		if strings.HasPrefix(strings.ToLower(s), p) {
			return true
		}
	}
	return false
}
