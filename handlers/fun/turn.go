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
	arg := data.SetChatStatusParams{ID: c.Chat().ID, Active: true}
	if err := h.Queries.SetChatStatus(ctx, arg); err != nil {
		return err
	}

	e := [...]string{"🔈", "🔔", "✅", "🆗", "▶️"}
	return c.Send(e[rand.N(len(e))] + " Робот включен.")
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
	arg := data.SetChatStatusParams{ID: c.Chat().ID, Active: false}
	if err := h.Queries.SetChatStatus(ctx, arg); err != nil {
		return err
	}

	e := [...]string{"🔇", "🔕", "💤", "❌", "⛔️", "🚫", "⏹"}
	return c.Send(e[rand.N(len(e))] + " Робот выключен.")
}
