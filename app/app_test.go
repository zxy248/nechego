package app

import (
	"fmt"
	"nechego/model"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

// newTestApp returns a new App for testing purposes.
func newTestApp() *App {
	bot, err := tele.NewBot(tele.Settings{Offline: true})
	failOn(err)
	mod := model.NewModel(sqlx.MustOpen("sqlite3", ":memory:"))
	logger, err := zap.NewDevelopment()
	failOn(err)
	defer logger.Sync()
	return NewApp(bot, mod, logger)
}

func failOn(err error) {
	if err != nil {
		panic(fmt.Errorf("fail on test init: %v", err))
	}
}
