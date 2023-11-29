package economy

import (
	"fmt"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Stack struct {
	Universe *game.Universe
}

var stackRe = handlers.Regexp("^!сложить")

func (h *Stack) Match(c tele.Context) bool {
	return stackRe.MatchString(c.Text())
}

func (h *Stack) Handle(c tele.Context) error {
	world, user := tu.Lock(c, h.Universe)
	defer world.Unlock()

	user.Inventory.Stack()
	l := tu.Link(c, user)
	s := fmt.Sprintf("🗄 <b>%s</b> складывает вещи.", l)
	return c.Send(s, tele.ModeHTML)
}
