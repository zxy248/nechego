package main

import (
	"log"
	"nechego/avatar"
	"nechego/game"
	"nechego/handlers"
	"os"
	"runtime/debug"
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
	avatars := &avatar.Storage{Dir: "avatar", MaxWidth: 1500, MaxHeight: 1500, Bot: bot}
	hello := &handlers.Hello{Path: "data/hello.json"}
	roll := &handlers.Roll{Universe: universe}
	router := &Router{}
	router.Handlers = []handlers.Handler{
		&handlers.Pic{Path: "data/pic"},
		&handlers.Basili{Path: "data/basili"},
		&handlers.Casper{Path: "data/casper"},
		&handlers.Zeus{Path: "data/zeus"},
		&handlers.Mouse{Path: "data/mouse.mp4"},
		&handlers.Tiktok{Path: "data/tiktok/"},
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
		&handlers.Soy{},
		&handlers.Danbooru{},
		&handlers.Fap{},
		&handlers.Masyunya{},
		&handlers.Poppy{},
		&handlers.Sima{},
		&handlers.Calculator{},
		&handlers.Name{},
		&handlers.Who{Universe: universe},
		&handlers.List{Universe: universe},
		&handlers.DailyEblan{Universe: universe},
		&handlers.DailyAdmin{Universe: universe},
		&handlers.DailyPair{Universe: universe},
		&handlers.Inventory{Universe: universe},
		&handlers.Sort{Universe: universe},
		&handlers.Catch{Universe: universe},
		&handlers.Drop{Universe: universe},
		&handlers.Pick{Universe: universe},
		&handlers.Floor{Universe: universe},
		&handlers.Market{Universe: universe},
		&handlers.Buy{Universe: universe},
		&handlers.Eat{Universe: universe},
		&handlers.EatQuick{Universe: universe},
		&handlers.Fish{Universe: universe},
		&handlers.Craft{Universe: universe},
		&handlers.Status{Universe: universe},
		&handlers.Sell{Universe: universe},
		&handlers.Stack{Universe: universe},
		&handlers.Cashout{Universe: universe},
		&handlers.Fight{Universe: universe},
		&handlers.Profile{Universe: universe, Avatars: avatars},
		&handlers.Avatar{Avatars: avatars},
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
		&handlers.Balance{Universe: universe},
		&handlers.Energy{Universe: universe},
		&handlers.NameMarket{Universe: universe},
		&handlers.NamePet{Universe: universe},
		&handlers.SendSMS{Universe: universe},
		&handlers.ReceiveSMS{Universe: universe},
		hello,
	}
	router.Middleware = []Wrapper{
		&RandomPhoto{Avatars: avatars},
		&MessageIncrementer{Universe: universe},
		&IgnoreBanned{Universe: universe},
		&DeleteMessage{},
		&LogMessage{},
		&IgnoreForwarded{},
		&RequireSupergroup{},
		WrapperFunc(middleware.Recover(func(err error) {
			log.Print(err)
			debug.PrintStack()
		})),
	}
	go refillMarket(universe)
	go restoreEnergy(universe)
	done := stopper(bot, universe)

	bot.Handle(tele.OnDice, router.wrap(roll.Handle))
	bot.Handle(tele.OnUserJoined, router.wrap(hello.Handle))
	bot.Handle(tele.OnText, router.OnText)
	bot.Handle(tele.OnPhoto, router.OnText)
	bot.Start()

	<-done
	log.Println("Successful shutdown.")
}
