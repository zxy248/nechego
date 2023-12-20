package main

import tele "gopkg.in/telebot.v3"

type Wrapper interface {
	Wrap(next tele.HandlerFunc) tele.HandlerFunc
}

type Wrapped struct {
	Service
	handler func(c tele.Context) error
}

func (h *Wrapped) Handle(c tele.Context) error {
	return h.handler(c)
}

func Wrap(s Service, ws []Wrapper) Service {
	r := &Wrapped{s, s.Handle}
	for _, w := range ws {
		r.handler = w.Wrap(r.handler)
	}
	return r
}
