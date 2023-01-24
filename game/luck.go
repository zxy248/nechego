package game

import (
	"encoding/binary"
	"hash/fnv"
	"math"
	"nechego/dates"
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
