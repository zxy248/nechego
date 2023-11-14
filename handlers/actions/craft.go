package actions

import (
	"nechego/format"
	"nechego/game"
	"nechego/game/recipes"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Craft struct {
	Universe *game.Universe
}

var craftRe = handlers.Regexp("^!крафт ([0-9 ]+)")

func (h *Craft) Match(c tele.Context) bool {
	return craftRe.MatchString(c.Text())
}

func (h *Craft) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	keys := craftKeys(c.Text())
	recipe := handlers.GetItems(user.Inventory, keys)
	crafted, ok := recipes.Craft(user.Inventory, recipe)
	if !ok {
		return c.Send(format.CannotCraft)
	}
	l := tu.Link(c, user)
	s := format.Crafted(l, crafted...)
	return c.Send(s, tele.ModeHTML)
}

func craftKeys(s string) []int {
	m := craftRe.FindStringSubmatch(s)[1]
	return handlers.Numbers(m)
}
