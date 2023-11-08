package actions

import (
	"nechego/fishing"
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Fish struct {
	Universe *game.Universe
}

var fishRe = handlers.Regexp("^!(р[ыі]балка|ловля рыб)")

func (h *Fish) Match(c tele.Context) bool {
	return fishRe.MatchString(c.Text())
}

func (h *Fish) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if fullInventory(user.Inventory) {
		return c.Send(format.InventoryOverflow)
	}
	rod, ok := user.FishingRod()
	if !ok {
		return c.Send(format.BuyFishingRod)
	}
	if !user.Energy.Spend(0.2) {
		return c.Send(format.NoEnergy)
	}
	item, caught := user.Fish(rod)
	if rod.Broken() {
		c.Send(format.FishingRodBroke)
	}
	if !caught {
		return c.Send(format.BadFishOutcome())
	}
	if f, ok := item.Value.(*fishing.Fish); ok {
		world.History.Add(user.ID, f)
		announceRecordCatch(c, world, f)
	}
	user.Inventory.Add(item)
	s := format.FishCatch(tu.Link(c, user), item)
	return c.Send(s, tele.ModeHTML)
}

func announceRecordCatch(c tele.Context, w *game.World, f *fishing.Fish) {
	for p, e := range w.History.Records {
		if e.Fish == f {
			s := format.RecordCatch(p, e)
			c.Send(s, tele.ModeHTML)
		}
	}
}
