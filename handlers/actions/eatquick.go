package actions

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	"nechego/item"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type EatQuick struct {
	Universe *game.Universe
}

var eatQuickRe = handlers.Regexp("^!еда")

func (h *EatQuick) Match(c tele.Context) bool {
	return eatQuickRe.MatchString(c.Text())
}

func (h *EatQuick) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	if user.Energy.Full() {
		return c.Send(format.NotHungry)
	}

	var ate []*item.Item
	for !user.Energy.Full() {
		x, ok := user.EatQuick()
		if !ok {
			break
		}
		ate = append(ate, x)
	}
	l := tu.Link(c, user)
	s1 := format.Eaten(l, ate) + "\n\n"
	s2 := format.EnergyRemaining(user.Energy)
	return c.Send(s1+s2, tele.ModeHTML)
}
