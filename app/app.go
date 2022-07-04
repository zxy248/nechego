package app

import (
	"nechego/model"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

const dataPath = "data"

type App struct {
	bot   *tele.Bot
	model *model.Model
	log   *zap.Logger
}

// NewApp returns a new app.
func NewApp(b *tele.Bot, m *model.Model, l *zap.Logger) *App {
	return &App{b, m, l}
}

// Start starts the bot, routing all preprocessed text messages to appropriate handlers.
func (a *App) Start() {
	go a.restoreEnergyEvery(restoreEnergyCooldown)
	a.bot.Handle(tele.OnText, a.route, a.preprocess, a.logMessage)
	a.bot.Handle(tele.OnUserJoined, a.handleJoin, a.preprocess)
	a.bot.Start()
}

// sugar is a shorthand for a.log.Sugar.
func (a *App) sugar() *zap.SugaredLogger {
	return a.log.Sugar()
}
