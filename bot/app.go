package main

import (
	"nechego/bot/middleware"
	"nechego/danbooru"
	"nechego/game"
	"nechego/handlers"
	"nechego/handlers/command"
	"nechego/handlers/daily"
	"nechego/handlers/fun"
	"nechego/handlers/pictures"
	"path/filepath"
	"time"
)

type App struct {
	Universe *game.Universe
	Danbooru *danbooru.Danbooru
}

func (a *App) Shutdown() error {
	return a.Universe.SaveAll()
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
		&middleware.RandomReact{Prob: 0.03},
		&middleware.IgnoreWorldInactive{
			Universe: a.Universe,
			Immune:   (&fun.TurnOn{}).Match,
		},
		&middleware.Throttle{Duration: 400 * time.Millisecond},
		&middleware.LogMessage{Wait: 5 * time.Second},
		&middleware.IgnoreMessageForwarded{},
		&middleware.RequireSupergroup{},
		&middleware.AddUser{Universe: a.Universe},
	}
}

func (a *App) dailyHandlers() []Handler {
	return []Handler{
		&daily.Eblan{Universe: a.Universe},
		&daily.Admin{Universe: a.Universe},
		&daily.Pair{Universe: a.Universe},
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
		&fun.Who{Universe: a.Universe},
		&fun.List{Universe: a.Universe},
		&fun.Top{Universe: a.Universe},
		&fun.Clock{},
		&fun.Date{},
		&fun.TurnOn{Universe: a.Universe},
		&fun.TurnOff{Universe: a.Universe},
		&fun.NewYear{},
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
		&pictures.Hello{Path: assetPath("hello.json")},
		&pictures.Anime{},
		&pictures.Furry{},
		&pictures.Flag{},
		&pictures.Car{},
		&pictures.Soy{},
		&pictures.Danbooru{API: a.Danbooru},
		&pictures.Fap{API: a.Danbooru},
		&pictures.Masyunya{},
		&pictures.Poppy{},
		&pictures.Sima{},
		&pictures.Cat{},
	}
}

func (a *App) commandHandlers() []Handler {
	return []Handler{
		&command.Add{Universe: a.Universe},
		&command.Remove{Universe: a.Universe},
		&command.Use{Universe: a.Universe},
	}
}

func (a *App) otherHandlers() []Handler {
	return []Handler{
		&handlers.Pass{},
	}
}

func assetPath(s string) string {
	return filepath.Join(assetsDirectory, s)
}
