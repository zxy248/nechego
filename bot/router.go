package main

import (
	"strings"

	tele "gopkg.in/telebot.v3"
)

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

type TextHandler struct {
	StringService
}

func (h *TextHandler) Match(c tele.Context) bool {
	return h.StringService.Match(c.Text())
}

type CallbackHandler struct {
	StringService
}

func (h *CallbackHandler) Match(c tele.Context) bool {
	cb := c.Callback()
	if cb == nil {
		return false
	}
	cb.Data = strings.TrimSpace(cb.Data)
	return h.StringService.Match(cb.Data)
}

type ClosureHandler struct {
	ContextService
	H func(tele.Context) error
}

func (h *ClosureHandler) Handle(c tele.Context) error {
	return h.H(c)
}

type Wrapper interface {
	Wrap(tele.HandlerFunc) tele.HandlerFunc
}

func wrap(s ContextService, w ...Wrapper) *ClosureHandler {
	h := &ClosureHandler{ContextService: s, H: s.Handle}
	for i := len(w) - 1; i >= 0; i-- {
		h.H = w[i].Wrap(h.H)
	}
	return h
}

type router []ContextService

func (r router) dispatch(c tele.Context) error {
	for _, h := range r {
		if h.Match(c) {
			return h.Handle(c)
		}
	}
	return nil
}

func dispatch(b *tele.Bot, h tele.HandlerFunc, endpoints ...any) {
	for _, e := range endpoints {
		b.Handle(e, h)
	}
}
