package farm

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
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

	n := user.FarmUpgradeCost()
	if n == 0 {
		return c.Send(format.MaxSizeFarm)
	}
	if !user.UpgradeFarm(n) {
		return c.Send(format.NoMoney)
	}
	s := format.FarmUpgraded(tu.Mention(c, user), user.Farm, n)
	return c.Send(s, tele.ModeHTML)
}
