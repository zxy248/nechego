package app

import (
	"nechego/model"

	tele "gopkg.in/telebot.v3"
)

const dataPath = "data"

type App struct {
	bot   *tele.Bot
	model *model.Model
}

func NewApp(b *tele.Bot, m *model.Model) *App {
	return &App{b, m}
}

func (a *App) Start() {
	a.bot.Handle(tele.OnText, a.route, a.preprocess)
	a.bot.Handle(tele.OnUserJoined, a.handleJoin, a.preprocess)
	a.bot.Start()
}
