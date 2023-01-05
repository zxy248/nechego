package handlers

import (
	"nechego/game"
	"regexp"

	tele "gopkg.in/telebot.v3"
)

type Save struct {
	Universe *game.Universe
}

var saveRe = regexp.MustCompile("!сохран")

func (h *Save) Regexp() *regexp.Regexp {
	return saveRe
}

func (h *Save) Handle(c tele.Context) error {
	return h.Universe.SaveAll()
}
