package middleware

import (
	"time"

	tele "gopkg.in/telebot.v3"
)

type AutoDelete struct {
	After time.Duration
}

func (m *AutoDelete) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		return next(autoDeleteContext{c, m.After})
	}
}

type autoDeleteContext struct {
	tele.Context
	after time.Duration
}

func (c autoDeleteContext) Send(what interface{}, opts ...interface{}) error {
	userMsg := c.Message()
	botMsg, err := c.Bot().Send(c.Recipient(), what, opts...)
	if err != nil {
		return err
	}
	c.deleteMessageAfter(userMsg, c.after)
	c.deleteMessageAfter(botMsg, c.after)
	return nil
}

func (c autoDeleteContext) deleteMessageAfter(m *tele.Message, d time.Duration) {
	time.AfterFunc(d, func() { c.Bot().Delete(m) })
}
