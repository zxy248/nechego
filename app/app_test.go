package app

import (
	"nechego/model/mock"

	tele "gopkg.in/telebot.v3"
)

func newTestApp() *App {
	b, err := tele.NewBot(tele.Settings{Offline: true})
	if err != nil {
		panic(err)
	}
	m := mock.NewModel()
	return NewApp(b, m)
}
