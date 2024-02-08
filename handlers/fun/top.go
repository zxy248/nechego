package fun

import (
	"fmt"
	"math/rand/v2"

	"github.com/zxy248/nechego/game"
	"github.com/zxy248/nechego/handlers"
	tu "github.com/zxy248/nechego/teleutil"

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
	list := world.RandomUsers(3 + rand.N(3))
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
