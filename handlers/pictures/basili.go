package pictures

import (
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type Basili struct {
	Path string
}

func (h *Basili) Match(s string) bool {
	return handlers.MatchRegexp("!(муся|марс|(кот|кошка) василия)", s)
}

func (h *Basili) Handle(c tele.Context) error {
	f, err := randomFile(h.Path)
	if err != nil {
		return err
	}
	return c.Send(&tele.Photo{File: tele.FromDisk(f)})
}
