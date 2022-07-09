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

// Start starts the bot.
func (a *App) Start() {
	go a.restoreEnergyEvery(restoreEnergyCooldown)
	a.bot.Handle(tele.OnText, a.route, a.pipeline)
	a.bot.Handle(tele.OnDice, a.handleRoll, a.pipeline)
	a.bot.Handle(tele.OnUserJoined, a.handleJoin, a.pipeline)

	a.SugarLog().Info("The bot has started.")
	a.bot.Start()
}

// SugarLog is a shorthand for a SugaredLogger.
func (a *App) SugarLog() *zap.SugaredLogger {
	return a.log.Sugar()
}
