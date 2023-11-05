package handlers

import (
	"nechego/handlers/parse"
	"regexp"
	"strconv"
	"strings"
)

func numCommand(prefix parse.G, s string) (keys []int, ok bool) {
	ok = parse.Seq(
		prefix,
		parse.All(parse.Or(
			parse.Int(func(n int) {
				keys = append(keys, n)
			}),
			parse.Interval(func(min, max int) {
				const lim = 20
				if max-min > lim {
					max = min + lim
				}
				for i := min; i <= max; i++ {
					keys = append(keys, i)
				}
			}),
		)),
	)(s)
	return
}

func textCommand(prefix parse.G, s string) (text string, ok bool) {
	ok = parse.Seq(prefix, parse.Str(parse.Assign(&text)))(s)
	return
}

func Regexp(pattern string) *regexp.Regexp {
	return regexp.MustCompile("(?i)" + pattern)
}

func MatchRegexp(pattern, s string) bool {
	return Regexp(pattern).MatchString(s)
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
