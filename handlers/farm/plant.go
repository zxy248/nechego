package farm

import (
	"nechego/farm/plant"
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	"nechego/item"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Plant struct {
	Universe *game.Universe
}

var plantRe = handlers.Regexp("^!посадить ([0-9 ]+)")

func (h *Plant) Match(c tele.Context) bool {
	return plantRe.MatchString(c.Text())
}

func (h *Plant) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	planted := []*plant.Plant{}
	keys := plantKeys(c.Text())
	for _, key := range keys {
		item, ok := user.Inventory.ByKey(key)
		if !ok {
			c.Send(format.BadKey(key), tele.ModeHTML)
			break
		}
		p := plantItem(user, item)
		if p == nil {
			c.Send(format.CannotPlant(item), tele.ModeHTML)
			break
		}
		planted = append(planted, p)
	}
	s := format.Planted(tu.LinkSender(c), planted...)
	return c.Send(s, tele.ModeHTML)
}

func plantKeys(s string) []int {
	m := plantRe.FindStringSubmatch(s)[1]
	return handlers.Numbers(m)
}

func plantItem(u *game.User, i *item.Item) *plant.Plant {
	p, ok := i.Value.(*plant.Plant)
	if !ok {
		return nil
	}
	planted := 0
	for p.Count > 0 && u.Farm.Plant(p.Type) {
		planted++
		p.Count--
	}
	if p.Count == 0 && !u.Inventory.Remove(i) {
		panic("cannot remove zero Plant from inventory")
	}
	return &plant.Plant{Type: p.Type, Count: planted}
}
