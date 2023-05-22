package middleware

import tele "gopkg.in/telebot.v3"

type Context func(tele.Context) tele.Context

func (ctx Context) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		c = ctx(c)
		return next(c)
	}
}
