package numbers

import (
	crand "crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"math/rand"

	"golang.org/x/exp/constraints"
)

func Normal(min, max float64) float64 {
	sigma := (max - min) / 6
	mu := (min + max) / 2
	r := rand.NormFloat64()*sigma + mu
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

func Min[T constraints.Ordered](x T, n ...T) T {
	for _, v := range n {
		if v < x {
			x = v
		}
	}
	return x
}

func Max[T constraints.Ordered](x T, n ...T) T {
	for _, v := range n {
		if v > x {
			x = v
		}
	}
	return x
}

type Interval struct {
	a, b int
}

func MakeInterval(a, b int) Interval {
	return Interval{a, b}
}

func (i Interval) Min() int {
	return Min(i.a, i.b)
}

func (i Interval) Max() int {
	return Max(i.a, i.b)
}

func (i Interval) Random() int {
	return InRange(i.Min(), i.Max())
}

func LuckyByte(data []byte) byte {
	return sha1.Sum(data)[0]
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
