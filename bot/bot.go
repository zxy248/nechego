package main

import (
	"log"
	"nechego/game"
	"nechego/handlers"
	"os"
	"os/signal"
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
		&handlers.Pic{Path: "data/pic"},
		&handlers.Basili{Path: "data/basili"},
		&handlers.Casper{Path: "data/casper"},
		&handlers.Zeus{Path: "data/zeus"},
		&handlers.Mouse{Path: "data/mouse.mp4"},
		&handlers.Tiktok{Path: "data/tiktok/"},
		&handlers.Hello{Path: "data/hello.json"},
		&handlers.Game{},
		&handlers.Infa{},
		&handlers.Weather{},
		&handlers.Cat{},
		&handlers.Anime{},
		&handlers.Furry{},
		&handlers.Flag{},
		&handlers.Person{},
		&handlers.Horse{},
		&handlers.Art{},
		&handlers.Car{},
		&handlers.Masyunya{},
		&handlers.Poppy{},
		&handlers.Sima{},
		&handlers.Calculator{},
		&handlers.Name{},
		&handlers.Who{Universe: universe},
		&handlers.Top{Universe: universe},
		&handlers.List{Universe: universe},
		&handlers.Save{Universe: universe},
		&handlers.DailyEblan{Universe: universe},
		&handlers.DailyAdmin{Universe: universe},
		&handlers.DailyPair{Universe: universe},
		&handlers.Inventory{Universe: universe},
		&handlers.Drop{Universe: universe},
		&handlers.Pick{Universe: universe},
		&handlers.Floor{Universe: universe},
		&handlers.Market{Universe: universe},
		&handlers.Buy{Universe: universe},
		&handlers.Eat{Universe: universe},
		&handlers.Fish{Universe: universe},
		&handlers.Status{Universe: universe},
		&handlers.Sell{Universe: universe},
		&handlers.Stack{Universe: universe},
	}
	router.Middleware = []Wrapper{
		&MessageIncrementer{Universe: universe},
		&UserAdder{Universe: universe},
		&RequireSupergroup{},
	}
	go func() {
		for range time.NewTicker(time.Second * 30).C {
			universe.ForEachWorld(func(w *game.World) {
				w.RestoreEnergy()
				w.Market.Refill()
			})
		}
	}()

	interrupt := make(chan os.Signal, 1)
	stop := make(chan struct{}, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		log.Println("Stopping the bot...")
		bot.Stop()
		log.Println("Saving the universe...")
		if err := universe.SaveAll(); err != nil {
			log.Fatal(err)
		}
		stop <- struct{}{}
	}()
	bot.Handle(tele.OnText, router.OnText)
	bot.Start()
	<-stop
	log.Println("Successful shutdown.")
}
