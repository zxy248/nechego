package main

import (
	"log"
	"nechego/game"
	"nechego/handlers"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Router struct {
	Handlers   []handlers.Handler
	Middleware []Wrapper
}

func (r *Router) OnText(c tele.Context) error {
	for _, h := range r.Handlers {
		if h.Match(c.Message().Text) {
			f := h.Handle
			for _, w := range r.Middleware {
				f = w.Wrap(f)
			}
			return f(c)
		}
	}
	return nil
}

func main() {
	pref := tele.Settings{
		Token:  os.Getenv("NECHEGO_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	universe := game.NewUniverse("universe")
	router := &Router{}
	router.Handlers = []handlers.Handler{
		&handlers.Mouse{Path: "data/mouse.mp4"},
		&handlers.Tiktok{Path: "data/tiktok/"},
		&handlers.Game{},
		&handlers.Infa{},
		&handlers.Who{Universe: universe},
		&handlers.Save{Universe: universe},
		&handlers.Weather{},
	}
	router.Middleware = []Wrapper{
		&MessageIncrementer{Universe: universe},
		&UserAdder{Universe: universe},
	}

	bot.Handle(tele.OnText, router.OnText)
	bot.Start()
}
