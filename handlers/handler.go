package handlers

import (
	"regexp"

	tele "gopkg.in/telebot.v3"
)

type Handler interface {
	Regexp() *regexp.Regexp
	Handle(c tele.Context) error
}
