package fun

import (
	"math/rand"
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type List struct {
	Universe *game.Universe
}

var listRe = handlers.Regexp("^!список ?(.*)")

func (h *List) Match(c tele.Context) bool {
	return listRe.MatchString(c.Text())
}

func (h *List) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	var links []string
	ids := world.RandomUserIDs(3 + rand.Intn(3))
	for _, id := range ids {
		l := tu.Link(c, id)
		links = append(links, l)
	}
	title := listTitle(c.Text())
	s := format.List(title, links)
	return c.Send(s, tele.ModeHTML)
}

func listTitle(s string) string {
	return listRe.FindStringSubmatch(s)[1]
}
