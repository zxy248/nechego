package pictures

import (
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type Casper struct {
	Path string
}

func (h *Casper) Match(s string) bool {
	return handlers.MatchRegexp("^!касп[ие]р", s)
}

func (h *Casper) Handle(c tele.Context) error {
	f, err := randomFile(h.Path)
	if err != nil {
		return err
	}
	return c.Send(&tele.Photo{File: tele.FromDisk(f)})
}
