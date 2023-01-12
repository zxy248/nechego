package main

import (
	"log"
	"math/rand"
	"nechego/game"
	"nechego/teleutil"
	"strings"

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

type RandomPhoto struct{}

func (m *RandomPhoto) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		p, err := c.Bot().ProfilePhotosOf(c.Sender())
		if err != nil {
			return err
		}
		if len(p) > 0 && rand.Float64() < 0.02 {
			c.Send(&p[0])
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
		log.Printf("[%s] %s: %s\n",
			c.Chat().Title,
			strings.TrimSpace(c.Sender().FirstName+" "+c.Sender().LastName),
			c.Text())
		return next(c)
	}
}
