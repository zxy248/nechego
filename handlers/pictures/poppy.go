package pictures

import (
	"math/rand/v2"

	"github.com/zxy248/nechego/handlers"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Poppy struct{}

func (h *Poppy) Match(c tele.Context) bool {
	return handlers.HasPrefix(c.Text(), "!паппи")
}

func (h *Poppy) Handle(c tele.Context) error {
	sets := [...]string{"pappy2_vk", "poppy_vk"}
	set := sets[rand.N(len(sets))]
	s, err := randomSticker(c, set)
	if err != nil {
		return err
	}
	return c.Send(s)
}
