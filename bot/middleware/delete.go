package middleware

import (
	"time"

	tele "gopkg.in/zxy248/telebot.v3"
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
	deleteAfter time.Duration
}

func (c autoDeleteContext) Send(what interface{}, opts ...interface{}) error {
	userMessage := c.Message()
	botMessage, err := c.Bot().Send(c.Recipient(), what, opts...)
	if err != nil {
		return err
	}
	c.deleteMessageAfter(userMessage, c.deleteAfter)
	c.deleteMessageAfter(botMessage, c.deleteAfter)
	return nil
}

func (c autoDeleteContext) deleteMessageAfter(m *tele.Message, d time.Duration) {
	time.AfterFunc(d, func() { c.Bot().Delete(m) })
}
