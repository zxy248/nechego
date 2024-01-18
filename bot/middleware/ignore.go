package middleware

import (
	"nechego/game"
	tu "nechego/teleutil"

	tele "gopkg.in/zxy248/telebot.v3"
)

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
		w := tu.Lock(c, m.Universe)
		off := w.Inactive
		w.Unlock()

		if off && !m.Immune(c) {
			return nil
		}
		return next(c)
	}
}
