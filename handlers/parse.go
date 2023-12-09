package handlers

import (
	"regexp"
	"strconv"
	"strings"
)

func Regexp(pattern string) *regexp.Regexp {
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

func Numbers(s string) []int {
	const lim = 10
	var ns []int
	for i, x := range strings.Fields(s) {
		if i == lim {
			break
		}
		n, err := strconv.Atoi(x)
		if err != nil {
			break
		}
		ns = append(ns, n)
	}
	return ns
}
