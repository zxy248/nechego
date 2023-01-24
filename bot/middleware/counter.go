package middleware

import (
	"nechego/game"
	"nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type IncrementCounters struct {
	Universe *game.Universe
}

func (m *IncrementCounters) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		world, user := teleutil.Lock(c, m.Universe)
		world.Messages++
		user.Messages++
		world.Unlock()
		return next(c)
	}
}
