package middleware

import (
	"nechego/game"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type CacheName struct {
	Universe *game.Universe
}

func (m *CacheName) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		defer m.updateName(c)
		return next(c)
	}
}

func (m *CacheName) updateName(c tele.Context) {
	go func() {
		n := tu.Name(tu.Member(c, c.Sender()))
		w, u := tu.Lock(c, m.Universe)
		u.Name = n
		w.Unlock()
	}()
}
