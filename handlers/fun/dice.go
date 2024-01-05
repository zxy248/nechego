package fun

import (
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type Dice struct{}

var diceRe = handlers.NewRegexp("^!кости")

func (h *Dice) Match(c tele.Context) bool {
	return diceRe.MatchString(c.Text())
}

func (h *Dice) Handle(c tele.Context) error {
	return c.Send(tele.Cube)
}
