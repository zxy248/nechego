package fun

import (
	"github.com/zxy248/nechego/handlers"
	"math/rand"
	"strings"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Choose struct{}

var chooseRe = handlers.NewRegexp("^!(выбор|выбрать) (.+)")

func (h *Choose) Match(c tele.Context) bool {
	return chooseRe.MatchString(c.Text())
}

func (h *Choose) Handle(c tele.Context) error {
	vars := strings.Split(parseChoose(c.Text()), " или ")
	v := vars[rand.Intn(len(vars))]
	return c.Send(v)
}

func parseChoose(s string) string {
	return chooseRe.FindStringSubmatch(s)[2]
}
