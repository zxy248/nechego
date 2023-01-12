package main

import (
	"errors"
	"log"
	"math/rand"
	"nechego/game"
	"strings"

	tele "gopkg.in/telebot.v3"
)

type Wrapper interface {
	Wrap(tele.HandlerFunc) tele.HandlerFunc
}

type UserAdder struct {
	Universe *game.Universe
}

func (m *UserAdder) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		user := c.Sender()
		w := m.Universe.MustWorld(c.Chat().ID)
		w.Lock()
		_, ok := w.UserByID(user.ID)
		if !ok {
			user := game.NewUser(user.ID)
			w.AddUser(user)
			user.Inventory.Add(&game.Item{
				Type:         game.ItemTypeWallet,
				Transferable: true,
				Value:        &game.Wallet{Money: 5000},
			})
		}
		w.Unlock()
		return next(c)
	}
}

type MessageIncrementer struct {
	Universe *game.Universe
}

func (m *MessageIncrementer) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		w := m.Universe.MustWorld(c.Chat().ID)
		w.Lock()
		u, ok := w.UserByID(c.Sender().ID)
		if !ok {
			w.Unlock()
			return errors.New("user not found")
		}
		w.Messages++
		u.Messages++
		w.Unlock()
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
		world := m.Universe.MustWorld(c.Chat().ID)
		world.Lock()
		user, ok := world.UserByID(c.Sender().ID)
		if !ok {
			world.Unlock()
			return errors.New("user not found")
		}
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
