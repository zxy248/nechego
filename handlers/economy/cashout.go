package economy

import (
	"errors"
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	"nechego/money"
	tu "nechego/teleutil"
	"strconv"

	tele "gopkg.in/telebot.v3"
)

type Cashout struct {
	Universe *game.Universe
}

var cashoutRe = handlers.Regexp("^!(отложить|обнал|снять) ([0-9]+)")

func (h *Cashout) Match(c tele.Context) bool {
	return cashoutRe.MatchString(c.Text())
}

func (h *Cashout) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	n := cashoutAmount(c.Text())
	err := user.Balance().Cashout(n)
	if errors.Is(err, money.ErrBadMoney) {
		return c.Send(format.BadMoney)
	} else if errors.Is(err, money.ErrNoMoney) {
		return c.Send(format.NoMoney)
	} else if err != nil {
		return err
	}
	l := tu.Link(c, user)
	s := format.Cashout(l, n)
	return c.Send(s, tele.ModeHTML)
}

func cashoutAmount(s string) int {
	m := cashoutRe.FindStringSubmatch(s)[2]
	n, _ := strconv.Atoi(m)
	return n
}
