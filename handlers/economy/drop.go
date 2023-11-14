package economy

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Drop struct {
	Universe *game.Universe
}

var dropRe = handlers.Regexp("^!(выкинуть|выбросить|выложить|дроп|положить) ([0-9 ]+)")

func (h *Drop) Match(c tele.Context) bool {
	return dropRe.MatchString(c.Text())
}

func (h *Drop) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	keys := dropKeys(c.Text())
	items := handlers.GetItems(user.Inventory, keys)
	drop, bad := handlers.MoveItems(world.Floor, user.Inventory, items)
	if bad != nil {
		c.Send(format.CannotDrop(bad), tele.ModeHTML)
	}
	world.Floor.Trim(10)
	l := tu.Link(c, user)
	s := format.Dropped(l, drop)
	return c.Send(s, tele.ModeHTML)
}

func dropKeys(s string) []int {
	m := dropRe.FindStringSubmatch(s)[2]
	return handlers.Numbers(m)
}
