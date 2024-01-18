package pictures

import (
	"github.com/zxy248/nechego/handlers"

	tele "gopkg.in/zxy248/telebot.v3"
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
