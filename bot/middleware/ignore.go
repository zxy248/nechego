package middleware

import (
	"nechego/game"
	"nechego/teleutil"
	"time"

	tele "gopkg.in/telebot.v3"
)

type IgnoreUserBanned struct {
	Universe *game.Universe
}

func (m *IgnoreUserBanned) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
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

type IgnoreMessageForwarded struct{}

func (m *IgnoreMessageForwarded) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if c.Message().OriginalUnixtime != 0 {
			return nil
		}
		return next(c)
	}
}

type IgnoreWorldInactive struct {
	Universe *game.Universe
	Immune   func(tele.Context) bool
}

func (m *IgnoreWorldInactive) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		world, _ := teleutil.Lock(c, m.Universe)
		inactive := world.Inactive
		world.Unlock()

		if inactive && !m.Immune(c) {
			return nil
		}
		return next(c)
	}
}
