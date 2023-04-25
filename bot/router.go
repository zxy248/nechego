package main

import tele "gopkg.in/telebot.v3"

type Handler interface {
	Handle(tele.Context) error
}

type ContextService interface {
	Handler
	Match(tele.Context) bool
}

type StringService interface {
	Handler
	Match(string) bool
}

type Router struct {
	Handlers []ContextService
}

func (r *Router) Dispatch(c tele.Context) error {
	for _, h := range r.Handlers {
		if h.Match(c) {
			return h.Handle(c)
		}
	}
	return nil
}
