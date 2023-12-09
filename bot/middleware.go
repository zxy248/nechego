package main

import tele "gopkg.in/telebot.v3"

type Wrapper interface {
	Wrap(next tele.HandlerFunc) tele.HandlerFunc
}

type Wrapped struct {
	Service
	Middleware []Wrapper
}

func (w *Wrapped) Handle(c tele.Context) error {
	h := w.Service.Handle
	for i := len(w.Middleware) - 1; i >= 0; i-- {
		h = w.Middleware[i].Wrap(h)
	}
	return h(c)
}
