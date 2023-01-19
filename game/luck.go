package game

import (
	"encoding/binary"
	"hash/fnv"
	"math"
	"nechego/dates"
	"nechego/modifier"
	"time"
)

const luckMagic = 497611803913981554

func (u *User) Luck() float64 {
	return luck(dates.Today(), u.TUID)
}

func luck(t time.Time, id int64) float64 {
	return checksum(t.UnixNano(), id, luckMagic)
}

func checksum(x ...any) float64 {
	h := fnv.New32()
	for _, v := range x {
		binary.Write(h, binary.LittleEndian, v)
	}
	return float64(h.Sum32()) / math.MaxUint32
}

func luckModifier(l float64) (m *modifier.Mod, ok bool) {
	var x *modifier.Mod
	switch {
	case l < 0.05:
		x = modifier.TerribleLuck
	case l < 0.20:
		x = modifier.BadLuck
	case l > 0.95:
		x = modifier.ExcellentLuck
	case l > 0.80:
		x = modifier.GoodLuck
	default:
		return nil, false
	}
	return &modifier.Mod{
		Emoji:       x.Emoji,
		Multiplier:  0,
		Description: x.Description,
	}, true
}
