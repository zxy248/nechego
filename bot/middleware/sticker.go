package middleware

import (
	"math/rand/v2"

	"github.com/zxy248/nechego/game"
	tu "github.com/zxy248/nechego/teleutil"

	tele "gopkg.in/zxy248/telebot.v3"
)

type RandomSticker struct {
	Universe *game.Universe
	Prob     float64
}

func (m *RandomSticker) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if rand.Float64() < m.Prob {
			w := tu.Lock(c, m.Universe)
			if len(w.Stickers) > 0 {
				var s tele.Sticker
				s.FileID = w.Stickers[rand.N(len(w.Stickers))]
				c.Send(&s)
			}
			w.Unlock()
		}
		return next(c)
	}
}
