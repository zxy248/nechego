package fun

import (
	"github.com/zxy248/nechego/handlers"
	"math/rand"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Game struct{}

var gameRe = handlers.NewRegexp("^!игр")

func (h *Game) Match(c tele.Context) bool {
	return gameRe.MatchString(c.Text())
}

func (h *Game) Handle(c tele.Context) error {
	games := [...]*tele.Dice{tele.Dart, tele.Ball, tele.Goal, tele.Slot, tele.Bowl}
	return c.Send(games[rand.Intn(len(games))])
}
