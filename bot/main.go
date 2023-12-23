package main

import (
	"fmt"
	"log"
	"nechego/avatar"
	"nechego/bot/middleware"
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
	"nechego/handlers/profile"
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
	srv := &Server{
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
			addService(w, refreshMarket, time.Minute)
			addService(w, restoreEnergy, time.Minute)
			w.Market.OnBuy = onBuyHandler(w)
			w.Market.OnSell = onSellHandler(w)
		}),
		avatars:  &avatar.Storage{Dir: avatarDirectory},
		danbooru: danbooru.New(danbooru.URL, 5*time.Second, 3),
	}, nil
}

func (a *app) shutdown() error {
	return a.universe.SaveAll()
}

func (a *app) services() []Service {
	global := a.globalMiddleware()
	spam := []Wrapper{&middleware.AutoDelete{After: 15 * time.Minute}}

	groups := []struct {
		services   []Service
		middleware []Wrapper
	}{
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

	var handlers []Service
	for _, g := range groups {
		for _, s := range g.services {
			var w []Wrapper
			w = append(w, g.middleware...)
			w = append(w, global...)
			h := Wrap(s, w)
			handlers = append(handlers, h)
		}
	}
	return handlers
}

func (a *app) globalMiddleware() []Wrapper {
	return []Wrapper{
		&middleware.Recover{},
		&middleware.RandomPhoto{Avatars: a.avatars, Prob: 1. / 200},
		&middleware.IncrementCounters{Universe: a.universe},
		&middleware.CacheName{Universe: a.universe},
		&middleware.IgnoreWorldInactive{
			Universe: a.universe,
			Immune: func(c tele.Context) bool {
				var h fun.TurnOn
				return h.Match(c)
			}},
		&middleware.Throttle{Duration: 800 * time.Millisecond},
		&middleware.LogMessage{Wait: 5 * time.Second},
		&middleware.IgnoreMessageForwarded{},
		&middleware.RequireSupergroup{},
	}
}

func (a *app) dailyServices() []Service {
	return []Service{
		&daily.Eblan{Universe: a.universe},
		&daily.Admin{Universe: a.universe},
		&daily.Pair{Universe: a.universe},
	}
}

func (a *app) economyServices() []Service {
	return []Service{
		&economy.Inventory{Universe: a.universe},
		&economy.Send{Universe: a.universe},
		&economy.Mail{Universe: a.universe},
		&economy.Sort{Universe: a.universe},
		&economy.Drop{Universe: a.universe},
		&economy.Pick{Universe: a.universe},
		&economy.Floor{Universe: a.universe},
		&economy.Stack{Universe: a.universe},
		&economy.Split{Universe: a.universe},
		&economy.Cashout{Universe: a.universe},
		&economy.Capital{Universe: a.universe},
		&economy.Balance{Universe: a.universe},
	}
}

func (a *app) farmServices() []Service {
	return []Service{
		&farm.Farm{Universe: a.universe},
		&farm.Plant{Universe: a.universe},
		&farm.Harvest{Universe: a.universe},
		&farm.Upgrade{Universe: a.universe},
		&farm.Name{Universe: a.universe},
	}
}

func (a *app) marketServices() []Service {
	return []Service{
		&market.Market{Universe: a.universe},
		&market.PriceList{Universe: a.universe},
		&market.Buy{Universe: a.universe},
		&market.Sell{Universe: a.universe},
		&market.SellQuick{Universe: a.universe},
		&market.Name{Universe: a.universe},
		&market.Job{Universe: a.universe},
	}
}

func (a *app) actionsServices() []Service {
	return []Service{
		&actions.Fish{Universe: a.universe},
		&actions.Craft{Universe: a.universe},
		&actions.Fight{Universe: a.universe},
		&actions.Eat{Universe: a.universe},
		&actions.EatQuick{Universe: a.universe},
		&actions.Records{Universe: a.universe},
		&actions.Friends{Universe: a.universe},
		&actions.Write{Universe: a.universe},
		&actions.Open{Universe: a.universe},
	}
}

func (a *app) topServices() []Service {
	return []Service{
		top.Rating(a.universe),
		top.Rich(a.universe),
		top.Strength(a.universe),
	}
}

func (a *app) profileServices() []Service {
	return []Service{
		&profile.Status{Universe: a.universe, MaxLength: 140},
		&profile.Profile{Universe: a.universe, Avatars: a.avatars},
		&profile.Avatar{Avatars: a.avatars, MaxWidth: 1500, MaxHeight: 1500},
		&profile.Energy{Universe: a.universe},
		&profile.NamePet{Universe: a.universe},
	}
}

func (a *app) funServices() []Service {
	return []Service{
		&fun.Game{},
		&fun.Infa{},
		&fun.Choose{},
		&fun.Weather{},
		&fun.Calc{},
		&fun.Name{},
		&fun.CheckName{},
		&fun.Who{Universe: a.universe},
		&fun.List{Universe: a.universe},
		&fun.Top{Universe: a.universe},
		&fun.Clock{},
		&fun.Date{},
		&fun.TurnOn{Universe: a.universe},
		&fun.TurnOff{Universe: a.universe},
		&fun.Reputation{Universe: a.universe},
		&fun.UpdateReputation{Universe: a.universe},
		&fun.NewYear{},
	}
}

func (a *app) pictureServices() []Service {
	return []Service{
		&pictures.Pic{Path: assetPath("pic")},
		&pictures.Basili{Path: assetPath("basili")},
		&pictures.Casper{Path: assetPath("casper")},
		&pictures.Zeus{Path: assetPath("zeus")},
		&pictures.Mouse{Path: assetPath("mouse.mp4")},
		&pictures.Tiktok{Path: assetPath("tiktok")},
		&pictures.Hello{Path: assetPath("hello.json")},
		&pictures.Anime{},
		&pictures.Furry{},
		&pictures.Flag{},
		&pictures.Car{},
		&pictures.Soy{},
		&pictures.Danbooru{API: a.danbooru},
		&pictures.Fap{API: a.danbooru},
		&pictures.Masyunya{},
		&pictures.Poppy{},
		&pictures.Sima{},
		&pictures.Cat{},
	}
}

func (a *app) casinoServices() []Service {
	return []Service{
		&casino.DiceRoll{Universe: a.universe},
		&casino.SlotRoll{Universe: a.universe},
		&casino.Dice{Universe: a.universe, MinBet: 100},
		&casino.Slot{Universe: a.universe, MinBet: 100},
	}
}

func (a *app) commandServices() []Service {
	return []Service{
		&command.Add{Universe: a.universe},
		&command.Remove{Universe: a.universe},
		&command.Use{Universe: a.universe},
	}
}

func (a *app) otherServices() []Service {
	return []Service{
		&handlers.Help{},
		&handlers.Pass{},
	}
}
