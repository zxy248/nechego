package economy

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Funds struct {
	Universe *game.Universe
}

var fundsRe = handlers.Regexp("^!(зарплата|средства|получить)")

func (h *Funds) Match(c tele.Context) bool {
	return fundsRe.MatchString(c.Text())
}

func (h *Funds) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if handlers.FullInventory(user.Inventory) {
		return c.Send(format.InventoryOverflow)
	}
	funds := user.Funds.Collect()
	for _, f := range funds {
		user.Inventory.Add(f.Item)
	}
	user.Inventory.Stack()
	l := tu.Link(c, user)
	s := format.FundsCollected(l, funds)
	return c.Send(s, tele.ModeHTML)
}
