package farm

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Harvest struct {
	Universe *game.Universe
}

var harvestRe = handlers.Regexp("^!(урожай|собрать)")

func (h *Harvest) Match(c tele.Context) bool {
	return harvestRe.MatchString(c.Text())
}

func (h *Harvest) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	p := user.Harvest()
	s := format.Harvested(tu.Mention(c, user), p...)
	return c.Send(s, tele.ModeHTML)
}
