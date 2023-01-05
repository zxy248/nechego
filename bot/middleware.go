package main

import (
	"errors"
	"nechego/game"

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
			w.AddUser(game.NewUser(user.ID))
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
		u.IncrementMessages()
		w.Unlock()
		return next(c)
	}
}
