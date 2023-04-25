package pictures

import (
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type Pic struct {
	Path string
}

func (h *Pic) Match(s string) bool {
	return handlers.MatchPrefix("!пик", s)
}

func (h *Pic) Handle(c tele.Context) error {
	f, err := randomFileFromSubdir(h.Path)
	if err != nil {
		return err
	}
	return c.Send(&tele.Photo{File: tele.FromDisk(f)})
}
