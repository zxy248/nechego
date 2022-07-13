package numbers

import (
	"math/rand"
)

func Normal(min, max float64) float64 {
	dev := (max - min) / 6
	mean := (min + max) / 2
	r := rand.NormFloat64()*dev + mean
	if r < min || r > max {
		return Normal(min, max)
	}
	return r
}

func RandMidNormal(min, mid, max float64) float64 {
	j := max - mid
	n := Normal(-j, j)
	if n < mid {
		return rand.Float64()*(mid-min) + min
	}
	return mid + n
}
