package market

import (
	"errors"
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	"nechego/item"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Buy struct {
	Universe *game.Universe
}

var buyRe = handlers.Regexp("^!купить ([0-9 ]+)")

func (h *Buy) Match(c tele.Context) bool {
	return buyRe.MatchString(c.Text())
}

func (h *Buy) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if handlers.FullInventory(user.Inventory) {
		return c.Send(format.InventoryOverflow)
	}
	var bought []*item.Item
	total := 0
	keys := buyKeys(c.Text())
	for _, key := range keys {
		p, err := user.Buy(world.Market, key)
		if errors.Is(err, game.ErrNoKey) {
			c.Send(format.BadKey(key), tele.ModeHTML)
			break
		} else if err != nil {
			c.Send(format.NoMoney, tele.ModeHTML)
			break
		}
		total += p.Price
		bought = append(bought, p.Item)
	}
	l := tu.Link(c, user)
	s := format.Bought(l, total, bought)
	return c.Send(s, tele.ModeHTML)
}

func buyKeys(s string) []int {
	m := buyRe.FindStringSubmatch(s)[1]
	return handlers.Numbers(m)
}
