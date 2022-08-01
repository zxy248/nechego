package app

import (
	"nechego/numbers"
	"nechego/service"
	"nechego/statistics"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

type Preferences struct {
	DataPath     string
	EnergyPeriod time.Duration
	ListLength   numbers.Interval
	HelloChance  float64
}

type App struct {
	bot      *tele.Bot
	log      *zap.Logger
	stat     *statistics.Statistics
	service  *service.Service
	stickers *Stickers
	pref     Preferences
}

func New(
	b *tele.Bot,
	l *zap.Logger,
	t *statistics.Statistics,
	s *service.Service,
	p Preferences,
) *App {
	a := &App{
		bot:     b,
		log:     l,
		service: s,
		stat:    t,
		pref:    p,
	}
	return a
}

func (a *App) Start() {
	a.bot.Handle(tele.OnText, a.route, a.pipeline)
	a.bot.Handle(tele.OnDice, a.handleRoll, a.pipeline)
	a.bot.Handle(tele.OnUserJoined, a.handleJoin, a.pipeline)
	handleEnergy = a.energyHandler()

	a.log.Sugar().Info("The bot has started.")
	a.bot.Start()
}

func (a *App) Locate(path string) string {
	return filepath.Join(a.pref.DataPath, path)
}
