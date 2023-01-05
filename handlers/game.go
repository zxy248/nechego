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

func (h *Save) Match(s string) bool {
	return saveRe.MatchString(s)
}

func (h *Save) Handle(c tele.Context) error {
	return h.Universe.SaveAll()
}
