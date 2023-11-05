package pictures

import (
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type Mouse struct {
	Path string
}

func (h *Mouse) Match(s string) bool {
	return handlers.HasPrefix(s, "!мыш")
}

func (h *Mouse) Handle(c tele.Context) error {
	return c.Send(&tele.Video{File: tele.FromDisk(h.Path)})
}
