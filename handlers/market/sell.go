package market

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	"nechego/item"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Sell struct {
	Universe *game.Universe
}

var sellRe = handlers.Regexp("^!продать ([0-9 ]+)")

func (h *Sell) Match(c tele.Context) bool {
	return sellRe.MatchString(c.Text())
}

func (h *Sell) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	var sold []*item.Item
	total := 0
	keys := sellKeys(c.Text())
	for _, key := range keys {
		i, ok := user.Inventory.ByKey(key)
		if !ok {
			c.Send(format.BadKey(key), tele.ModeHTML)
			break
		}
		n, ok := user.Sell(world.Market, i)
		if !ok {
			c.Send(format.CannotSell(i), tele.ModeHTML)
			break
		}
		total += n
		sold = append(sold, i)
	}
	l := tu.Link(c, user)
	s := format.Sold(l, total, sold)
	return c.Send(s, tele.ModeHTML)
}

func sellKeys(s string) []int {
	m := sellRe.FindStringSubmatch(s)[1]
	return handlers.Numbers(m)
}
