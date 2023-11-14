package game

import (
	"math/rand"
	"time"
)

func (u *User) Luck() float64 {
	return luck(time.Now(), u.ID)
}

func luck(t time.Time, id int64) float64 {
	const magic = 2042238053
	yd := int64(t.YearDay())
	seed := magic + id + yd
	rng := rand.New(rand.NewSource(seed))
	return rng.Float64()
}
