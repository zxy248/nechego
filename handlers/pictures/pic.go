package pictures

import (
	"nechego/handlers"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Pic struct {
	Path string
}

func (h *Pic) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!пик")
}

func (h *Pic) Handle(c tele.Context) error {
	f, err := randomSubdirFile(h.Path)
	if err != nil {
		return err
	}
	return c.Send(&tele.Photo{File: tele.FromDisk(f)})
}
