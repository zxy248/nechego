package app

import (
	"nechego/model/mock"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

// newTestApp returns a new App for testing purposes.
func newTestApp() *App {
	bot, err := tele.NewBot(tele.Settings{Offline: true})
	if err != nil {
		panic(err)
	}
	model := mock.NewModel()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	return NewApp(bot, model, logger)
}
