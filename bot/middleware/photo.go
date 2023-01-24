package middleware

import (
	"math/rand"
	"nechego/avatar"

	tele "gopkg.in/telebot.v3"
)

type RandomPhoto struct {
	Avatars *avatar.Storage
}

func (m *RandomPhoto) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if rand.Float64() < 0.02 {
			var r []*tele.Photo
			p, err := c.Bot().ProfilePhotosOf(c.Sender())
			if err != nil {
				return err
			}
			if len(p) > 0 {
				r = append(r, &p[0])
			}
			if a, ok := m.Avatars.Get(c.Sender().ID); ok {
				r = append(r, a)
			}
			if len(r) > 0 {
				c.Send(r[rand.Intn(len(r))])
			}
		}
		return next(c)
	}
}
