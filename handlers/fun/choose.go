package fun

import (
	"math/rand"
	"nechego/handlers"
	"strings"

	tele "gopkg.in/telebot.v3"
)

type Choose struct{}

var chooseRe = handlers.Regexp("^!(выбор|выбрать) (.+)")

func (h *Choose) Match(c tele.Context) bool {
	return chooseRe.MatchString(c.Text())
}

func (h *Choose) Handle(c tele.Context) error {
	vars := strings.Split(parseChoose(c.Text()), "или")
	v := vars[rand.Intn(len(vars))]
	return c.Send(v)
}

func parseChoose(s string) string {
	return chooseRe.FindStringSubmatch(s)[2]
}
