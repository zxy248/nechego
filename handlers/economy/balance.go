package economy

import (
	"fmt"
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Balance struct {
	Universe *game.Universe
}

var balanceRe = handlers.Regexp("^!(баланс|деньги)")

func (h *Balance) Match(c tele.Context) bool {
	return balanceRe.MatchString(c.Text())
}

func (h *Balance) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	n := user.Balance().Total()
	s := fmt.Sprintf("💵 Ваш баланс: %s", format.Money(n))
	return c.Send(s, tele.ModeHTML)
}
