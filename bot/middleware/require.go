package middleware

import tele "gopkg.in/telebot.v3"

type RequireSupergroup struct{}

func (m *RequireSupergroup) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if c.Chat().Type != tele.ChatSuperGroup {
			return nil
		}
		return next(c)
	}
}
