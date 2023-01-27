package game

import (
	"encoding/binary"
	"hash/fnv"
	"math"
	"nechego/dates"
	"nechego/modifier"
	"time"
)

// Luck is a dynamic luck modifier.
// If the luck value is low, the resulting modifier is good.
// If the luck value is high, the resulting modifier is bad.
// If the luck value is average, the resulting modifier is nil.
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
			Description: "Вам не везет.",
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
			Description: "Вам везет.",
		}, true
	}
	return nil, false
}

func (u *User) Luck() float64 {
	return luck(dates.Today(), u.TUID)
}

func luck(t time.Time, id int64) float64 {
	const magic = 497611803913981554
	return checksum(t.UnixNano(), id, magic)
}

func checksum(x ...any) float64 {
	h := fnv.New32()
	for _, v := range x {
		binary.Write(h, binary.LittleEndian, v)
	}
	return float64(h.Sum32()) / math.MaxUint32
}
