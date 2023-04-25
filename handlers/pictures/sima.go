package pictures

import (
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type Sima struct{}

func (h *Sima) Match(s string) bool {
	return handlers.MatchPrefix("!сима", s)
}

func (h *Sima) Handle(c tele.Context) error {
	s, err := randomSticker(c, "catsima_vk")
	if err != nil {
		return err
	}
	return c.Send(s)
}
