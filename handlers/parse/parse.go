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

// Seq returns a function that sequentially parses the sequence of
// tokens in a string with the provided parsers.
func Seq(p ...G) F {
	return func(s string) bool {
		seq := strings.Fields(s)
		for _, q := range p {
			if len(seq) == 0 {
				return false
			}
			rest, ok := q(seq)
			if !ok {
				return false
			}
			seq = rest
		}
		return true
	}
}

// Parse returns a function that parses a sequence of tokens.
func Parse(p ...G) F {
	return func(s string) bool {
		seq := strings.Fields(s)
		for _, q := range p {
			rest, ok := q(seq)
			if !ok {
				return false
			}
			seq = rest
		}
		return true
	}
}

// Or returns a function that parses its argument with the provided
// parsers, returning true on the first match.
func Or(p ...G) G {
	return func(seq []string) ([]string, bool) {
		for _, q := range p {
			if rest, ok := q(seq); ok {
				return rest, true
			}
		}
		return seq, false
	}
}

// Match returns a function that matches one of the given strings.
func Match(s ...string) G {
	return func(seq []string) ([]string, bool) {
		if len(seq) == 0 {
			return nil, false
		}
		for _, t := range s {
			if strings.EqualFold(t, car(seq)) {
				return cdr(seq), true
			}
		}
		return seq, false
	}
}

// Prefix returns a function that matches one of the given prefixes.
func Prefix(p ...string) G {
	return func(seq []string) ([]string, bool) {
		if len(seq) == 0 {
			return nil, false
		}
		a := car(seq)
		for _, q := range p {
			if len(a) >= len(q) && strings.EqualFold(q, a[:len(q)]) {
				return cdr(seq), true
			}
		}
		return seq, false
	}
}

// Int returns a function that parses an integer and calls f with the
// result.
func Int(f func(int)) G {
	return func(seq []string) ([]string, bool) {
		if len(seq) == 0 {
			return nil, false
		}
		n, err := strconv.Atoi(car(seq))
		if err != nil {
			return seq, false
		}
		f(n)
		return cdr(seq), true
	}
}

// Str returns a function that joins all tokens of the sequence with a
// single whitespace and calls f with the result.
func Str(f func(string)) G {
	return func(seq []string) ([]string, bool) {
		if len(seq) == 0 {
			return nil, false
		}
		f(strings.Join(seq, " "))
		return nil, true
	}
}

// Interval returns a function that parses an interval in the form
// `n-m` where n and m are integers and calls f with the result.
func Interval(f func(min, max int)) G {
	return func(seq []string) ([]string, bool) {
		if len(seq) == 0 {
			return nil, false
		}
		before, after, found := strings.Cut(car(seq), "-")
		if !found {
			return seq, false
		}
		min, err := strconv.Atoi(before)
		if err != nil {
			return seq, false
		}
		max, err := strconv.Atoi(after)
		if err != nil {
			return seq, false
		}
		if min > max {
			min, max = max, min
		}
		f(min, max)
		return cdr(seq), true
	}
}

// All returns a function that parses all tokens in the sequence with
// p until it is empty.
func All(p G) G {
	return func(seq []string) ([]string, bool) {
		for len(seq) > 0 {
			rest, ok := p(seq)
			if !ok {
				return seq, false
			}
			seq = rest
		}
		return nil, true
	}
}

// Assign returns a function that assigns its argument to x.
func Assign[T any](x *T) func(T) {
	return func(y T) {
		*x = y
	}
}

// Any returns a function that matches anything.
func Any() G {
	return func(seq []string) ([]string, bool) {
		return seq, true
	}
}

// Maybe returns a function that optionally matches with p.
func Maybe(p G) G {
	return Or(p, Any())
}

// car returns the head of the list.
func car[T any](list []T) T {
	return list[0]
}

// cdr returns the tail of the list.
func cdr[T any](list []T) []T {
	return list[1:]
}
