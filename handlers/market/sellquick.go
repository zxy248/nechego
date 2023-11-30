package market

import (
	"nechego/fishing"
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	"nechego/item"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type SellQuick struct {
	Universe *game.Universe
}

var sellQuickRe = handlers.Regexp("^!продать")

func (h *SellQuick) Match(c tele.Context) bool {
	return sellQuickRe.MatchString(c.Text())
}

func (h *SellQuick) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	total := 0
	var sold []*item.Item
	for _, i := range user.Inventory.List() {
		f, ok := i.Value.(*fishing.Fish)
		if !ok || f.Price() < 2000 {
			continue
		}
		n, ok := user.Sell(world, i)
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
