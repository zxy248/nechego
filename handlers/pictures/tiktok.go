package pictures

import (
	"github.com/zxy248/nechego/handlers"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Tiktok struct {
	Path string
}

func (h *Tiktok) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!тикток")
}

func (h *Tiktok) Handle(c tele.Context) error {
	name, err := randomFile(h.Path)
	if err != nil {
		return err
	}
	return c.Send(&tele.Video{File: tele.FromDisk(name)})
}
