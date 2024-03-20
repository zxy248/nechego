package main

import (
	"path/filepath"

	"github.com/zxy248/nechego/bot/middleware"
	"github.com/zxy248/nechego/data"
	"github.com/zxy248/nechego/handlers"
	"github.com/zxy248/nechego/handlers/command"
	"github.com/zxy248/nechego/handlers/daily"
	"github.com/zxy248/nechego/handlers/fun"
	"github.com/zxy248/nechego/handlers/pictures"
)

type App struct {
	Queries *data.Queries
}

func (a *App) Handlers() []Handler {
	var hs []Handler
	hs = append(hs, a.dailyHandlers()...)
	hs = append(hs, a.funHandlers()...)
	hs = append(hs, a.pictureHandlers()...)
	hs = append(hs, a.commandHandlers()...)
	hs = append(hs, a.otherHandlers()...)

	mw := a.middleware()
	var r []Handler
	for _, h := range hs {
		h = &InstrumentedHandler{a.Queries, h}
		r = append(r, Wrap(h, mw))
	}
	return r
}

func (a *App) middleware() []Wrapper {
	return []Wrapper{
		&middleware.Recover{},
		&middleware.RandomPhoto{Prob: 0.005},
		&middleware.RandomReact{Prob: 0.033},
		&middleware.RandomSticker{Queries: a.Queries, Prob: 0.02},
		&middleware.IgnoreInactive{Queries: a.Queries, Immune: (&fun.TurnOn{}).Match},
		&middleware.UpdateDaily{Queries: a.Queries},
		&middleware.LogMessage{Queries: a.Queries},
		&middleware.UpdateInfo{Queries: a.Queries},
		&middleware.IgnoreMessageForwarded{},
		&middleware.RequireSupergroup{},
	}
}

func (a *App) dailyHandlers() []Handler {
	return []Handler{
		&daily.Eblan{Queries: a.Queries},
		&daily.Admin{Queries: a.Queries},
		&daily.Pair{Queries: a.Queries},
	}
}

func (a *App) funHandlers() []Handler {
	return []Handler{
		&fun.Game{},
		&fun.Infa{},
		&fun.Choose{},
		&fun.Weather{},
		&fun.Calc{},
		&fun.Name{},
		&fun.Who{Queries: a.Queries},
		&fun.List{Queries: a.Queries},
		&fun.Top{Queries: a.Queries},
		&fun.Clock{},
		&fun.Date{},
		&fun.TurnOn{Queries: a.Queries},
		&fun.TurnOff{Queries: a.Queries},
		&fun.NewYear{},
		&fun.Photon{},
		&fun.Avatar{},
	}
}

func (a *App) pictureHandlers() []Handler {
	return []Handler{
		&pictures.FromDir{
			Path:   assetPath("pic"),
			Regexp: handlers.NewRegexp("^!пик"),
		},
		&pictures.FromDir{
			Path:   assetPath("basili"),
			Regexp: handlers.NewRegexp("^!(муся|марс|(кот|кош)[а-я]* василия)"),
		},
		&pictures.FromDir{
			Path:   assetPath("casper"),
			Regexp: handlers.NewRegexp("^!касп[ие]р"),
		},
		&pictures.FromDir{
			Path:   assetPath("zeus"),
			Regexp: handlers.NewRegexp("^!зевс"),
		},
		&pictures.FromDir{
			Path:   assetPath("mouse"),
			Regexp: handlers.NewRegexp("!мыш"),
		},
		&pictures.FromDir{
			Path:   assetPath("tiktok"),
			Regexp: handlers.NewRegexp("^!тикток"),
		},
		&pictures.FromURL{
			Locator: &pictures.Anime{},
			Regexp:  handlers.NewRegexp("^!(аним|мульт)"),
		},
		&pictures.FromURL{
			Locator: &pictures.Furry{},
			Regexp:  handlers.NewRegexp("^!фур"),
		},
		&pictures.FromURL{
			Locator: &pictures.Flag{},
			Regexp:  handlers.NewRegexp("^!флаг"),
		},
		&pictures.FromURL{
			Locator: &pictures.Soy{},
			Regexp:  handlers.NewRegexp("^!сой"),
		},
		&pictures.FromURL{
			Locator: &pictures.Cat{},
			Regexp:  handlers.NewRegexp("!(кот|кош)"),
		},
		&pictures.Car{},
		&pictures.Danbooru{},
		&pictures.Fap{},
		&pictures.FromStickerPack{
			Source: []string{"masyunya_vk"},
			Regexp: handlers.NewRegexp("^!ма[нс]ю[нс][а-я]*[пая]"),
		},
		&pictures.FromStickerPack{
			Source: []string{"pappy2_vk", "poppy_vk"},
			Regexp: handlers.NewRegexp("^!паппи"),
		},
		&pictures.FromStickerPack{
			Source: []string{"catsima_vk"},
			Regexp: handlers.NewRegexp("^!сима"),
		},
		&pictures.FromStickerPack{
			Source: []string{"Vjopeneogurez_by_fStikBot"},
			Regexp: handlers.NewRegexp("^!зюз[а-я]*"),
		},
		&pictures.FromStickerPack{
			Source: []string{"jdjsjakwkek_by_fStikBot"},
			Regexp: handlers.NewRegexp("^!с[еи]л[еи]н[еи]?[еэ]ль?"),
		},
		&pictures.Hello{},
	}
}

func (a *App) commandHandlers() []Handler {
	return []Handler{
		&command.Add{Queries: a.Queries},
		&command.Remove{Queries: a.Queries},
		&command.Use{Queries: a.Queries},
	}
}

func (a *App) otherHandlers() []Handler {
	return []Handler{
		&fun.Speak{Queries: a.Queries, Attempts: 50},
		&handlers.Pass{Queries: a.Queries},
	}
}

func assetPath(s string) string {
	return filepath.Join(assetsDirectory, s)
}
