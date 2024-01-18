package middleware

import (
	tu "nechego/teleutil"

	tele "gopkg.in/zxy248/telebot.v3"
)

type RequireSupergroup struct{}

func (m *RequireSupergroup) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if !tu.SuperGroup(c) {
			return nil
		}
		return next(c)
	}
}
