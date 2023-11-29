package fun

import (
	"math/rand"
	"nechego/game"
	"nechego/handlers"
	tu "nechego/teleutil"

	tele "gopkg.in/telebot.v3"
)

type TurnOn struct {
	Universe *game.Universe
}

var turnOnRe = handlers.Regexp("^!(вкл|подкл|подруб)")

func (h *TurnOn) Match(c tele.Context) bool {
	return turnOnRe.MatchString(c.Text())
}

func (h *TurnOn) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	world.Inactive = false
	es := [...]string{"🔈", "🔔", "✅", "🆗", "▶️"}
	return c.Send(es[rand.Intn(len(es))] + " Робот включен.")
}

type TurnOff struct {
	Universe *game.Universe
}

var turnOffRe = handlers.Regexp("^!(выкл|откл|отруб)")

func (h *TurnOff) Match(c tele.Context) bool {
	return turnOffRe.MatchString(c.Text())
}

func (h *TurnOff) Handle(c tele.Context) error {
	world, _ := tu.Lock(c, h.Universe)
	defer world.Unlock()

	world.Inactive = true
	e := [...]string{"🔇", "🔕", "💤", "❌", "⛔️", "🚫", "⏹"}
	return c.Send(e[rand.Intn(len(e))] + " Робот выключен.")
}
