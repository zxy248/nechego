package economy

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	"nechego/item"
	tu "nechego/teleutil"
	"strconv"

	tele "gopkg.in/telebot.v3"
)

type Split struct {
	Universe *game.Universe
}

var splitRe = handlers.Regexp("^!(отложить|разделить) ([0-9]+) ([0-9]+)")

func (h *Split) Match(c tele.Context) bool {
	return splitRe.MatchString(c.Text())
}

func (h *Split) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if handlers.FullInventory(user.Inventory) {
		return c.Send(format.InventoryOverflow)
	}
	key, count, _ := splitKeyCount(c.Text())
	orig, ok := user.Inventory.ByKey(key)
	if !ok {
		return c.Send(format.BadKey(key), tele.ModeHTML)
	}
	part, ok := item.Split(orig, count)
	if !ok {
		return c.Send(format.CannotSplit(orig), tele.ModeHTML)
	}
	user.Inventory.Add(part)
	l := tu.Link(c, user)
	s := format.Splitted(l, part)
	return c.Send(s, tele.ModeHTML)
}

func splitKeyCount(s string) (key, count int, ok bool) {
	m := splitRe.FindStringSubmatch(s)
	if m == nil {
		return 0, 0, false
	}
	key, _ = strconv.Atoi(m[2])
	count, _ = strconv.Atoi(m[3])
	return key, count, true
}
