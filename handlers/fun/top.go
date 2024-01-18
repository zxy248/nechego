package fun

import (
	"fmt"
	"math/rand"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Top struct {
	Universe *game.Universe
}

var topRe = handlers.NewRegexp("^!топ ?(.*)")

func (h *Top) Match(c tele.Context) bool {
	return topRe.MatchString(c.Text())
}

func (h *Top) Handle(c tele.Context) error {
	world := tu.Lock(c, h.Universe)
	defer world.Unlock()

	name := topName(c.Text())
	n := 3 + rand.Intn(3)
	list := world.RandomUsers(n)
	out := fmt.Sprintf("<b>🏆 Топ %s</b>\n", name)
	for i, id := range list {
		l := tu.Link(c, id)
		out += fmt.Sprintf("%d. <b>%s</b>\n", i+1, l)
	}
	return c.Send(out, tele.ModeHTML)
}

func topName(s string) string {
	return topRe.FindStringSubmatch(s)[1]
}
