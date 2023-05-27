package adapter

import (
	"nechego/bot/server"
	"strings"

	tele "gopkg.in/telebot.v3"
)

type TextService interface {
	Match(string) bool
	Handle(tele.Context) error
}

type Text struct {
	TextService
}

func (s *Text) Match(c tele.Context) bool {
	return s.TextService.Match(c.Text())
}

type Callback struct {
	TextService
}

func (s *Callback) Match(c tele.Context) bool {
	cb := c.Callback()
	if cb == nil {
		return false
	}
	cb.Data = strings.TrimSpace(cb.Data)
	return s.TextService.Match(cb.Data)
}

type Closure struct {
	server.Service
	handle func(tele.Context) error
}

func (s *Closure) Handle(c tele.Context) error {
	return s.handle(c)
}

type Wrapper interface {
	Wrap(tele.HandlerFunc) tele.HandlerFunc
}

func Wrap(s server.Service, w ...Wrapper) *Closure {
	h := &Closure{s, s.Handle}
	for i := len(w) - 1; i >= 0; i-- {
		h.handle = w[i].Wrap(h.handle)
	}
	return h
}
