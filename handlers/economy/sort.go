package economy

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Sort struct {
	Universe *game.Universe
}

var sortRe = handlers.Regexp("^!сорт ([0-9 ]+)")

func (h *Sort) Match(c tele.Context) bool {
	return sortRe.MatchString(c.Text())
}

func (h *Sort) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	keys := sortKeys(c.Text())
	items := handlers.GetItems(user.Inventory, keys)
	if len(items) == 0 {
		return c.Send(format.ItemNotFound)
	}
	user.Inventory.Remove(items...)
	user.Inventory.AddFront(items...)
	return c.Send(format.InventorySorted)
}

func sortKeys(s string) []int {
	m := sortRe.FindStringSubmatch(s)[1]
	return handlers.Numbers(m)
}
