package danbooru

import "fmt"

// Rating of a picture.
type Rating int

const (
	General Rating = iota
	Sensitive
	Questionable
	Explicit
)

// MarshalText implements the encoding.TextMarshaler interface.
func (r Rating) MarshalText() (text []byte, err error) {
	switch r {
	case General:
		return []byte("g"), nil
	case Sensitive:
		return []byte("s"), nil
	case Questionable:
		return []byte("q"), nil
	case Explicit:
		return []byte("e"), nil
	}
	return nil, fmt.Errorf("cannot marshal %v", r)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (r *Rating) UnmarshalText(text []byte) error {
	switch string(text) {
	case "g":
		*r = General
	case "s":
		*r = Sensitive
	case "q":
		*r = Questionable
	case "e":
		*r = Explicit
	default:
		return fmt.Errorf("cannot unmarshal %s", text)
	}
	return nil
}

// NSFW returns true if the rating is not safe for work.
func (r Rating) NSFW() bool {
	return r == Questionable || r == Explicit
}
