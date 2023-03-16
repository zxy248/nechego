package main

import (
	"log"
	"math/rand"
	"nechego/avatar"
	"nechego/bot/context"
	"nechego/bot/middleware"
	"nechego/game"
	"nechego/handlers"
	"os"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

var botToken = os.Getenv("NECHEGO_TOKEN")

// init seeds the random number generator and validates the global
// variables.
func init() {
	rand.Seed(time.Now().UnixNano())
	if botToken == "" {
		log.Fatal("$NECHEGO_TOKEN not set")
	}
}

func main() {
	bot, err := tele.NewBot(tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}
	as := &avatar.Storage{Dir: "avatar", MaxWidth: 1500, MaxHeight: 1500, Bot: bot}
	u := game.NewUniverse("universe", worldInitializer(bot))

	router := NewRouter()
	for _, m := range middlewareWrappers(u, as) {
		router.AddMidleware(m)
	}
	for _, h := range textHandlers(u, as) {
		router.AddHandler(tele.OnText, h)
		router.AddHandler(tele.OnPhoto, h)
	}
	router.AddHandler(tele.OnDice, &handlers.Roll{Universe: u})
	router.AddHandler(tele.OnCallback, &handlers.HarvestInline{Universe: u})
	router.AddHandler(tele.OnCallback, &handlers.AuctionBuy{Universe: u})
	router.Set(bot)

	done := notifyStop(bot, u)
	bot.Start()
	<-done
	log.Println("Successful shutdown.")
}

type Router struct {
	// Lists of handlers indexed by their corresponding endpoints.
	Handlers map[string][]handlers.Handler

	// Request wrappers.
	Middleware []middleware.Wrapper
}

// NewRouter returns a new Router that does nothing.
func NewRouter() *Router {
	return &Router{
		Handlers:   map[string][]handlers.Handler{},
		Middleware: []middleware.Wrapper{},
	}
}

// Set sets the handlers for each endpoint.
func (r *Router) Set(bot *tele.Bot) {
	for endpoint := range r.Handlers {
		bot.Handle(endpoint, r.HandlerFunc(endpoint))
	}
}

// AddMidleware registers the specified middleware.
func (r *Router) AddMidleware(m middleware.Wrapper) {
	r.Middleware = append([]middleware.Wrapper{m}, r.Middleware...)
}

// AddHandler registers the specified handler for the given endpoint.
func (r *Router) AddHandler(endpoint string, h handlers.Handler) {
	r.Handlers[endpoint] = append(r.Handlers[endpoint], h)
}

// HandlerFunc returns the HandlerFunc that telebot can handle.
func (r *Router) HandlerFunc(endpoint string) tele.HandlerFunc {
	// This function delegates the given context to the first
	// matching handler. If there is no match, the function is a
	// no-op.
	f := func(c tele.Context) error {
		s := c.Text()
		if endpoint == tele.OnCallback {
			callback := c.Callback()
			s = strings.TrimSpace(callback.Data)
			callback.Data = s
		}
		for _, h := range r.Handlers[endpoint] {
			if h.Match(s) {
				context.SetHandlerID(c, h.Self())
				return h.Handle(c)
			}
		}
		return nil
	}
	for _, w := range r.Middleware {
		f = w.Wrap(f)
	}
	return f
}

func textHandlers(u *game.Universe, as *avatar.Storage) []handlers.Handler {
	return []handlers.Handler{
		&handlers.Help{},

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
		&handlers.Lageona{},

		// Daily.
		&handlers.DailyEblan{Universe: u},
		&handlers.DailyAdmin{Universe: u},
		&handlers.DailyPair{Universe: u},

		// Economy.
		&handlers.Inventory{Universe: u},
		&handlers.Funds{Universe: u},
		&handlers.Sort{Universe: u},
		&handlers.Drop{Universe: u},
		&handlers.Pick{Universe: u},
		&handlers.Floor{Universe: u},
		&handlers.Stack{Universe: u},
		&handlers.Split{Universe: u},
		&handlers.Cashout{Universe: u},
		&handlers.Capital{Universe: u},
		&handlers.Balance{Universe: u},

		// Farm.
		&handlers.Farm{Universe: u},
		&handlers.Plant{Universe: u},
		&handlers.Harvest{Universe: u},
		&handlers.PriceList{Universe: u},
		&handlers.UpgradeFarm{Universe: u},

		// Market.
		&handlers.Market{Universe: u},
		&handlers.Buy{Universe: u},
		&handlers.Sell{Universe: u},
		&handlers.SellQuick{Universe: u},
		&handlers.NameMarket{Universe: u},
		&handlers.GetJob{Universe: u},
		&handlers.QuitJob{Universe: u},

		// Auction.
		&handlers.Auction{Universe: u},
		&handlers.AuctionSell{Universe: u},

		// Actions.
		&handlers.Craft{Universe: u},
		&handlers.Fish{Universe: u},
		&handlers.DrawNet{Universe: u},
		&handlers.CastNet{Universe: u},
		&handlers.Net{Universe: u},
		&handlers.Catch{Universe: u},
		&handlers.Dice{Universe: u},
		&handlers.Fight{Universe: u},
		&handlers.PvP{Universe: u},
		&handlers.Eat{Universe: u},
		&handlers.EatQuick{Universe: u},
		&handlers.FishingRecords{Universe: u},
		&handlers.Friends{Universe: u},
		&handlers.Transfer{Universe: u},

		// Top.
		&handlers.TopStrong{Universe: u},
		&handlers.TopRating{Universe: u},
		&handlers.TopRich{Universe: u},

		// Profile.
		&handlers.Status{Universe: u, MaxLength: 120},
		&handlers.Profile{Universe: u, Avatars: as},
		&handlers.Avatar{Universe: u, Avatars: as},
		&handlers.Energy{Universe: u},
		&handlers.NamePet{Universe: u},

		// Phone.
		&handlers.SendSMS{Universe: u},
		&handlers.ReceiveSMS{Universe: u},
		&handlers.Contacts{Universe: u},
		&handlers.Spam{Universe: u},

		// Fun.
		&handlers.Hello{Path: "data/hello.json"},
		&handlers.Game{},
		&handlers.Infa{},
		&handlers.Weather{},
		&handlers.Calculator{},
		&handlers.Name{},
		&handlers.Who{Universe: u},
		&handlers.List{Universe: u},
		&handlers.TurnOn{Universe: u},
		&handlers.TurnOff{Universe: u},
		&handlers.Top{Universe: u},
	}
}

func middlewareWrappers(u *game.Universe, as *avatar.Storage) []middleware.Wrapper {
	return []middleware.Wrapper{
		middleware.Recover,
		&middleware.RequireSupergroup{},
		&middleware.IgnoreForwarded{},
		&middleware.LogMessage{},
		&middleware.DeleteMessage{},
		&middleware.IgnoreBanned{Universe: u},
		&middleware.IgnoreSpam{Universe: u},
		&middleware.IncrementCounters{Universe: u},
		&middleware.RandomPhoto{Avatars: as},
	}
}
