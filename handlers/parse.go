package handlers

import "nechego/handlers/parse"

type callback interface {
	encode() string
	decode(s string) error
}

func callbackMatch(c callback, s string) bool {
	return c.decode(s) == nil
}

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
