package handlers

import (
	tele "gopkg.in/telebot.v3"
)

type Handler interface {
	Match(s string) bool
	Handle(c tele.Context) error
}
