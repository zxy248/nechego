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
		start := time.Now()
		errc := make(chan error, 1)
		go func() {
			errc <- next(c)
		}()
		select {
		case err := <-errc:
			log.Printf("[%s] %s", time.Since(start), contextSummary(c))
			return err
		case <-time.After(m.Wait):
			log.Printf("[TOO LONG] %s", contextSummary(c))
			return nil
		}
	}
}

func contextSummary(c tele.Context) string {
	return fmt.Sprintf("<%s> %s :: %s\n", c.Chat().Title, userName(c.Sender()), c.Text())
}

func userName(u *tele.User) string {
	return strings.TrimSpace(u.FirstName + " " + u.LastName)
}
