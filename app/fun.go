package app

import (
	"math/rand"

	tele "gopkg.in/telebot.v3"
)

var mouseVideo = "mouse.mp4"

// !мыш
func (a *App) handleMouse(c tele.Context) error {
	return respondVideo(c, a.Locate(mouseVideo))
}

var tikTokVideo = "tiktok.mp4"

// !тикток
func (a *App) handleTikTok(c tele.Context) error {
	return respondVideo(c, a.Locate(tikTokVideo))
}

// !игра
func (a *App) handleGame(c tele.Context) error {
	game := games[rand.Intn(len(games))]
	return c.Send(game)
}

var games = []*tele.Dice{tele.Dart, tele.Ball, tele.Goal, tele.Slot, tele.Bowl}

func randomGame() *tele.Dice {
	return games[rand.Intn(len(games))]
}
