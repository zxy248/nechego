package profile

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Energy struct {
	Universe *game.Universe
}

var energyRe = handlers.Regexp("^!энергия")

func (h *Energy) Match(c tele.Context) bool {
	return energyRe.MatchString(c.Text())
}

func (h *Energy) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	s := format.EnergyLevel(user.Energy)
	return c.Send(s, tele.ModeHTML)
}
