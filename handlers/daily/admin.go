package daily

import (
	"fmt"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele"gopkg.in/telebot.v3"
)

type Admin struct {
	Universe *game.Universe
}

var adminRe = handlers.Regexp("^!админ")

func (h *Admin) Match(c tele.Context) bool {
	return adminRe.MatchString(c.Text())
}

func (h *Admin) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	u, ok := world.DailyAdmin()
	if !ok {
		return c.Send("👑")
	}
	l := tu.Link(c, u)
	s := fmt.Sprintf("<b>Админ дня</b> — %s 👑", l)
	return c.Send(s, tele.ModeHTML)
}
