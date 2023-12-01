package market

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"
	"nechego/valid"
	"strings"

	tele "gopkg.in/telebot.v3"
)

type Name struct {
	Universe *game.Universe
}

var nameRe = handlers.Regexp("^!назвать магазин (.+)")

func (h *Name) Match(c tele.Context) bool {
	return nameRe.MatchString(c.Text())
}

func (h *Name) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	n := marketName(c.Text())
	if !valid.Name(n) {
		return c.Send(format.BadMarketName)
	}
	world.Market.Name = formatName(n)
	return c.Send(format.MarketRenamed)
}

func marketName(s string) string {
	return nameRe.FindStringSubmatch(s)[1]
}

func formatName(s string) string {
	return strings.Title(s)
}
