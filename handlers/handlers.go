package handlers

import (
	"regexp"

	tele "gopkg.in/telebot.v3"
)

//go:generate handlerplate basic.go calc.go daily.go game.go phone.go

type Handler interface {
	Match(s string) bool
	Handle(c tele.Context) error
	Self() HandlerID
}

func re(s string) *regexp.Regexp {
	return regexp.MustCompile("(?i)" + s)
}
