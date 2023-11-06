package top

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Strength struct {
	Universe *game.Universe
}

var strengthRe = handlers.Regexp("^!топ сил")

func (h *Strength) Match(c tele.Context) bool {
	return strengthRe.MatchString(c.Text())
}

func (h *Strength) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	us := trim(world.SortedUsers(game.ByStrength))
	s := format.TopStrength(whoFunc(c), world, us)
	return c.Send(s, tele.ModeHTML)
}
