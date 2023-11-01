package handlers

import (
	"nechego/farm"
	"nechego/farm/plant"
	"nechego/format"
	"nechego/game"
	"nechego/handlers/parse"
	tu "nechego/teleutil"
	"nechego/valid"
	"strings"

	tele "gopkg.in/telebot.v3"
)

type Farm struct {
	Universe *game.Universe
}

var farmRe = Regexp("^!(ферма|огород|грядка)")

func (h *Farm) Match(s string) bool {
	return farmRe.MatchString(s)
}

func (h *Farm) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	upgradeCost, _ := user.FarmUpgradeCost()
	return c.Send(format.Farm(tu.Mention(c, user), user.Farm, upgradeCost),
		farmInlineKeyboard(user.TUID, user.Farm), tele.ModeHTML)
}

func farmInlineKeyboard(id int64, f *farm.Farm) *tele.ReplyMarkup {
	grid := &tele.ReplyMarkup{}
	rows := []tele.Row{}
	for r := 0; r < f.Rows; r++ {
		buttons := []tele.Btn{}
		for c := 0; c < f.Columns; c++ {
			emoji := f.Grid[farm.Plot{Row: r, Column: c}].String()
			btn := grid.Data(emoji, "0")
			buttons = append(buttons, btn)
		}
		rows = append(rows, grid.Row(buttons...))
		buttons = []tele.Btn{}
	}
	grid.Inline(rows...)
	return grid
}

type Plant struct {
	Universe *game.Universe
}

var plantRe = Regexp(`^!посадить (.*)`)

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

var harvestRe = Regexp("^!(урожай|собрать)")

func (h *Harvest) Match(s string) bool {
	return harvestRe.MatchString(s)
}

func (h *Harvest) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	harvested := user.Harvest()
	return c.Send(format.Harvested(tu.Mention(c, user), harvested...), tele.ModeHTML)
}

type UpgradeFarm struct {
	Universe *game.Universe
}

var upgradeFarmRe = Regexp("^!апгрейд")

func (h *UpgradeFarm) Match(s string) bool {
	return upgradeFarmRe.MatchString(s)
}

func (h *UpgradeFarm) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	cost, ok := user.FarmUpgradeCost()
	if !ok {
		return c.Send(format.MaxSizeFarm)
	}
	if !user.UpgradeFarm(cost) {
		return c.Send(format.NoMoney)
	}
	return c.Send(format.FarmUpgraded(tu.Mention(c, user), user.Farm, cost), tele.ModeHTML)
}

type NameFarm struct {
	Universe *game.Universe
}

func (h *NameFarm) Match(s string) bool {
	_, ok := nameFarmCommand(s)
	return ok
}

func (h *NameFarm) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	name, ok := nameFarmCommand(c.Text())
	if !ok {
		panic("bad name farm command")
	}
	if !valid.Name(name) {
		return c.Send(format.BadFarmName)
	}
	user.Farm.Name = strings.Title(name)
	return c.Send(format.FarmNamed(tu.Mention(c, user), user.Farm), tele.ModeHTML)
}

func nameFarmCommand(s string) (name string, ok bool) {
	ok = parse.Seq(
		parse.Match("!назвать"), parse.Match("ферму"),
		parse.Str(parse.Assign(&name)),
	)(s)
	return
}
