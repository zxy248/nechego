package farm

import (
	"nechego/farm/plant"
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
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

	ps := []*plant.Plant{}
	keys := handlers.Numbers(plantRe.FindStringSubmatch(c.Text())[1])
	for _, key := range keys {
		item, ok := user.Inventory.ByKey(key)
		if !ok {
			c.Send(format.BadKey(key), tele.ModeHTML)
			break
		}
		p := user.Plant(item)
		if p.Count == 0 {
			c.Send(format.CannotPlant(item), tele.ModeHTML)
			break
		}
		ps = append(ps, p)
	}
	s := format.Planted(tu.Mention(c, user), ps...)
	return c.Send(s, tele.ModeHTML)
}
