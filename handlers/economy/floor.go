package economy

import (
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Floor struct {
	Universe *game.Universe
}

var floorRe = handlers.Regexp("^!(пол|мусор|вещи|предметы)")

func (h *Floor) Match(c tele.Context) bool {
	return floorRe.MatchString(c.Text())
}

func (h *Floor) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	head := "<b>🗃️ Предметы</b>\n"
	list := format.Items(world.Floor.HkList())
	return c.Send(head+list, tele.ModeHTML)
}
