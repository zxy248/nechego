package main

import (
	tele "gopkg.in/telebot.v3"
)

type Handler interface {
	Match(c tele.Context) bool
	Handle(c tele.Context) error
}

type Server struct {
	Bot      *tele.Bot
	Handlers []Handler
}

func (s *Server) Start() {
	eps := []string{tele.OnText, tele.OnPhoto}
	hf := handleFirstMatch(s.Handlers)
	for _, ep := range eps {
		s.Bot.Handle(ep, hf)
	}
	s.Bot.Start()
}

func (s *Server) Stop() {
	s.Bot.Stop()
}

func handleFirstMatch(hs []Handler) tele.HandlerFunc {
	return func(c tele.Context) error {
		for _, h := range hs {
			if h.Match(c) {
				return h.Handle(c)
			}
		}
		return nil
	}
}
