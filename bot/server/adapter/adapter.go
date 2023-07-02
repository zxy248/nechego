package adapter

import (
	"nechego/bot/server"
	"strings"

	tele "gopkg.in/telebot.v3"
)

type TextService interface {
	Match(s string) bool
	Handle(c tele.Context) error
}

type Text struct{ TextService }

func (s *Text) Match(c tele.Context) bool {
	return s.TextService.Match(c.Text())
}

type Callback struct{ TextService }

func (s *Callback) Match(c tele.Context) bool {
	cb := c.Callback()
	if cb == nil {
		return false
	}
	cb.Data = strings.TrimSpace(cb.Data)
	return s.TextService.Match(cb.Data)
}

type Wrapped struct {
	server.Service
	handle func(tele.Context) error
}

func (s *Wrapped) Handle(c tele.Context) error {
	return s.handle(c)
}

type Wrapper interface {
	Wrap(next tele.HandlerFunc) tele.HandlerFunc
}

func Wrap(s server.Service, w Wrapper) *Wrapped {
	return &Wrapped{s, w.Wrap(s.Handle)}
}
