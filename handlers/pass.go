package handlers

import (
	tele "gopkg.in/telebot.v3"
)

type Pass struct{}

func (h *Pass) Match(c tele.Context) bool {
	return true
}

func (h *Pass) Handle(c tele.Context) error {
	return nil
}
