package main

import (
	"log"
	"nechego/avatar"
	"nechego/bot/middleware"
	"nechego/game"
	"runtime/debug"
	"time"

	tele "gopkg.in/telebot.v3"
	telemw "gopkg.in/telebot.v3/middleware"
)

type Wrapper interface {
	Wrap(tele.HandlerFunc) tele.HandlerFunc
}

type TeleWrapper tele.MiddlewareFunc

func (f TeleWrapper) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return f(next)
}

func globalMiddleware(u *game.Universe, as *avatar.Storage) []Wrapper {
	return []Wrapper{
		TeleWrapper(telemw.Recover(func(err error) {
			log.Print(err)
			debug.PrintStack()
		})),
		&middleware.RequireSupergroup{},
		&middleware.IgnoreForwarded{},
		&middleware.IgnoreBanned{Universe: u},
		&middleware.LogMessage{},
		&middleware.Throttle{Duration: 800 * time.Millisecond},
		&middleware.IncrementCounters{Universe: u},
		&middleware.RandomPhoto{Avatars: as, Prob: 1. / 200},
	}
}
