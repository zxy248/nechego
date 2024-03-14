package fun

import (
	"context"
	"math/rand/v2"

	"github.com/zxy248/nechego/data"
	"github.com/zxy248/nechego/handlers"

	tele "gopkg.in/zxy248/telebot.v3"
)

type TurnOn struct {
	Queries *data.Queries
}

var turnOnRe = handlers.NewRegexp("^!(вкл|подкл|подруб)")

func (h *TurnOn) Match(c tele.Context) bool {
	return turnOnRe.MatchString(c.Text())
}

func (h *TurnOn) Handle(c tele.Context) error {
	ctx := context.Background()
	if err := h.Queries.ActivateChat(ctx, c.Chat().ID); err != nil {
		return err
	}
	es := [...]string{"🔈", "🔔", "✅", "🆗", "▶️"}
	e := es[rand.N(len(es))]
	return c.Send(e + " Робот включен.")
}

type TurnOff struct {
	Queries *data.Queries
}

var turnOffRe = handlers.NewRegexp("^!(выкл|откл|отруб)")

func (h *TurnOff) Match(c tele.Context) bool {
	return turnOffRe.MatchString(c.Text())
}

func (h *TurnOff) Handle(c tele.Context) error {
	ctx := context.Background()
	if err := h.Queries.DisableChat(ctx, c.Chat().ID); err != nil {
		return err
	}
	es := [...]string{"🔇", "🔕", "💤", "❌", "⛔️", "🚫", "⏹"}
	e := es[rand.N(len(es))]
	return c.Send(e + " Робот выключен.")
}
