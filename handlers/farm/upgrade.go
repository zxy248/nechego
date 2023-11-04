package farm

import (
	"nechego/farm"
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

const (
	maxRows    = 6
	maxColumns = 6
	costFactor = 100000
)

type Upgrade struct {
	Universe *game.Universe
}

var upgradeRe = handlers.Regexp("^!апгрейд")

func (h *Upgrade) Match(c tele.Context) bool {
	return upgradeRe.MatchString(c.Text())
}

func (h *Upgrade) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	n := upgradeCost(user.Farm)
	if n == 0 {
		return c.Send(format.MaxSizeFarm)
	}
	if !upgradeFarm(user, n) {
		return c.Send(format.NoMoney)
	}
	s := format.FarmUpgraded(tu.LinkSender(c), user.Farm, n)
	return c.Send(s, tele.ModeHTML)
}

func upgradeCost(f *farm.Farm) int {
	if f.Rows >= maxRows && f.Columns >= maxColumns {
		return 0
	}
	return costFactor * fib(f.Columns+f.Rows)
}

func upgradeFarm(u *game.User, cost int) bool {
	if !u.Balance().Spend(cost) {
		return false
	}
	u.Farm.Grow()
	return true
}

func fib(n int) int {
	prev, curr := 0, 1
	for ; n > 0; n-- {
		prev, curr = curr, prev+curr
	}
	return prev
}
