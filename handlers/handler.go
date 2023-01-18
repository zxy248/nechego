package handlers

import (
	"regexp"

	tele "gopkg.in/telebot.v3"
)

type Handler interface {
	Match(s string) bool
	Handle(c tele.Context) error
}

func re(s string) *regexp.Regexp {
	return regexp.MustCompile("(?i)" + s)
}
