package middleware

import (
	"nechego/game"
	"nechego/teleutil"
	"time"

	tele "gopkg.in/telebot.v3"
)

type IgnoreBanned struct {
	Universe *game.Universe
}

func (m *IgnoreBanned) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		world, user := teleutil.Lock(c, m.Universe)
		banned := user.BannedUntil
		world.Unlock()

		if time.Now().Before(banned) {
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

type IgnoreSpam struct {
	Universe *game.Universe
}

func (m *IgnoreSpam) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		world, user := teleutil.Lock(c, m.Universe)
		lastMessage := user.LastMessage
		world.Unlock()

		if time.Since(lastMessage) < time.Second {
			return nil
		}
		return next(c)
	}
}
