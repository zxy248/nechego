package main

import (
	"nechego/avatar"
	"nechego/bot/middleware"
	"nechego/danbooru"
	"nechego/game"
	"nechego/handlers"
	"nechego/handlers/casino"
	"nechego/handlers/pictures"
	"time"

	tele "gopkg.in/telebot.v3"
)

type deps struct {
	teleBot       *tele.Bot
	gameUniverse  *game.Universe
	avatarStorage *avatar.Storage
	danbooru      *danbooru.Danbooru
}

func (d *deps) router() router {
	remaining := d.remainingServices()
	picture := d.pictureServices()
	casino := d.casinoServices()
	callback := d.callbackServices()

	groups := [][]ContextService{
		remaining,
		picture,
		casino,
		callback,
	}
	handlers := []ContextService{}
	for _, g := range groups {
		handlers = append(handlers, g...)
	}

	mw := d.middleware()
	for i, h := range handlers {
		handlers[i] = wrap(h, mw...)
	}
	return handlers
}

func (d *deps) middleware() []Wrapper {
	return []Wrapper{
		&middleware.Recover{},
		&middleware.RequireSupergroup{},
		&middleware.IgnoreForwarded{},
		&middleware.IgnoreBanned{Universe: d.gameUniverse},
		&middleware.LogMessage{Wait: 5 * time.Second},
		&middleware.Throttle{Duration: 800 * time.Millisecond},
		&middleware.IncrementCounters{Universe: d.gameUniverse},
		&middleware.RandomPhoto{Avatars: d.avatarStorage, Prob: 1. / 200},
	}
}

func (d *deps) pictureServices() []ContextService {
	r := []ContextService{}
	for _, s := range []StringService{
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
		&pictures.Danbooru{API: d.danbooru}, // TODO: add Settings, singular design
		&pictures.Fap{API: d.danbooru},      // TODO: same
		&pictures.Masyunya{},
		&pictures.Poppy{},
		&pictures.Sima{},
	} {
		r = append(r, &TextHandler{s})
	}
	return r
}

func (d *deps) casinoServices() []ContextService {
	r := []ContextService{
		&casino.Roll{Universe: d.gameUniverse},
	}
	for _, s := range []StringService{
		&casino.Dice{Universe: d.gameUniverse},
		&casino.Slot{Universe: d.gameUniverse, MinBet: 100},
	} {
		r = append(r, &TextHandler{s})
	}
	return r
}

func (d *deps) remainingServices() []ContextService {
	// TODO: group handlers
	r := []ContextService{}
	for _, s := range []StringService{
		&handlers.Help{},

		// Daily.
		&handlers.DailyEblan{Universe: d.gameUniverse},
		&handlers.DailyAdmin{Universe: d.gameUniverse},
		&handlers.DailyPair{Universe: d.gameUniverse},

		// Economy.
		&handlers.Inventory{Universe: d.gameUniverse},
		&handlers.Funds{Universe: d.gameUniverse},
		&handlers.Sort{Universe: d.gameUniverse},
		&handlers.Drop{Universe: d.gameUniverse},
		&handlers.Pick{Universe: d.gameUniverse},
		&handlers.Floor{Universe: d.gameUniverse},
		&handlers.Stack{Universe: d.gameUniverse},
		&handlers.Split{Universe: d.gameUniverse},
		&handlers.Cashout{Universe: d.gameUniverse},
		&handlers.Capital{Universe: d.gameUniverse},
		&handlers.Balance{Universe: d.gameUniverse},

		// Farm.
		&handlers.Farm{Universe: d.gameUniverse},
		&handlers.Plant{Universe: d.gameUniverse},
		&handlers.Harvest{Universe: d.gameUniverse},
		&handlers.PriceList{Universe: d.gameUniverse},
		&handlers.UpgradeFarm{Universe: d.gameUniverse},
		&handlers.NameFarm{Universe: d.gameUniverse},

		// Market.
		&handlers.Market{Universe: d.gameUniverse},
		&handlers.Buy{Universe: d.gameUniverse},
		&handlers.Sell{Universe: d.gameUniverse},
		&handlers.SellQuick{Universe: d.gameUniverse},
		&handlers.NameMarket{Universe: d.gameUniverse},
		&handlers.GetJob{Universe: d.gameUniverse},
		&handlers.QuitJob{Universe: d.gameUniverse},

		// Auction.
		&handlers.Auction{Universe: d.gameUniverse},
		&handlers.AuctionSell{Universe: d.gameUniverse},

		// Actions.
		&handlers.Craft{Universe: d.gameUniverse},
		&handlers.Fish{Universe: d.gameUniverse},
		&handlers.DrawNet{Universe: d.gameUniverse},
		&handlers.CastNet{Universe: d.gameUniverse},
		&handlers.Net{Universe: d.gameUniverse},
		&handlers.Catch{Universe: d.gameUniverse},
		&handlers.Fight{Universe: d.gameUniverse},
		&handlers.PvP{Universe: d.gameUniverse},
		&handlers.Eat{Universe: d.gameUniverse},
		&handlers.EatQuick{Universe: d.gameUniverse},
		&handlers.FishingRecords{Universe: d.gameUniverse},
		&handlers.Friends{Universe: d.gameUniverse},
		&handlers.Transfer{Universe: d.gameUniverse},
		&handlers.Use{Universe: d.gameUniverse},

		// Top.
		&handlers.TopStrong{Universe: d.gameUniverse},
		&handlers.TopRating{Universe: d.gameUniverse},
		&handlers.TopRich{Universe: d.gameUniverse},

		// Profile.
		&handlers.Status{Universe: d.gameUniverse, MaxLength: 120},
		&handlers.Profile{Universe: d.gameUniverse, Avatars: d.avatarStorage},
		&handlers.Avatar{Universe: d.gameUniverse, Avatars: d.avatarStorage},
		&handlers.Energy{Universe: d.gameUniverse},
		&handlers.NamePet{Universe: d.gameUniverse},

		// Phone.
		&handlers.SendSMS{Universe: d.gameUniverse},
		&handlers.ReceiveSMS{Universe: d.gameUniverse},
		&handlers.Contacts{Universe: d.gameUniverse},
		&handlers.Spam{Universe: d.gameUniverse},

		// Fun.
		&handlers.Game{},
		&handlers.Infa{},
		&handlers.Weather{},
		&handlers.Calculator{},
		&handlers.Name{},
		&handlers.Who{Universe: d.gameUniverse},
		&handlers.List{Universe: d.gameUniverse},
		&handlers.TurnOn{Universe: d.gameUniverse},
		&handlers.TurnOff{Universe: d.gameUniverse},
		&handlers.Top{Universe: d.gameUniverse},
	} {
		r = append(r, &TextHandler{s})
	}
	return r
}

func (d *deps) callbackServices() []ContextService {
	r := []ContextService{}
	for _, s := range []StringService{
		&handlers.HarvestInline{Universe: d.gameUniverse},
		&handlers.AuctionBuy{Universe: d.gameUniverse},
	} {
		r = append(r, &CallbackHandler{s})
	}
	return r
}
