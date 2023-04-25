package pictures

import (
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type Masyunya struct{}

func (h *Masyunya) Match(s string) bool {
	return handlers.MatchRegexp("^!ма[нс]ю[нс][а-я]*[пая]", s)
}

func (h *Masyunya) Handle(c tele.Context) error {
	s, err := randomSticker(c, "masyunya_vk")
	if err != nil {
		return err
	}
	return c.Send(s)
}
