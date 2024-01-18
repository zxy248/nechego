package middleware

import (
	"math/rand"

	tele "gopkg.in/zxy248/telebot.v3"
)

type RandomPhoto struct {
	Prob float64
}

func (m *RandomPhoto) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if rand.Float64() < m.Prob {
			go func() {
				pp, err := c.Bot().ProfilePhotosOf(c.Sender())
				if err != nil || len(pp) == 0 {
					return
				}
				c.Send(&pp[0])
			}()
		}
		return next(c)
	}
}
