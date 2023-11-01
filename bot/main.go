package main

import (
	"fmt"
	"log"
	"nechego/avatar"
	"nechego/bot/middleware"
	"nechego/bot/server"
	"nechego/bot/server/adapter"
	"nechego/danbooru"
	"nechego/game"
	"nechego/handlers"
	"nechego/handlers/casino"
	"nechego/handlers/fun"
	"nechego/handlers/pictures"
	"nechego/services/router"
	"os"
	"path/filepath"
	"time"

	tele "gopkg.in/telebot.v3"
)

const (
	universeDirectory = "universe"
	avatarDirectory   = "avatar"
)

var (
	botToken        = getEnv("NECHEGO_TOKEN")
	assetsDirectory = getEnv("NECHEGO_ASSETS")
)

func main() {
	app, err := setup()
	if err != nil {
		log.Fatal("cannot setup: ", err)
	}
	srv := &server.Server{
		Bot:      app.bot,
		Handlers: app.services(),
	}
	srv.Run()
	if err := app.shutdown(); err != nil {
		log.Fatal("cannot shutdown: ", err)
	}
	log.Println("successful shutdown")
}

func assetPath(s string) string {
	return filepath.Join(assetsDirectory, s)
}

func getEnv(s string) string {
	e := os.Getenv(s)
	if e == "" {
		panic(fmt.Sprintf("%s not set", s))
	}
	return e
}

type app struct {
	bot      *tele.Bot
	universe *game.Universe
	avatars  *avatar.Storage
	danbooru *danbooru.Danbooru
}

func setup() (*app, error) {
	bot, err := tele.NewBot(tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}
	return &app{
		bot: bot,
		universe: game.NewUniverse(universeDirectory, func(w *game.World) {
			const want = 10
			for i := len(w.Market.Products()); i < want; i++ {
				w.Market.Refill()
			}
			w.History.Announce(handlers.RecordAnnouncer(bot, tele.ChatID(w.TGID)))
		}),
		avatars: &avatar.Storage{
			Bot:       bot,
			Dir:       avatarDirectory,
			MaxWidth:  1500,
			MaxHeight: 1500,
		},
		danbooru: danbooru.New(danbooru.URL, 5*time.Second, 3),
	}, nil
}

func (a *app) shutdown() error {
	return a.universe.SaveAll()
}

func (a *app) services() []server.Service {
	global := a.globalMiddleware()
	spam := []adapter.Wrapper{&middleware.AutoDelete{After: 5 * time.Minute}}

	groups := []struct {
		services   []server.Service
		middleware []adapter.Wrapper
	}{
		{a.informationServices(), nil},
		{a.dailyServices(), nil},
		{a.economyServices(), spam},
		{a.farmServices(), spam},
		{a.marketServices(), spam},
		{a.actionsServices(), spam},
		{a.topServices(), spam},
		{a.profileServices(), spam},
		{a.funServices(), nil},
		{a.pictureServices(), nil},
		{a.casinoServices(), spam},
		{a.nextServices(), nil},
	}

	var handlers []server.Service
	for _, g := range groups {
		for _, s := range g.services {
			h := wrap(s, concat(global, g.middleware)...)
			handlers = append(handlers, h)
		}
	}
	return handlers
}

func (a *app) globalMiddleware() []adapter.Wrapper {
	return []adapter.Wrapper{
		&middleware.Recover{},
		&middleware.RequireSupergroup{},
		&middleware.IgnoreMessageForwarded{},
		&middleware.IgnoreUserBanned{Universe: a.universe},
		&middleware.LogMessage{Wait: 5 * time.Second},
		&middleware.Throttle{Duration: 800 * time.Millisecond},
		&middleware.IgnoreWorldInactive{
			Universe: a.universe,
			Immune: func(c tele.Context) bool {
				return fun.MatchTurnOn(c.Message().Text)
			},
		},
		&middleware.IncrementCounters{Universe: a.universe},
		&middleware.RandomPhoto{Avatars: a.avatars, Prob: 1. / 200},
	}
}

func (a *app) informationServices() []server.Service {
	return []server.Service{
		text(&handlers.Help{}),
	}
}

func (a *app) dailyServices() []server.Service {
	return []server.Service{
		text(&handlers.DailyEblan{Universe: a.universe}),
		text(&handlers.DailyAdmin{Universe: a.universe}),
		text(&handlers.DailyPair{Universe: a.universe}),
	}
}

func (a *app) economyServices() []server.Service {
	return []server.Service{
		text(&handlers.Inventory{Universe: a.universe}),
		text(&handlers.Funds{Universe: a.universe}),
		text(&handlers.Sort{Universe: a.universe}),
		text(&handlers.Drop{Universe: a.universe}),
		text(&handlers.Pick{Universe: a.universe}),
		text(&handlers.Floor{Universe: a.universe}),
		text(&handlers.Stack{Universe: a.universe}),
		text(&handlers.Split{Universe: a.universe}),
		text(&handlers.Cashout{Universe: a.universe}),
		text(&handlers.Capital{Universe: a.universe}),
		text(&handlers.Balance{Universe: a.universe}),
	}
}

