package middleware

import (
	"nechego/game"
	tu "nechego/teleutil"
	"time"

	tele "gopkg.in/telebot.v3"
)

type IncrementCounters struct {
	Universe *game.Universe
}

func (m *IncrementCounters) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		tu.ContextWorld(c, m.Universe, func(w *game.World) {
			u := tu.CurrentUser(c, w)
			u.Messages++
			u.LastMessage = time.Now()
			w.Messages++
		})
		return next(c)
	}
}
