package pictures

import (
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type Zeus struct {
	Path string
}

func (h *Zeus) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!зевс")
}

func (h *Zeus) Handle(c tele.Context) error {
	f, err := randomFile(h.Path)
	if err != nil {
		return err
	}
	return c.Send(&tele.Photo{File: tele.FromDisk(f)})
}
