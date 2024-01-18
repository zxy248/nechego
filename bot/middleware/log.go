package middleware

import (
	"fmt"
	"log"
	"strings"
	"time"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Log struct {
	Timeout time.Duration
}

func (m *Log) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		summary := contextSummary(c)
		timeout := time.After(m.Timeout)
		start := time.Now()
		errc := make(chan error, 1)
		go func() { errc <- next(c) }()

		select {
		case err := <-errc:
			d := time.Since(start)
			log.Printf("%s %s", d, summary)
			return err
		case <-timeout:
			log.Printf("∞ %s", summary)
			return nil
		}
	}
}

func contextSummary(c tele.Context) string {
	title := c.Chat().Title
	user := userName(c.Sender())
	text := c.Text()
	if m := c.Message(); m != nil && m.Sticker != nil {
		text = m.Sticker.Emoji
	}
	return fmt.Sprintf("|%s| %s -> %s\n", title, user, text)
}

func userName(u *tele.User) string {
	return strings.TrimSpace(u.FirstName + " " + u.LastName)
}
