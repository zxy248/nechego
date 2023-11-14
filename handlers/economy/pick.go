package economy

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Pick struct {
	Universe *game.Universe
}

var pickRe = handlers.Regexp("^!(взять|подобрать|поднять) ([0-9 ]+)")

func (h *Pick) Match(c tele.Context) bool {
	return pickRe.MatchString(c.Text())
}

func (h *Pick) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if handlers.FullInventory(user.Inventory) {
		return c.Send(format.InventoryOverflow)
	}
	keys := pickKeys(c.Text())
	items := handlers.GetItems(world.Floor, keys)
	pick, bad := handlers.MoveItems(user.Inventory, world.Floor, items)
	if bad != nil {
		c.Send(format.CannotPick(bad), tele.ModeHTML)
	}
	l := tu.Link(c, user)
	s := format.Picked(l, pick)
	return c.Send(s, tele.ModeHTML)
}

func pickKeys(s string) []int {
	m := pickRe.FindStringSubmatch(s)[2]
	return handlers.Numbers(m)
}
