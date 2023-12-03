package pictures

import (
	"math/rand"
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type Poppy struct{}

func (h *Poppy) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!паппи")
}

func (h *Poppy) Handle(c tele.Context) error {
	sets := [...]string{"pappy2_vk", "poppy_vk"}
	set := sets[rand.Intn(len(sets))]
	s, err := randomSticker(c, set)
	if err != nil {
		return err
	}
	return c.Send(s)
}
