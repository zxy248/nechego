package fun

import (
	"math/rand"
	"nechego/game"
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type TurnOn struct {
	Universe *game.Universe
}

func MatchTurnOn(s string) bool {
	return handlers.MatchRegexp("^!(вкл|подкл|подруб)", s)
}

func (h *TurnOn) Match(s string) bool {
	return MatchTurnOn(s)
}

func (h *TurnOn) Handle(c tele.Context) error {
	return handlers.HandleWorld(c, h.Universe, h)
}

func (h *TurnOn) HandleWorld(c tele.Context, w *game.World) error {
	w.Inactive = false
	emoji := [...]string{"🔈", "🔔", "✅", "🆗", "▶️"}
	return c.Send(emoji[rand.Intn(len(emoji))])
}

type TurnOff struct {
	Universe *game.Universe
}

func (h *TurnOff) Match(s string) bool {
	return handlers.MatchRegexp("^!(выкл|откл)", s)
}

func (h *TurnOff) Handle(c tele.Context) error {
	return handlers.HandleWorld(c, h.Universe, h)
}

func (h *TurnOff) HandleWorld(c tele.Context, w *game.World) error {
	w.Inactive = true
	emoji := [...]string{"🔇", "🔕", "💤", "❌", "⛔️", "🚫", "⏹"}
	return c.Send(emoji[rand.Intn(len(emoji))])
}
