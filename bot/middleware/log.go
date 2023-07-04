package middleware

import (
	"fmt"
	"log"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

type LogMessage struct {
	Wait time.Duration
}

func (m *LogMessage) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		var prefix string
		var err error

		start := time.Now()
		select {
		case err = <-runHandler(c, next):
			prefix = time.Since(start).String()
		case <-time.After(m.Wait):
			prefix = "TOO LONG"
		}

		log.Printf("[%s] %s", prefix, contextSummary(c))
		return err
	}
}

func runHandler(c tele.Context, f tele.HandlerFunc) <-chan error {
	errc := make(chan error, 1)
	go func() {
		errc <- f(c)
	}()
	return errc
}

func contextSummary(c tele.Context) string {
	return fmt.Sprintf("<%s> %s :: %s\n", c.Chat().Title, userName(c.Sender()), c.Text())
}

func userName(u *tele.User) string {
	return strings.TrimSpace(u.FirstName + " " + u.LastName)
}
