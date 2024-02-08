package fun

import (
	"math/rand/v2"

	"github.com/zxy248/nechego/handlers"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Game struct{}

var gameRe = handlers.NewRegexp("^!игр")

func (h *Game) Match(c tele.Context) bool {
	return gameRe.MatchString(c.Text())
}

func (h *Game) Handle(c tele.Context) error {
	games := [...]*tele.Dice{tele.Dart, tele.Ball, tele.Goal, tele.Slot, tele.Bowl}
	return c.Send(games[rand.N(len(games))])
}
