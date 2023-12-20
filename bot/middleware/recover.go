package middleware

import (
	"log"
	"runtime/debug"

	tele "gopkg.in/telebot.v3"
)

type Recover struct{}

func (r *Recover) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		defer func() {
			if err := recover(); err != nil {
				log.Print(err)
				debug.PrintStack()
			}
		}()
		return next(c)
	}
}
