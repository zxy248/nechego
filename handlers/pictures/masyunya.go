package pictures

import (
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type Masyunya struct{}

var masyunyaRe = handlers.Regexp("^!ма[нс]ю[нс][а-я]*[пая]")

func (h *Masyunya) Match(c tele.Context) bool {
	return masyunyaRe.MatchString(c.Text())
}

func (h *Masyunya) Handle(c tele.Context) error {
	s, err := randomSticker(c, "masyunya_vk")
	if err != nil {
		return err
	}
	return c.Send(s)
}
