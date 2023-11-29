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
	"nechego/handlers/actions"
	"nechego/handlers/casino"
	"nechego/handlers/command"
	"nechego/handlers/daily"
	"nechego/handlers/economy"
	"nechego/handlers/farm"
	"nechego/handlers/fun"
	"nechego/handlers/market"
	"nechego/handlers/pictures"
	"nechego/handlers/top"
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
	botToken        = getenv("NECHEGO_TOKEN")
	assetsDirectory = getenv("NECHEGO_ASSETS")
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

func getenv(s string) string {
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
			refreshMarket(w)
			resetEnergy(w)
			addService(w, refreshMarket, time.Minute)
			addService(w, restoreEnergy, time.Minute)
			addService(w, fillNets, time.Minute)
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
		{a.commandServices(), nil},
		{a.otherServices(), nil},
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
				var h fun.TurnOn
				return h.Match(c)
			},
		},
		&middleware.CacheName{Universe: a.universe},
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
		&daily.Eblan{Universe: a.universe},
		&daily.Admin{Universe: a.universe},
		&daily.Pair{Universe: a.universe},
	}
}

func (a *app) economyServices() []server.Service {
	return []server.Service{
		&economy.Inventory{Universe: a.universe},
		&economy.Funds{Universe: a.universe},
		&economy.Sort{Universe: a.universe},
		&economy.Drop{Universe: a.universe},
		&economy.Pick{Universe: a.universe},
		&economy.Floor{Universe: a.universe},
		&economy.Stack{Universe: a.universe},
		&economy.Split{Universe: a.universe},
		&economy.Cashout{Universe: a.universe},
		text(&handlers.Capital{Universe: a.universe}),
		text(&handlers.Balance{Universe: a.universe}),
	}
}

func (a *app) farmServices() []server.Service {
	return []server.Service{
		&farm.Farm{Universe: a.universe},
		&farm.Plant{Universe: a.universe},
		&farm.Harvest{Universe: a.universe},
		&farm.Upgrade{Universe: a.universe},
		&farm.Name{Universe: a.universe},
	}
}

func (a *app) marketServices() []server.Service {
	return []server.Service{
		text(&handlers.Market{Universe: a.universe}),
		text(&handlers.PriceList{Universe: a.universe}),
		&market.Buy{Universe: a.universe},
		text(&handlers.Sell{Universe: a.universe}),
		text(&handlers.SellQuick{Universe: a.universe}),
		text(&handlers.NameMarket{Universe: a.universe}),
		text(&handlers.GetJob{Universe: a.universe}),
		text(&handlers.QuitJob{Universe: a.universe}),
	}
}

func (a *app) actionsServices() []server.Service {
	return []server.Service{
		&actions.Fish{Universe: a.universe},
		&actions.Craft{Universe: a.universe},
		text(&handlers.DrawNet{Universe: a.universe}),
		text(&handlers.CastNet{Universe: a.universe}),
		text(&handlers.Net{Universe: a.universe}),
		text(&handlers.Catch{Universe: a.universe}),
		text(&handlers.Fight{Universe: a.universe}),
		&actions.Eat{Universe: a.universe},
		&actions.EatQuick{Universe: a.universe},
		text(&handlers.FishingRecords{Universe: a.universe}),
		text(&handlers.Friends{Universe: a.universe}),
		text(&handlers.Transfer{Universe: a.universe}),
	}
}

func (a *app) topServices() []server.Service {
	return []server.Service{
		top.Rating(a.universe),
		top.Rich(a.universe),
		top.Strength(a.universe),
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
		&fun.Game{},
		&fun.Infa{},
		&fun.Weather{},
		&fun.Calc{},
		text(&fun.Name{}),
		text(&fun.CheckName{}),
		text(&handlers.Who{Universe: a.universe}),
		text(&handlers.List{Universe: a.universe}),
		text(&handlers.Top{Universe: a.universe}),
		&fun.Clock{},
		&fun.TurnOn{Universe: a.universe},
		&fun.TurnOff{Universe: a.universe},
		&fun.Reputation{Universe: a.universe},
		&fun.UpdateReputation{Universe: a.universe},
	}
}

func (a *app) pictureServices() []server.Service {
	return []server.Service{
		text(&pictures.Pic{Path: assetPath("pic")}),
		&pictures.Basili{Path: assetPath("basili")},
		&pictures.Casper{Path: assetPath("casper")},
		text(&pictures.Zeus{Path: assetPath("zeus")}),
		text(&pictures.Mouse{Path: assetPath("mouse.mp4")}),
		text(&pictures.Tiktok{Path: assetPath("tiktok")}),
		&pictures.Hello{Path: assetPath("hello.json")},
		text(&pictures.Anime{}),
		text(&pictures.Furry{}),
		text(&pictures.Flag{}),
		text(&pictures.Car{}),
		text(&pictures.Soy{}),
		text(&pictures.Danbooru{API: a.danbooru}),
		&pictures.Fap{API: a.danbooru},
		&pictures.Masyunya{},
		text(&pictures.Poppy{}),
		text(&pictures.Sima{}),
		text(&pictures.Cat{}),
	}
}

func (a *app) casinoServices() []server.Service {
	return []server.Service{
		&casino.DiceRoll{Universe: a.universe},
		&casino.SlotRoll{Universe: a.universe},
		&casino.Dice{Universe: a.universe, MinBet: 100},
		&casino.Slot{Universe: a.universe, MinBet: 100},
	}
}

func (a *app) commandServices() []server.Service {
	return []server.Service{
		&command.Add{Universe: a.universe},
		&command.Remove{Universe: a.universe},
		&command.Use{Universe: a.universe},
	}
}

func (a *app) otherServices() []server.Service {
	return []server.Service{
		&handlers.Pass{},
	}
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
