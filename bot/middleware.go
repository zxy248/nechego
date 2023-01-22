package main

import (
	"log"
	"math/rand"
	"nechego/avatar"
	"nechego/game"
	"nechego/teleutil"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Wrapper interface {
	Wrap(tele.HandlerFunc) tele.HandlerFunc
}

type MessageIncrementer struct {
	Universe *game.Universe
}

func (m *MessageIncrementer) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		world, user := teleutil.Lock(c, m.Universe)
		world.Messages++
		user.Messages++
		world.Unlock()
		return next(c)
	}
}

type RequireSupergroup struct{}

func (m *RequireSupergroup) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if c.Chat().Type != tele.ChatSuperGroup {
			return nil
		}
		return next(c)
	}
}

type WrapperFunc tele.MiddlewareFunc

func (f WrapperFunc) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return f(next)
}

type RandomPhoto struct {
	Avatars *avatar.Storage
}

func (m *RandomPhoto) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if rand.Float64() < 0.02 {
			r := make([]*tele.Photo, 0, 2)
			p, err := c.Bot().ProfilePhotosOf(c.Sender())
			if err != nil {
				return err
			}
			if len(p) > 0 {
				r = append(r, &p[0])
			}
			if a, ok := m.Avatars.Get(c.Sender().ID); ok {
				r = append(r, a)
			}
			if len(r) > 0 {
				c.Send(r[rand.Intn(len(r))])
			}
		}
		return next(c)
	}
}

type IgnoreBanned struct {
	Universe *game.Universe
}

func (m *IgnoreBanned) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		world, user := teleutil.Lock(c, m.Universe)
		if user.Banned {
			world.Unlock()
			return nil
		}
		world.Unlock()
		return next(c)
	}
}

type IgnoreForwarded struct{}

func (m *IgnoreForwarded) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if c.Message().OriginalUnixtime != 0 {
			return nil
		}
		return next(c)
	}
}

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

type DeleterContext struct {
	tele.Context
}

func (c DeleterContext) Send(what interface{}, opts ...interface{}) error {
	msg, err := c.Bot().Send(c.Recipient(), what, opts...)
	if err != nil {
		return err
	}
	deleteLater(c.Bot(), msg)
	return nil
}

type DeleteMessage struct{}

func (m *DeleteMessage) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		deleteLater(c.Bot(), c.Message())
		return next(DeleterContext{c})
	}
}

func deleteLater(b *tele.Bot, m *tele.Message) {
	time.AfterFunc(5*time.Minute, func() { b.Delete(m) })
}
