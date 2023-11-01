package farm

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Farm struct {
	Universe *game.Universe
}

var farmRe = handlers.Regexp("^!(ферма|огород|грядка)")

func (h *Farm) Match(c tele.Context) bool {
	return farmRe.MatchString(c.Text())
}

func (h *Farm) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	n := user.FarmUpgradeCost()
	s := format.Farm(tu.Mention(c, user), user.Farm, n)
	return c.Send(s, tele.ModeHTML)
}
