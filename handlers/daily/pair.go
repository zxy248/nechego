package daily

import (
	"fmt"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Pair struct {
	Universe *game.Universe
}

var pairRe = handlers.NewRegexp("^!пара")

func (h *Pair) Match(c tele.Context) bool {
	return pairRe.MatchString(c.Text())
}

func (h *Pair) Handle(c tele.Context) error {
	world := tu.Lock(c, h.Universe)
	defer world.Unlock()

	p := world.DailyPair()
	l1 := tu.Link(c, p[0])
	l2 := tu.Link(c, p[1])
	s := fmt.Sprintf("<b>✨ Пара дня</b> — %s 💘 %s", l1, l2)
	return c.Send(s, tele.ModeHTML)
}
