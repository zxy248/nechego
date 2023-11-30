package market

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Market struct {
	Universe *game.Universe
}

var marketRe = handlers.Regexp("^!(магаз|шоп)")

func (h *Market) Match(c tele.Context) bool {
	return marketRe.MatchString(c.Text())
}

func (h *Market) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	var l string
	if id, ok := world.Market.Shift.Worker(); ok {
		l = tu.Link(c, id)
	}
	s := format.Market(l, world.Market)
	return c.Send(s, tele.ModeHTML)
}
