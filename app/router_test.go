package app

import (
	"nechego/input"
	"nechego/statistics"
	"testing"

	"fmt"
	"nechego/model"
	"nechego/service"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func TestCommandHandler(t *testing.T) {
	a := testApp()
	for _, c := range input.AllCommands() {
		if c == input.CommandUnknown {
			if h := a.commandHandler(c); h != nil {
				t.Error("want nil")
			}
		} else {
			if h := a.commandHandler(c); h == nil {
				t.Error("want not nil")
			}
		}
	}
}

func testApp() *App {
	b := testBot()
	m := testModel()
	l := testLog()
	t := testStatistics(m)
	s := testService(m, t)
	p := Preferences{}
	return NewApp(b, l, t, s, p)
}

func testBot() *tele.Bot {
	b, err := tele.NewBot(tele.Settings{Offline: true})
	failOn(err)
	return b
}

func testModel() *model.Model {
	db := sqlx.MustOpen("sqlite3", ":memory:")
	err := db.Ping()
	failOn(err)
	return model.New(db)
}

func testService(m *model.Model, s *statistics.Statistics) *service.Service {
	return service.New(m, s, service.Config{})
}

func testStatistics(m *model.Model) *statistics.Statistics {
	return statistics.New(m, statistics.Settings{})
}

func testLog() *zap.Logger {
	l, err := zap.NewDevelopment()
	failOn(err)
	return l
}

func testConfig() service.Config {
	return service.Config{}
}

func failOn(err error) {
	if err != nil {
		panic(fmt.Errorf("fail on test init: %v", err))
	}
}
