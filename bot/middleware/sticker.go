package middleware

import (
	"context"
	"math/rand/v2"

	"github.com/zxy248/nechego/data"

	tele "gopkg.in/zxy248/telebot.v3"
)

type RandomSticker struct {
	Queries *data.Queries
	Prob    float64
}

func (m *RandomSticker) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if rand.Float64() < m.Prob {
			ctx := context.Background()
			stickers, err := m.Queries.RecentStickers(ctx, c.Chat().ID)
			if err != nil {
				return err
			}
			if len(stickers) > 0 {
				var s tele.Sticker
				s.FileID = stickers[rand.N(len(stickers))].FileID
				c.Send(&s)
			}
		}
		return next(c)
	}
}
