package handlers

type callback interface {
	encode() string
	decode(s string) error
}

func callbackMatch(c callback, s string) bool {
	return c.decode(s) == nil
}
