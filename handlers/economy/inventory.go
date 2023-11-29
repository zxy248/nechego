package economy

import (
	"fmt"
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Inventory struct {
	Universe *game.Universe
}

var inventoryRe = handlers.Regexp("^!(инвентарь|лут)")

func (h *Inventory) Match(c tele.Context) bool {
	return inventoryRe.MatchString(c.Text())
}

func (h *Inventory) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	l := tu.Link(c, user)
	head := fmt.Sprintf("<b>🗄 %s: Инвентарь</b>\n", l)
	items := user.Inventory.HkList()
	list := format.Items(items)
	return c.Send(head+list, tele.ModeHTML)
}
