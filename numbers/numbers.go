package numbers

import (
	crand "crypto/rand"
	"encoding/binary"
	"errors"
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

func UniNormal(min, mid, max float64) float64 {
	j := max - mid
	n := Normal(-j, j)
	if n < mid {
		return rand.Float64()*(mid-min) + min
	}
	return mid + n
}

func InRange(min, max int) int {
	return rand.Intn(max-min+1) + min
}

type cryptSource struct{}

func (c cryptSource) Int63() int64 {
	return int64(c.Uint64() & ^uint64(1<<63))
}

func (c cryptSource) Uint64() uint64 {
	var n uint64
	if err := binary.Read(crand.Reader, binary.LittleEndian, &n); err != nil {
		panic(errors.New("cannot read a random number"))
	}
	return n
}

func (c cryptSource) Seed(_ int64) {}

var CryptSource rand.Source64 = cryptSource{}
