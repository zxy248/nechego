package danbooru

import "fmt"

type Rating int

const (
	General Rating = iota
	Sensitive
	Questionable
	Explicit
)

func rate(s string) Rating {
	switch s {
	case "g":
		return General
	case "s":
		return Sensitive
	case "q":
		return Questionable
	case "e":
		return Explicit
	}
	panic(fmt.Errorf("unexpected rating %s", s))
}

func (r Rating) NSFW() bool {
	return r == Questionable || r == Explicit
}
