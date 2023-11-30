package fun

import (
	"math/rand"
	"nechego/format"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type Top struct {
	Universe *game.Universe
}

var topRe = handlers.Regexp("^!топ ?(.*)")

func (h *Top) Match(c tele.Context) bool {
	return topRe.MatchString(c.Text())
}

func (h *Top) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	title := topTitle(c.Text())
	n := 3 + rand.Intn(3)
	us := randomUsers(world, n)
	s := format.TopPlain(title, us)
	return c.Send(s, tele.ModeHTML)
}

func topTitle(s string) string {
	return topRe.FindStringSubmatch(s)[1]
}

func randomUsers(w *game.World, n int) []*game.User {
	var us []*game.User
	ids := w.RandomUserIDs(n)
	for _, id := range ids {
		us = append(us, w.User(id))
	}
	return us
}
