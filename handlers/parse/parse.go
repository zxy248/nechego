package parse

import (
	"strconv"
	"strings"
)

// F parses a string, returning true on success or false on failure.
type F func(s string) (ok bool)

// G parses a sequence of strings, returning the remaining unparsed
// part of the sequence and a status.
type G func(seq []string) (rest []string, ok bool)

// Fseq returns a function that parses each token of a string with a
// corresponding parser.
func Fseq(p ...F) F {
	return func(s string) bool {
		tokens := strings.Fields(s)
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

// Or returns a function that sequentially parses its argument with
// the provided parsers, returning true on the first match.
func Or(p ...F) F {
	return func(s string) bool {
		for _, q := range p {
			if q(s) {
				return true
			}
		}
		return false
	}
}

// Str returns a function that matches a string.
func Str(s string) F {
	return func(t string) bool {
		return strings.EqualFold(s, t)
	}
}

// Prefix returns a function that matches a prefix.
func Prefix(s string) F {
	return func(t string) bool {
		return len(t) >= len(s) && strings.EqualFold(s, t[:len(s)])
	}
}

// Int returns a function that parses an integer and calls f with the
// result.
func Int(f func(int)) F {
	return func(s string) bool {
		n, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		f(n)
		return true
	}
}

// Interval returns a function that parses an interval (in the form
// `n-m`, where n and m are integers) and calls f with the result.
func Interval(f func(min, max int)) F {
	return func(s string) bool {
		before, after, found := strings.Cut(s, "-")
		if !found {
			return false
		}
		min, err := strconv.Atoi(before)
		if err != nil {
			return false
		}
		max, err := strconv.Atoi(after)
		if err != nil {
			return false
		}
		if min > max {
			min, max = max, min
		}
		f(min, max)
		return true
	}
}

// Gseq returns a function that sequentially parses the sequence of
// tokens in a string with the provided parsers.
func Gseq(p ...G) F {
	return func(s string) bool {
		tokens := strings.Fields(s)
		for _, q := range p {
			if len(tokens) == 0 {
				return false
			}
			rest, ok := q(tokens)
			if !ok {
				return false
			}
			tokens = rest
		}
		return true
	}
}

// Ftog wraps F in G that consumes the first token of seq.
func Ftog(f F) G {
	return func(seq []string) ([]string, bool) {
		if !f(seq[0]) {
			return seq, false
		}
		return seq[1:], true
	}
}

// All returns a function that parses each token of seq with p until
// it is empty.
func All(p F) G {
	return func(seq []string) ([]string, bool) {
		for len(seq) > 0 {
			if !p(seq[0]) {
				return seq, false
			}
			seq = seq[1:]
		}
		return []string{}, true
	}
}

// Assign returns a function that assigns its argument to x.
func Assign[T any](x *T) func(T) {
	return func(y T) {
		*x = y
	}
}
