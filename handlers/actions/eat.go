package actions

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	"nechego/item"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Eat struct {
	Universe *game.Universe
}

var eatRe = handlers.Regexp("^!(с[ъь]есть|еда) ([0-9 ]+)")

func (h *Eat) Match(c tele.Context) bool {
	return eatRe.MatchString(c.Text())
}

func (h *Eat) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if user.Energy.Full() {
		return c.Send(format.NotHungry)
	}

	var ate []*item.Item
	keys := eatKeys(c.Text())
	items := handlers.GetItems(user.Inventory, keys)
	for _, x := range items {
		if !user.Eat(x) {
			c.Send(format.CannotEat(x), tele.ModeHTML)
			break
		}
		ate = append(ate, x)
	}
	l := tu.Link(c, user)
	s1 := format.Eaten(l, ate) + "\n\n"
	s2 := format.EnergyRemaining(user.Energy)
	return c.Send(s1+s2, tele.ModeHTML)
}

func eatKeys(s string) []int {
	m := eatRe.FindStringSubmatch(s)[2]
	return handlers.Numbers(m)
}
