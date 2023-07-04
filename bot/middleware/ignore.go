package middleware

import (
	"nechego/game"
	tu "nechego/teleutil"
	"time"

	tele "gopkg.in/telebot.v3"
)

type IgnoreUserBanned struct {
	Universe *game.Universe
}

func (m *IgnoreUserBanned) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		var banExpires time.Time
		tu.ContextWorld(c, m.Universe, func(w *game.World) {
			banExpires = w.UserByID(c.Sender().ID).BannedUntil
		})

		if time.Now().Before(banExpires) {
			return nil
		}
		return next(c)
	}
}

type IgnoreMessageForwarded struct{}

func (m *IgnoreMessageForwarded) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if tu.MessageForwarded(c.Message()) {
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
		var inactive bool
		tu.ContextWorld(c, m.Universe, func(w *game.World) {
			inactive = w.Inactive
		})

		if inactive && !m.Immune(c) {
			return nil
		}
		return next(c)
	}
}
