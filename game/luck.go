package game

import (
	"math/rand"
	"nechego/modifier"
	"time"
)

// Luck is a dynamic luck modifier.
// If the luck value is high, returns a positive modifier.
// If the luck value is low, returns a negative modifier.
// If the luck value is average, returns nil, false.
type Luck float64

func (l Luck) Mod() (m *modifier.Mod, ok bool) {
	switch {
	case l < 0.05:
		return &modifier.Mod{
			Emoji:       "☠️",
			Multiplier:  -0.04,
			Description: "Вас преследуют неудачи.",
		}, true
	case l < 0.2:
		return &modifier.Mod{
			Emoji:       "🌧",
			Multiplier:  -0.02,
			Description: "Вам не везёт.",
		}, true
	case l > 0.95:
		return &modifier.Mod{
			Emoji:       "🍀",
			Multiplier:  +0.02,
			Description: "Сегодня ваш день.",
		}, true
	case l > 0.8:
		return &modifier.Mod{
			Emoji:       "🤞",
			Multiplier:  +0.04,
			Description: "Вам везёт.",
		}, true
	}
	return nil, false
}

func (u *User) Luck() float64 {
	return luck(time.Now(), u.TUID)
}

func luck(t time.Time, id int64) float64 {
	const magic = 2042238053
	yd := int64(t.YearDay())
	seed := magic + id + yd
	rng := rand.New(rand.NewSource(seed))
	return rng.Float64()
}
