package handlers

import (
	"nechego/farm/plant"
	"nechego/format"
	"nechego/game"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Farm struct {
	Universe *game.Universe
}

var farmRe = re("^!(ферма|огород|грядка)")

func (h *Farm) Match(s string) bool {
	return farmRe.MatchString(s)
}

func (h *Farm) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	return c.Send(format.Farm(tu.Mention(c, user), user.Farm), tele.ModeHTML)
}

type Plant struct {
	Universe *game.Universe
}

var plantRe = re(`^!посадить (.*)`)

func (h *Plant) Match(s string) bool {
	return plantRe.MatchString(s)
}

func (h *Plant) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	planted := []*plant.Plant{}
	for _, key := range tu.NumArg(c, plantRe, 1) {
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
		planted = append(planted, p)
	}
	return c.Send(format.Planted(tu.Mention(c, user), planted...), tele.ModeHTML)
}

type Harvest struct {
	Universe *game.Universe
}

var gatherRe = re("^!урожай")

func (h *Harvest) Match(s string) bool {
	return gatherRe.MatchString(s)
}

func (h *Harvest) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	harvested := user.Harvest()
	return c.Send(format.Harvested(tu.Mention(c, user), harvested...), tele.ModeHTML)
}

type FarmSize struct {
	Universe *game.Universe
}

var farmSizeRe = re("^!земля")

func (h *FarmSize) Match(s string) bool {
	return farmSizeRe.MatchString(s)
}

func (h *FarmSize) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	return c.Send(format.FarmSize(user.Farm, user.FarmUpgradeCost()), tele.ModeHTML)
}

type UpgradeFarm struct {
	Universe *game.Universe
}

var upgradeFarmRe = re("^!апгрейд")

func (h *UpgradeFarm) Match(s string) bool {
	return upgradeFarmRe.MatchString(s)
}

func (h *UpgradeFarm) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	cost := user.FarmUpgradeCost()
	if !user.UpgradeFarm() {
		return c.Send(format.NoMoney)
	}
	return c.Send(format.FarmUpgraded(tu.Mention(c, user), user.Farm, cost), tele.ModeHTML)
}