func (a *app) farmServices() []server.Service {
	return []server.Service{
		text(&handlers.Farm{Universe: a.universe}),
		text(&handlers.Plant{Universe: a.universe}),
		text(&handlers.Harvest{Universe: a.universe}),
		text(&handlers.PriceList{Universe: a.universe}),
		text(&handlers.UpgradeFarm{Universe: a.universe}),
		text(&handlers.NameFarm{Universe: a.universe}),
	}
}

func (a *app) marketServices() []server.Service {
	return []server.Service{
		text(&handlers.Market{Universe: a.universe}),
		text(&handlers.Buy{Universe: a.universe}),
		text(&handlers.Sell{Universe: a.universe}),
		text(&handlers.SellQuick{Universe: a.universe}),
		text(&handlers.NameMarket{Universe: a.universe}),
		text(&handlers.GetJob{Universe: a.universe}),
		text(&handlers.QuitJob{Universe: a.universe}),
	}
}

func (a *app) actionsServices() []server.Service {
	return []server.Service{
		text(&handlers.Craft{Universe: a.universe}),
		text(&handlers.Fish{Universe: a.universe}),
		text(&handlers.DrawNet{Universe: a.universe}),
		text(&handlers.CastNet{Universe: a.universe}),
		text(&handlers.Net{Universe: a.universe}),
		text(&handlers.Catch{Universe: a.universe}),
		text(&handlers.Fight{Universe: a.universe}),
		text(&handlers.Eat{Universe: a.universe}),
		text(&handlers.EatQuick{Universe: a.universe}),
		text(&handlers.FishingRecords{Universe: a.universe}),
		text(&handlers.Friends{Universe: a.universe}),
		text(&handlers.Transfer{Universe: a.universe}),
		text(&handlers.Use{Universe: a.universe}),
	}
}

func (a *app) topServices() []server.Service {
	return []server.Service{
		text(&handlers.TopStrong{Universe: a.universe}),
		text(&handlers.TopRating{Universe: a.universe}),
		text(&handlers.TopRich{Universe: a.universe}),
	}
}

func (a *app) profileServices() []server.Service {
	return []server.Service{
		text(&handlers.Status{Universe: a.universe, MaxLength: 120}),
		text(&handlers.Profile{Universe: a.universe, Avatars: a.avatars}),
		text(&handlers.Avatar{Universe: a.universe, Avatars: a.avatars}),
		text(&handlers.Energy{Universe: a.universe}),
		text(&handlers.NamePet{Universe: a.universe}),
	}
}

func (a *app) funServices() []server.Service {
	return []server.Service{
		text(&handlers.Game{}),
		text(&handlers.Infa{}),
		text(&handlers.Weather{}),
		text(&handlers.Calculator{}),
		text(&fun.Name{}),
		text(&fun.CheckName{}),
		text(&handlers.Who{Universe: a.universe}),
		text(&handlers.List{Universe: a.universe}),
		text(&handlers.Top{Universe: a.universe}),
		text(&fun.Clock{}),
		text(&fun.TurnOn{Universe: a.universe}),
		text(&fun.TurnOff{Universe: a.universe}),
		text(&fun.Reputation{Universe: a.universe}),
		text(&fun.UpdateReputation{Universe: a.universe}),
	}
}

func (a *app) pictureServices() []server.Service {
	return []server.Service{
		text(&pictures.Pic{Path: assetPath("pic")}),
		text(&pictures.Basili{Path: assetPath("basili")}),
		text(&pictures.Casper{Path: assetPath("casper")}),
		text(&pictures.Zeus{Path: assetPath("zeus")}),
		text(&pictures.Mouse{Path: assetPath("mouse.mp4")}),
		text(&pictures.Tiktok{Path: assetPath("tiktok")}),
		text(&pictures.Hello{Path: assetPath("hello.json")}),
		text(&pictures.Anime{}),
		text(&pictures.Furry{}),
		text(&pictures.Flag{}),
		text(&pictures.Car{}),
		text(&pictures.Soy{}),
		text(&pictures.Danbooru{API: a.danbooru}),
		text(&pictures.Fap{API: a.danbooru}),
		text(&pictures.Masyunya{}),
		text(&pictures.Poppy{}),
		text(&pictures.Sima{}),
		text(&pictures.Cat{}),
	}
}

func (a *app) casinoServices() []server.Service {
	return []server.Service{
		&casino.Roll{Universe: a.universe},
		text(&casino.Dice{Universe: a.universe}),
		text(&casino.Slot{Universe: a.universe, MinBet: 100}),
	}
}

func (a *app) nextServices() []server.Service {
	return []server.Service{&router.Service{Universe: a.universe}}
}

func text(s adapter.TextService) server.Service {
	return &adapter.Text{TextService: s}
}

func wrap(s server.Service, w ...adapter.Wrapper) server.Service {
	for i := len(w) - 1; i >= 0; i-- {
		s = adapter.Wrap(s, w[i])
	}
	return s
}

func concat[T any](slices ...[]T) []T {
	var r []T
	for _, s := range slices {
		r = append(r, s...)
	}
	return r
}
