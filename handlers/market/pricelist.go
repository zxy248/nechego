package market

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type PriceList struct {
	Universe *game.Universe
}

var priceListRe = handlers.Regexp("^!(прайс-?лист|цен)")

func (h *PriceList) Match(c tele.Context) bool {
	return priceListRe.MatchString(c.Text())
}

func (h *PriceList) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	world.Market.PriceList.Refresh()
	s := format.PriceList(world.Market.PriceList)
	return c.Send(s, tele.ModeHTML)
}
