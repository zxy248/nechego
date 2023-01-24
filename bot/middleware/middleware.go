package middleware

import tele "gopkg.in/telebot.v3"

type Wrapper interface {
	Wrap(tele.HandlerFunc) tele.HandlerFunc
}

type Tele tele.MiddlewareFunc

func (f Tele) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return f(next)
}
