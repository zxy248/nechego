package top

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Rich struct {
	Universe *game.Universe
}

var richRe = handlers.Regexp("^!топ бога[тч]")

func (h *Rich) Match(c tele.Context) bool {
	return richRe.MatchString(c.Text())
}

func (h *Rich) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	us := trim(world.SortedUsers(game.ByWealth))
	s := format.TopRich(whoFunc(c), world, us)
	return c.Send(s, tele.ModeHTML)
}
