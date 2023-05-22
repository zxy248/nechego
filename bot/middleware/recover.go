package middleware

import (
	"log"
	"runtime/debug"

	tele "gopkg.in/telebot.v3"
	telemw "gopkg.in/telebot.v3/middleware"
)

type Recover struct{}

func (r *Recover) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	f := telemw.Recover(func(err error) {
		log.Print(err)
		debug.PrintStack()
	})
	return f(next)
}
