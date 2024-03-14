package fun

import (
	"fmt"
	"regexp"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Photon struct{}

var photonRe = regexp.MustCompile("^!фотон")

func (h *Photon) Match(c tele.Context) bool {
	return photonRe.MatchString(c.Text())
}

func (h *Photon) Handle(c tele.Context) error {
	candidate := [...]string{"Даванков", "Харитонов", "Слуцкий"}[c.Sender().ID%3]
	out := "<blockquote><b>🇷🇺 Выборы Президента РФ</b></blockquote>\n" +
		"Ваш кандидат: <tg-spoiler><b>%s ☑</b></tg-spoiler>"
	return c.Send(fmt.Sprintf(out, candidate), tele.ModeHTML)
}
