package server

import (
	"os"
	"os/signal"
	"syscall"

	tele "gopkg.in/telebot.v3"
)

type Service interface {
	Match(c tele.Context) bool
	Handle(c tele.Context) error
}

type Server struct {
	Bot      *tele.Bot
	Handlers []Service
}

func (s *Server) Run() {
	endpoints := []string{
		tele.OnText,
		tele.OnPhoto,
		tele.OnDice,
	}
	h := dispatcher(s.Handlers)
	for _, e := range endpoints {
		s.Bot.Handle(e, h)
	}

	x := shutdown(s.Bot)
	s.Bot.Start()
	<-x
}

func dispatcher(ss []Service) tele.HandlerFunc {
	return func(c tele.Context) error {
		for _, s := range ss {
			if s.Match(c) {
				return s.Handle(c)
			}
		}
		return nil
	}
}

func shutdown(b *tele.Bot) <-chan struct{} {
	x := make(chan struct{})
	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-interrupt
		b.Stop()
		x <- struct{}{}
	}()
	return x
}
