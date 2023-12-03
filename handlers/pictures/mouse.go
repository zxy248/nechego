package pictures

import (
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type Mouse struct {
	Path string
}

func (h *Mouse) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!мыш")
}

func (h *Mouse) Handle(c tele.Context) error {
	return c.Send(&tele.Video{File: tele.FromDisk(h.Path)})
}
