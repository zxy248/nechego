package middleware

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/zxy248/nechego/data"
	tele "gopkg.in/zxy248/telebot.v3"
)

type UpdateDaily struct {
	Queries *data.Queries
}

func (m *UpdateDaily) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		ctx := context.Background()
		chat, err := m.Queries.GetChat(ctx, c.Chat().ID)
		if err != nil {
			return err
		}
		if chat.Data.UpdatedAt.YearDay() != time.Now().YearDay() {
			users, err := m.Queries.ListUsers(ctx, c.Chat().ID)
			if err != nil {
				return err
			}
			arg := data.UpdateDailyParams{
				ID: c.Chat().ID,
				Data: data.ChatData{
					Admin:     users[rand.N(len(users))].ID,
					Eblan:     users[rand.N(len(users))].ID,
					Pair1:     users[rand.N(len(users))].ID,
					Pair2:     users[rand.N(len(users))].ID,
					Active:    chat.Data.Active,
					UpdatedAt: time.Now(),
				},
			}
			if err := m.Queries.UpdateDaily(ctx, arg); err != nil {
				return err
			}
		}
		return next(c)
	}
}
