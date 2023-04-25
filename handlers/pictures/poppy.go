package pictures

import (
	"math/rand"
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type Poppy struct{}

func (h *Poppy) Match(s string) bool {
	return handlers.MatchPrefix("!паппи", s)
}

func (h *Poppy) Handle(c tele.Context) error {
	names := []string{"pappy2_vk", "poppy_vk"}
	name := names[rand.Intn(len(names))]
	s, err := randomSticker(c, name)
	if err != nil {
		return err
	}
	return c.Send(s)
}
