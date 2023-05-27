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
			w.History.Announce(handlers.RecordAnnouncer(bot, tele.ChatID(w.TGID)))
		}),
		avatars:  &avatar.Storage{bot, avatarDirectory, 1500, 1500},
		danbooru: danbooru.New(danbooru.URL, 5*time.Second, 3),
	}, nil
}

func (a *app) shutdown() error {
	return a.universe.SaveAll()
}

func (a *app) services() []server.Service {
	spam := []adapter.Wrapper{&middleware.AutoDelete{After: 5 * time.Minute}}
	global := a.globalMiddleware()
	groups := []struct {
		services   []server.Service
		middleware []adapter.Wrapper
	}{
		{a.informationServices(), nil},
		{a.dailyServices(), nil},
		{a.economyServices(), spam},
		{a.farmServices(), spam},
		{a.marketServices(), spam},
		{a.auctionServices(), spam},
		{a.actionsServices(), spam},
		{a.topServices(), spam},
		{a.profileServices(), spam},
		{a.phoneServices(), spam},
		{a.funServices(), nil},
		{a.pictureServices(), nil},
		{a.casinoServices(), spam},
		{a.callbackServices(), nil},
	}
	handlers := []server.Service{}
	for _, g := range groups {
		for _, s := range g.services {
			s = adapter.Wrap(s, g.middleware...)
			s = adapter.Wrap(s, global...)
			handlers = append(handlers, s)
		}
	}
	return handlers
}

func (a *app) globalMiddleware() []adapter.Wrapper {
	return []adapter.Wrapper{
		&middleware.Recover{},
		&middleware.RequireSupergroup{},
		&middleware.IgnoreForwarded{},
		&middleware.IgnoreBanned{Universe: a.universe},
		&middleware.LogMessage{Wait: 5 * time.Second},
		&middleware.Throttle{Duration: 800 * time.Millisecond},
		&middleware.IncrementCounters{Universe: a.universe},
		&middleware.RandomPhoto{Avatars: a.avatars, Prob: 1. / 200},
	}
}

func (a *app) informationServices() []server.Service {
	r := []server.Service{}
	for _, s := range []adapter.TextService{
		&handlers.Help{},
	} {
		r = append(r, &adapter.Text{s})
	}
	return r
}

func (a *app) dailyServices() []server.Service {
	r := []server.Service{}
	for _, s := range []adapter.TextService{
		&handlers.DailyEblan{Universe: a.universe},
		&handlers.DailyAdmin{Universe: a.universe},
		&handlers.DailyPair{Universe: a.universe},
	} {
		r = append(r, &adapter.Text{s})
	}
	return r
}

func (a *app) economyServices() []server.Service {
	r := []server.Service{}
	for _, s := range []adapter.TextService{
		&handlers.Inventory{Universe: a.universe},
		&handlers.Funds{Universe: a.universe},
		&handlers.Sort{Universe: a.universe},
		&handlers.Drop{Universe: a.universe},
		&handlers.Pick{Universe: a.universe},
		&handlers.Floor{Universe: a.universe},
		&handlers.Stack{Universe: a.universe},
		&handlers.Split{Universe: a.universe},
		&handlers.Cashout{Universe: a.universe},
		&handlers.Capital{Universe: a.universe},
		&handlers.Balance{Universe: a.universe},
	} {
		r = append(r, &adapter.Text{s})
	}
	return r
}

func (a *app) farmServices() []server.Service {
	r := []server.Service{}
	for _, s := range []adapter.TextService{
		&handlers.Farm{Universe: a.universe},
		&handlers.Plant{Universe: a.universe},
		&handlers.Harvest{Universe: a.universe},
		&handlers.PriceList{Universe: a.universe},
		&handlers.UpgradeFarm{Universe: a.universe},
		&handlers.NameFarm{Universe: a.universe},
	} {
		r = append(r, &adapter.Text{s})
	}
	return r
}

func (a *app) marketServices() []server.Service {
	r := []server.Service{}
	for _, s := range []adapter.TextService{
		&handlers.Market{Universe: a.universe},
		&handlers.Buy{Universe: a.universe},
		&handlers.Sell{Universe: a.universe},
		&handlers.SellQuick{Universe: a.universe},
		&handlers.NameMarket{Universe: a.universe},
		&handlers.GetJob{Universe: a.universe},
		&handlers.QuitJob{Universe: a.universe},
	} {
		r = append(r, &adapter.Text{s})
	}
	return r
}

func (a *app) auctionServices() []server.Service {
	r := []server.Service{}
	for _, s := range []adapter.TextService{
		&handlers.Auction{Universe: a.universe},
		&handlers.AuctionSell{Universe: a.universe},
	} {
		r = append(r, &adapter.Text{s})
	}
	return r
}

