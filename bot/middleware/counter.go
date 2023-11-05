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
		w, u := tu.Lock(c, m.Universe)
		u.Messages++
		u.LastMessage = time.Now()
		w.Messages++
		w.Unlock()

		return next(c)
	}
}
