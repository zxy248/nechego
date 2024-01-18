package middleware

import (
	"nechego/game"
	tu "nechego/teleutil"
	"slices"

	tele "gopkg.in/zxy248/telebot.v3"
)

type AddUser struct {
	Universe *game.Universe
}

func (m *AddUser) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		w := tu.Lock(c, m.Universe)
		if id := c.Sender().ID; !slices.Contains(w.Users, id) {
			w.Users = append(w.Users, id)
		}
		w.Unlock()

		return next(c)
	}
}
