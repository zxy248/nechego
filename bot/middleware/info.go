package middleware

import (
	"context"

	"github.com/zxy248/nechego/data"
	tu "github.com/zxy248/nechego/teleutil"
	tele "gopkg.in/zxy248/telebot.v3"
)

type UpdateInfo struct {
	Queries *data.Queries
}

func (m *UpdateInfo) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		ctx := context.Background()
		if err := m.Queries.UpdateUser(ctx, data.UpdateUserParams{
			ID:        c.Sender().ID,
			FirstName: c.Sender().FirstName,
			LastName:  c.Sender().LastName,
			Username:  c.Sender().Username,
			IsPremium: c.Sender().IsPremium,
		}); err != nil {
			return err
		}
		if err := m.Queries.UpdateChat(ctx, data.UpdateChatParams{
			ID:    c.Chat().ID,
			Title: c.Chat().Title,
		}); err != nil {
			return err
		}
		if err := m.Queries.UpdateChatMember(ctx, data.UpdateChatMemberParams{
			UserID:      c.Sender().ID,
			ChatID:      c.Chat().ID,
			CustomTitle: tu.Member(c, c.Sender()).Title,
		}); err != nil {
			return err
		}
		return next(c)
	}
}
