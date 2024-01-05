package main

import (
	"fmt"
	"log"
	"nechego/bot/middleware"
	"nechego/danbooru"
	"nechego/game"
	"nechego/handlers"
	"nechego/handlers/command"
	"nechego/handlers/daily"
	"nechego/handlers/fun"
	"nechego/handlers/pictures"
	"os"
	"path/filepath"
	"time"

	tele "gopkg.in/telebot.v3"
)

var (
	botToken         = getenv("NECHEGO_TOKEN")
	assetsDirectory  = getenv("NECHEGO_ASSETS")
	storageDirectory = getenv("NECHEGO_STORAGE")
)

func main() {
	bot, err := tele.NewBot(tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal("cannot build bot: ", err)
	}
	app := &app{
		universe: game.NewUniverse(storageDirectory),
		danbooru: danbooru.New(danbooru.URL, 5*time.Second, 3),
	}
	srv := &Server{
		Bot:      bot,
		Handlers: app.handlers(),
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
	universe *game.Universe
	danbooru *danbooru.Danbooru
}

func (a *app) shutdown() error {
	return a.universe.SaveAll()
}

func (a *app) handlers() []Handler {
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

func (a *app) middleware() []Wrapper {
	return []Wrapper{
		&middleware.Recover{},
		&middleware.RandomPhoto{Prob: 0.005},
		&middleware.IgnoreWorldInactive{Universe: a.universe, Immune: turnOnMatch},
		&middleware.Throttle{Duration: 400 * time.Millisecond},
		&middleware.LogMessage{Wait: 5 * time.Second},
		&middleware.IgnoreMessageForwarded{},
		&middleware.RequireSupergroup{},
		&middleware.AddUser{Universe: a.universe},
	}
}

func turnOnMatch(c tele.Context) bool {
	var h fun.TurnOn
	return h.Match(c)
}

func (a *app) dailyHandlers() []Handler {
	return []Handler{
		&daily.Eblan{Universe: a.universe},
		&daily.Admin{Universe: a.universe},
		&daily.Pair{Universe: a.universe},
	}
}

func (a *app) funHandlers() []Handler {
	return []Handler{
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
		&fun.NewYear{},
	}
}

func (a *app) pictureHandlers() []Handler {
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
		&pictures.Danbooru{API: a.danbooru},
		&pictures.Fap{API: a.danbooru},
		&pictures.Masyunya{},
		&pictures.Poppy{},
		&pictures.Sima{},
		&pictures.Cat{},
	}
}

func (a *app) commandHandlers() []Handler {
	return []Handler{
		&command.Add{Universe: a.universe},
		&command.Remove{Universe: a.universe},
		&command.Use{Universe: a.universe},
	}
}

func (a *app) otherHandlers() []Handler {
	return []Handler{
		&handlers.Pass{},
	}
}
