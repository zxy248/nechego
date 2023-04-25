package middleware

import (
	"log"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

type LogMessage struct{}

func (m *LogMessage) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		// TODO: force log if too long
		start := time.Now()
		err := next(c)
		log.Printf("%s %s: %s: %s\n",
			time.Since(start),
			c.Chat().Title,
			strings.TrimSpace(c.Sender().FirstName+" "+c.Sender().LastName),
			c.Text())
		return err
	}
}
