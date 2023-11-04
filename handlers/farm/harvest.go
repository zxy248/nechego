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

type Harvest struct {
	Universe *game.Universe
}

var harvestRe = handlers.Regexp("^!(урожай|собрать)")

func (h *Harvest) Match(c tele.Context) bool {
	return harvestRe.MatchString(c.Text())
}

func (h *Harvest) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	ps := harvest(user)
	s := format.Harvested(tu.LinkSender(c), ps...)
	return c.Send(s, tele.ModeHTML)
}

func harvest(u *game.User) []*plant.Plant {
	ps := u.Farm.Harvest()
	for _, p := range ps {
		u.Inventory.Add(item.New(p))
	}
	return ps
}
