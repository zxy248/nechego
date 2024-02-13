package main

import (
	"path/filepath"
	"time"

	"github.com/zxy248/nechego/bot/middleware"
	"github.com/zxy248/nechego/game"
	"github.com/zxy248/nechego/handlers"
	"github.com/zxy248/nechego/handlers/command"
	"github.com/zxy248/nechego/handlers/daily"
	"github.com/zxy248/nechego/handlers/fun"
	"github.com/zxy248/nechego/handlers/pictures"
)

type App struct {
	Universe *game.Universe
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
		&middleware.RandomReact{Prob: 0.033},
		&middleware.IgnoreWorldInactive{
			Universe: a.Universe,
			Immune:   (&fun.TurnOn{}).Match,
		},
		&middleware.Log{Timeout: 10 * time.Second},
		&middleware.Throttle{Duration: 400 * time.Millisecond},
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
		&command.Add{Universe: a.Universe},
		&command.Remove{Universe: a.Universe},
		&command.Use{Universe: a.Universe},
	}
}

func (a *App) otherHandlers() []Handler {
	logger := &handlers.Logger{Dir: messagesDirectory}

	return []Handler{
		&fun.Speak{Universe: a.Universe, Logger: logger, Attempts: 50},
		&handlers.Pass{logger},
	}
}

func assetPath(s string) string {
	return filepath.Join(assetsDirectory, s)
}
