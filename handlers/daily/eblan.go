package daily

import (
	"fmt"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Eblan struct {
	Universe *game.Universe
}

var eblanRe = handlers.NewRegexp("^![ие][б6п]?л[ап]н[а-я]*")

func (h *Eblan) Match(c tele.Context) bool {
	return eblanRe.MatchString(c.Text())
}

func (h *Eblan) Handle(c tele.Context) error {
	world := tu.Lock(c, h.Universe)
	defer world.Unlock()

	u := world.DailyEblan()
	l := tu.Link(c, u)
	s := fmt.Sprintf("<b>Еблан дня</b> — %s 😸", l)
	return c.Send(s, tele.ModeHTML)
}
