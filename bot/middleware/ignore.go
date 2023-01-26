package middleware

import (
	"nechego/game"
	"nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type IgnoreBanned struct {
	Universe *game.Universe
}

func (m *IgnoreBanned) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		world, user := teleutil.Lock(c, m.Universe)
		banned := user.Banned
		world.Unlock()

		if banned {
			return nil
		}
		return next(c)
	}
}

type IgnoreForwarded struct{}

func (m *IgnoreForwarded) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if c.Message().OriginalUnixtime != 0 {
			return nil
		}
		return next(c)
	}
}