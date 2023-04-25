package main

import (
	"nechego/avatar"
	"nechego/danbooru"
	"nechego/game"
	"nechego/handlers"
	"nechego/handlers/casino"
	"nechego/handlers/pictures"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

type TextHandler struct {
	StringService
}

func (h *TextHandler) Match(c tele.Context) bool {
	return h.StringService.Match(c.Text())
}

type CallbackHandler struct {
	StringService
}

func (h *CallbackHandler) Match(c tele.Context) bool {
	cb := c.Callback()
	if cb == nil {
		return false
	}
	cb.Data = strings.TrimSpace(cb.Data)
	return h.StringService.Match(cb.Data)
}

type ClosureHandler struct {
	ContextService
	H func(tele.Context) error
}

func (h *ClosureHandler) Handle(c tele.Context) error {
	return h.H(c)
}

func Wrap(s ContextService, w ...Wrapper) *ClosureHandler {
	h := &ClosureHandler{ContextService: s, H: s.Handle}
	for i := len(w) - 1; i >= 0; i-- {
		h.H = w[i].Wrap(h.H)
	}
	return h
}

func PictureHandlers() []ContextService {
	d := danbooru.New(danbooru.URL, 5*time.Second, 3)
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
		&pictures.Danbooru{API: d}, // TODO: add Settings, singular design
		&pictures.Fap{API: d},
		&pictures.Masyunya{},
		&pictures.Poppy{},
		&pictures.Sima{},
	} {
		r = append(r, &TextHandler{s})
	}
	return r
}

func CasinoHandlers(u *game.Universe) []ContextService {
	r := []ContextService{
		&casino.Roll{Universe: u},
	}
	for _, s := range []StringService{
		&casino.Dice{Universe: u},
		&casino.Slot{Universe: u, MinBet: 100},
	} {
		r = append(r, &TextHandler{s})
	}
	return r
}

func RemainingHandlers(u *game.Universe, as *avatar.Storage) []StringService {
	return []StringService{
		&handlers.Help{},

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
		&handlers.NameFarm{Universe: u},

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
		&handlers.Fight{Universe: u},
		&handlers.PvP{Universe: u},
		&handlers.Eat{Universe: u},
		&handlers.EatQuick{Universe: u},
		&handlers.FishingRecords{Universe: u},
		&handlers.Friends{Universe: u},
		&handlers.Transfer{Universe: u},
		&handlers.Use{Universe: u},

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

func CallbackHandlers(u *game.Universe) []ContextService {
	r := []ContextService{}
	for _, s := range []StringService{
		&handlers.HarvestInline{Universe: u},
		&handlers.AuctionBuy{Universe: u},
	} {
		r = append(r, &CallbackHandler{s})
	}
	return r
}
