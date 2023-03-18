package parse

import (
	"strconv"
	"strings"
)

// F is a functions that tries to parse its argument.
// Returns true on success or false on failure.
type F func(string) bool

// Command returns a function that parses a sequence of tokens.
func Command(p ...F) F {
	return func(s string) bool {
		tokens := strings.Split(s, " ")
		if len(tokens) != len(p) {
			return false
		}
		for i, q := range p {
			if !q(tokens[i]) {
				return false
			}
		}
		return true
	}
}

// Str returns a function that matches a string.
func Str(s string) F {
	return func(t string) bool {
		return s == t
	}
}

// Int returns a function that parses an integer.
func Int(n *int) F {
	return func(s string) bool {
		m, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		*n = m
		return true
	}
}
