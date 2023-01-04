package handlers

import tele "gopkg.in/telebot.v3"

type Handler interface {
	Handle(c tele.Context) error
}

func Func(h Handler) tele.HandlerFunc {
	return h.Handle
}
