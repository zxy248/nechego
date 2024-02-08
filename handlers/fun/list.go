package fun

import (
	"fmt"
	"math/rand/v2"

	"github.com/zxy248/nechego/game"
	"github.com/zxy248/nechego/handlers"
	tu "github.com/zxy248/nechego/teleutil"

	tele "gopkg.in/zxy248/telebot.v3"
)

type List struct {
	Universe *game.Universe
}

var listRe = handlers.NewRegexp("^!список ?(.*)")

func (h *List) Match(c tele.Context) bool {
	return listRe.MatchString(c.Text())
}

func (h *List) Handle(c tele.Context) error {
	world := tu.Lock(c, h.Universe)
	defer world.Unlock()

	name := listName(c.Text())
	list := world.RandomUsers(3 + rand.N(3))
	out := fmt.Sprintf("<b>📝 Список %s</b>\n", name)
	for _, id := range list {
		l := tu.Link(c, id)
		out += fmt.Sprintf("• <b>%s</b>\n", l)
	}
	return c.Send(out, tele.ModeHTML)
}

func listName(s string) string {
	return listRe.FindStringSubmatch(s)[1]
}
