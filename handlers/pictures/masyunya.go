package pictures

import (
	"github.com/zxy248/nechego/handlers"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Masyunya struct{}

var masyunyaRe = handlers.NewRegexp("^!ма[нс]ю[нс][а-я]*[пая]")

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
