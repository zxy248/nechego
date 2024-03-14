package middleware

import (
	"context"

	"github.com/zxy248/nechego/data"
	"github.com/zxy248/nechego/handlers"
	tele "gopkg.in/zxy248/telebot.v3"
)

type LogMessage struct {
	Queries *data.Queries
}

func (m *LogMessage) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		ctx := context.Background()
		arg := data.AddMessageParams{
			UserID:  c.Sender().ID,
			ChatID:  c.Chat().ID,
			Content: c.Message().Text,
		}
		id, err := m.Queries.AddMessage(ctx, arg)
		if err != nil {
			return err
		}
		c.Set(handlers.MessageIDKey, id)
		return next(c)
	}
}
