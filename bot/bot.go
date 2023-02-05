package main

import (
	"fmt"
	"log"
	"math/rand"
	"nechego/avatar"
	"nechego/bot/context"
	"nechego/bot/middleware"
	"nechego/game"
	"nechego/handlers"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Router struct {
	Handlers   []handlers.Handler
	Special    map[handlers.HandlerID]handlers.Handler
	Middleware []middleware.Wrapper
}

func (r *Router) HandlerFunc(endpoint string) tele.HandlerFunc {
	var f tele.HandlerFunc
	switch endpoint {
	case tele.OnText, tele.OnPhoto:
		special := []handlers.Handler{}
		for _, h := range r.Special {
			special = append(special, h)
		}
		all := append(r.Handlers, special...)

		f = func(c tele.Context) error {
			for _, h := range all {
				if h.Match(c.Text()) {
					context.SetHandlerID(c, h.Self())
					return h.Handle(c)
				}
			}
			return nil

		}
	case tele.OnDice:
		f = func(c tele.Context) error {
			h := handlers.RollHandler
			context.SetHandlerID(c, h)
			return r.Special[h].Handle(c)
		}
	case tele.OnUserJoined:
		f = func(c tele.Context) error {
			h := handlers.HelloHandler
			context.SetHandlerID(c, h)
			return r.Special[h].Handle(c)
		}
	default:
		panic(fmt.Sprintf("unexpected endpoint %s", endpoint))
	}
	for _, w := range r.Middleware {
		f = w.Wrap(f)
	}
	return f
}

var token = os.Getenv("NECHEGO_TOKEN")

func init() {
	rand.Seed(time.Now().UnixNano())

	if token == "" {
		log.Fatal("$NECHEGO_TOKEN not set")
	}
}

func main() {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	universe := game.NewUniverse("universe", worldInitializer(bot))
	avatars := &avatar.Storage{Dir: "avatar", MaxWidth: 1500, MaxHeight: 1500, Bot: bot}
	router := &Router{}
	router.Handlers = []handlers.Handler{
		// Pictures.
		&handlers.Pic{Path: "data/pic"},
		&handlers.Basili{Path: "data/basili"},
		&handlers.Casper{Path: "data/casper"},
		&handlers.Zeus{Path: "data/zeus"},
		&handlers.Mouse{Path: "data/mouse.mp4"},
		&handlers.Tiktok{Path: "data/tiktok/"},
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

		// Daily.
		&handlers.DailyEblan{Universe: universe},
		&handlers.DailyAdmin{Universe: universe},
		&handlers.DailyPair{Universe: universe},

		// Economy.
		&handlers.Inventory{Universe: universe},
		&handlers.Sort{Universe: universe},
		&handlers.Drop{Universe: universe},
		&handlers.Pick{Universe: universe},
		&handlers.Floor{Universe: universe},
		&handlers.Stack{Universe: universe},
		&handlers.Cashout{Universe: universe},
		&handlers.Capital{Universe: universe},
		&handlers.Balance{Universe: universe},

		// Market.
		&handlers.Market{Universe: universe},
		&handlers.Buy{Universe: universe},
		&handlers.Sell{Universe: universe},
		&handlers.NameMarket{Universe: universe},

		// Actions.
		&handlers.Craft{Universe: universe},
		&handlers.Fish{Universe: universe},
		&handlers.DrawNet{Universe: universe},
		&handlers.CastNet{Universe: universe},
		&handlers.Net{Universe: universe},
		&handlers.Catch{Universe: universe},
		&handlers.Dice{Universe: universe},
		&handlers.Fight{Universe: universe},
		&handlers.Eat{Universe: universe},
		&handlers.EatQuick{Universe: universe},
		&handlers.FishingRecords{Universe: universe},

		// Admin.
		&handlers.Ban{Universe: universe, DurationHr: 2},
		&handlers.Unban{Universe: universe},

		// Top.
		&handlers.TopStrong{Universe: universe},
		&handlers.TopRating{Universe: universe},
		&handlers.TopRich{Universe: universe},

		// Profile.
		&handlers.Status{Universe: universe},
		&handlers.Profile{Universe: universe, Avatars: avatars},
		&handlers.Avatar{Avatars: avatars},
		&handlers.Energy{Universe: universe},
		&handlers.NamePet{Universe: universe},

		// Phone.
		&handlers.SendSMS{Universe: universe},
		&handlers.ReceiveSMS{Universe: universe},
		&handlers.Contacts{Universe: universe},
		&handlers.Spam{Universe: universe},

		// Fun.
		&handlers.Game{},
		&handlers.Infa{},
		&handlers.Weather{},
		&handlers.Calculator{},
		&handlers.Name{},
		&handlers.Who{Universe: universe},
		&handlers.List{Universe: universe},
		&handlers.TurnOn{Universe: universe},
		&handlers.TurnOff{Universe: universe},
		&handlers.Top{Universe: universe},
	}
	router.Special = map[handlers.HandlerID]handlers.Handler{
		handlers.HelloHandler: &handlers.Hello{Path: "data/hello.json"},
		handlers.RollHandler:  &handlers.Roll{Universe: universe},
	}
	router.Middleware = []middleware.Wrapper{
		&middleware.RandomPhoto{Avatars: avatars},
		&middleware.IncrementCounters{Universe: universe},
		&middleware.IgnoreBanned{Universe: universe},
		&middleware.DeleteMessage{},
		&middleware.LogMessage{},
		&middleware.IgnoreForwarded{},
		&middleware.RequireSupergroup{},
		middleware.Recover,
	}
	go refillMarket(universe)
	go restoreEnergy(universe)
	go fillNet(universe)
	done := stopper(bot, universe)

	endpoints := [...]string{
		tele.OnText,
		tele.OnDice,
		tele.OnUserJoined,
		tele.OnPhoto,
	}
	for _, e := range endpoints {
		bot.Handle(e, router.HandlerFunc(e))
	}
	bot.Start()

	<-done
	log.Println("Successful shutdown.")
}