func (a *app) actionsServices() []server.Service {
	r := []server.Service{}
	for _, s := range []adapter.TextService{
		&handlers.Craft{Universe: a.universe},
		&handlers.Fish{Universe: a.universe},
		&handlers.DrawNet{Universe: a.universe},
		&handlers.CastNet{Universe: a.universe},
		&handlers.Net{Universe: a.universe},
		&handlers.Catch{Universe: a.universe},
		&handlers.Fight{Universe: a.universe},
		&handlers.PvP{Universe: a.universe},
		&handlers.Eat{Universe: a.universe},
		&handlers.EatQuick{Universe: a.universe},
		&handlers.FishingRecords{Universe: a.universe},
		&handlers.Friends{Universe: a.universe},
		&handlers.Transfer{Universe: a.universe},
		&handlers.Use{Universe: a.universe},
	} {
		r = append(r, &adapter.Text{s})
	}
	return r
}

func (a *app) topServices() []server.Service {
	r := []server.Service{}
	for _, s := range []adapter.TextService{
		&handlers.TopStrong{Universe: a.universe},
		&handlers.TopRating{Universe: a.universe},
		&handlers.TopRich{Universe: a.universe},
	} {
		r = append(r, &adapter.Text{s})
	}
	return r
}

func (a *app) profileServices() []server.Service {
	r := []server.Service{}
	for _, s := range []adapter.TextService{
		&handlers.Status{Universe: a.universe, MaxLength: 120},
		&handlers.Profile{Universe: a.universe, Avatars: a.avatars},
		&handlers.Avatar{Universe: a.universe, Avatars: a.avatars},
		&handlers.Energy{Universe: a.universe},
		&handlers.NamePet{Universe: a.universe},
	} {
		r = append(r, &adapter.Text{s})
	}
	return r
}

func (a *app) phoneServices() []server.Service {
	r := []server.Service{}
	for _, s := range []adapter.TextService{
		&handlers.SendSMS{Universe: a.universe},
		&handlers.ReceiveSMS{Universe: a.universe},
		&handlers.Contacts{Universe: a.universe},
		&handlers.Spam{Universe: a.universe},
	} {
		r = append(r, &adapter.Text{s})
	}
	return r
}

func (a *app) funServices() []server.Service {
	r := []server.Service{}
	for _, s := range []adapter.TextService{
		&handlers.Game{},
		&handlers.Infa{},
		&handlers.Weather{},
		&handlers.Calculator{},
		&handlers.Name{},
		&handlers.Who{Universe: a.universe},
		&handlers.List{Universe: a.universe},
		&handlers.Top{Universe: a.universe},
		&fun.Time{},
		&handlers.TurnOn{Universe: a.universe},
		&handlers.TurnOff{Universe: a.universe},
	} {
		r = append(r, &adapter.Text{s})
	}
	return r
}

func (a *app) pictureServices() []server.Service {
	r := []server.Service{}
	for _, s := range []adapter.TextService{
		&pictures.Pic{Path: assetPath("pic")},
		&pictures.Basili{Path: assetPath("basili")},
		&pictures.Casper{Path: assetPath("casper")},
		&pictures.Zeus{Path: assetPath("zeus")},
		&pictures.Mouse{Path: assetPath("mouse.mp4")},
		&pictures.Tiktok{Path: assetPath("tiktok")},
		&pictures.Hello{Path: assetPath("hello.json")}, // TODO: is cache initialized once?
		&pictures.Anime{},
		&pictures.Furry{},
		&pictures.Flag{},
		&pictures.Car{},
		&pictures.Soy{},
		&pictures.Danbooru{API: a.danbooru}, // TODO: add Settings, singular design
		&pictures.Fap{API: a.danbooru},      // TODO: same
		&pictures.Masyunya{},
		&pictures.Poppy{},
		&pictures.Sima{},
	} {
		r = append(r, &adapter.Text{s})
	}
	return r
}

func (a *app) casinoServices() []server.Service {
	r := []server.Service{
		&casino.Roll{Universe: a.universe},
	}
	for _, s := range []adapter.TextService{
		&casino.Dice{Universe: a.universe},
		&casino.Slot{Universe: a.universe, MinBet: 100},
	} {
		r = append(r, &adapter.Text{s})
	}
	return r
}

func (a *app) callbackServices() []server.Service {
	r := []server.Service{}
	for _, s := range []adapter.TextService{
		&handlers.HarvestInline{Universe: a.universe},
		&handlers.AuctionBuy{Universe: a.universe},
	} {
		r = append(r, &adapter.Callback{s})
	}
	return r
}
