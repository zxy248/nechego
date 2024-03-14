package middleware

import (
	"context"

	"github.com/zxy248/nechego/data"
	tu "github.com/zxy248/nechego/teleutil"

	tele "gopkg.in/zxy248/telebot.v3"
)

type IgnoreMessageForwarded struct{}

func (m *IgnoreMessageForwarded) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if tu.MessageForwarded(c.Message()) {
			return nil
		}
		return next(c)
	}
}

type IgnoreInactive struct {
	Queries *data.Queries
	Immune  func(tele.Context) bool
}

func (m *IgnoreInactive) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		ctx := context.Background()
		chat, err := m.Queries.GetChat(ctx, c.Chat().ID)
		if err != nil {
			return err
		}
		if !chat.Data.Active && !m.Immune(c) {
			return nil
		}
		return next(c)
	}
}
