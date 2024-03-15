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
		&fun.CheckName{},
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
		&pictures.Pic{Path: assetPath("pic")},
		&pictures.Basili{Path: assetPath("basili")},
		&pictures.Casper{Path: assetPath("casper")},
		&pictures.Zeus{Path: assetPath("zeus")},
		&pictures.Mouse{Path: assetPath("mouse.mp4")},
		&pictures.Tiktok{Path: assetPath("tiktok")},
		&pictures.Hello{Queries: a.Queries},
		&pictures.Anime{},
		&pictures.Furry{},
		&pictures.Flag{},
		&pictures.Car{},
		&pictures.Soy{},
		&pictures.Danbooru{},
		&pictures.Fap{},
		&pictures.Masyunya{},
		&pictures.Poppy{},
		&pictures.Sima{},
		&pictures.Cat{},
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
