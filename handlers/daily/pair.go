package daily

import (
	"fmt"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele"gopkg.in/telebot.v3"
)

type Pair struct {
	Universe *game.Universe
}

var pairRe = handlers.Regexp("^!пара")

func (h *Pair) Match(c tele.Context) bool {
	return pairRe.MatchString(c.Text())
}

func (h *Pair) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	p, ok := world.DailyPair()
	if !ok {
		return c.Send("💔")
	}
	l1 := tu.Link(c, p[0])
	l2 := tu.Link(c, p[1])
	s := fmt.Sprintf("<b>✨ Пара дня</b> — %s 💘 %s", l1, l2)
	return c.Send(s, tele.ModeHTML)
}
