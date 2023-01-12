package main

import (
	"log"
	"nechego/game"
	"nechego/handlers"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type Router struct {
	Handlers   []handlers.Handler
	Middleware []Wrapper
}

func (r *Router) wrap(f tele.HandlerFunc) tele.HandlerFunc {
	for _, w := range r.Middleware {
		f = w.Wrap(f)
	}
	return f
}

func (r *Router) OnText(c tele.Context) error {
	for _, h := range r.Handlers {
		if h.Match(c.Text()) {
			return r.wrap(h.Handle)(c)
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
	helloHandler := &handlers.Hello{Path: "data/hello.json"}
	router := &Router{}
	router.Handlers = []handlers.Handler{
		&handlers.Pic{Path: "data/pic"},
		&handlers.Basili{Path: "data/basili"},
		&handlers.Casper{Path: "data/casper"},
		&handlers.Zeus{Path: "data/zeus"},
		&handlers.Mouse{Path: "data/mouse.mp4"},
		&handlers.Tiktok{Path: "data/tiktok/"},
		helloHandler,
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
		&handlers.Fight{Universe: universe},
		&handlers.Profile{Universe: universe, AvatarPath: "data/avatar"},
		&handlers.Avatar{Path: "data/avatar"},
		&handlers.Dice{Universe: universe},
		&handlers.TurnOn{Universe: universe},
		&handlers.TurnOff{Universe: universe},
		&handlers.Ban{Universe: universe},
		&handlers.Unban{Universe: universe},
		&handlers.TopStrong{Universe: universe},
		&handlers.TopRating{Universe: universe},
		&handlers.TopRich{Universe: universe},
		&handlers.Top{Universe: universe},
		&handlers.Capital{Universe: universe},
	}
	router.Middleware = []Wrapper{
		&RandomPhoto{},
		&MessageIncrementer{Universe: universe},
		&IgnoreBanned{Universe: universe},
		&UserAdder{Universe: universe},
		&LogMessage{},
		&IgnoreForwarded{},
		&RequireSupergroup{},
		WrapperFunc(middleware.Recover(func(err error) {
			log.Print(err)
			debug.PrintStack()
		})),
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
	done := make(chan struct{}, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-interrupt
		log.Println("Stopping the bot...")
		bot.Stop()
		log.Println("Saving the universe...")
		if err := universe.SaveAll(); err != nil {
			log.Fatal(err)
		}
		done <- struct{}{}
	}()
	rollHandler := &handlers.Roll{Universe: universe}
	bot.Handle(tele.OnDice, router.wrap(rollHandler.Handle))
	bot.Handle(tele.OnUserJoined, router.wrap(helloHandler.Handle))
	bot.Handle(tele.OnText, router.OnText)
	bot.Handle(tele.OnPhoto, router.OnText)
	bot.Start()
	<-done
	log.Println("Successful shutdown.")
}
