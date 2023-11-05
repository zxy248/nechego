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
		w, u := tu.Lock(c, m.Universe)
		e := u.BannedUntil
		w.Unlock()

		if time.Now().Before(e) {
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
		w, _ := tu.Lock(c, m.Universe)
		off := w.Inactive
		w.Unlock()

		if off && !m.Immune(c) {
			return nil
		}
		return next(c)
	}
}
