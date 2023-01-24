package middleware

import (
	"log"
	"runtime/debug"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

var Recover = Tele(middleware.Recover(func(err error) {
	log.Print(err)
	debug.PrintStack()
}))

type LogMessage struct{}

func (m *LogMessage) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		start := time.Now()
		err := next(c)
		log.Printf("(%s) [%s] %s: %s\n",
			time.Since(start),
			c.Chat().Title,
			strings.TrimSpace(c.Sender().FirstName+" "+c.Sender().LastName),
			c.Text())
		return err
	}
}
