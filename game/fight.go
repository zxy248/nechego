package game

import (
	"math"
	"math/rand"
)

func (u *User) Strength() float64 {
	const base, mult = 1, 2
	c := base + u.Activity + u.Luck() + ModSum(u.Mods())
	return mult * c
}

func Fight(u1, u2 *User) (uw, ul *User, dr float64) {
	s1 := u1.Strength() * rand.Float64()
	s2 := u2.Strength() * rand.Float64()
	if s1 > s2 {
		uw, ul = u1, u2
	} else {
		uw, ul = u2, u1
	}
	dr = eloDelta(uw.Rating, ul.Rating, kDefault, scoreWin)
	uw.Rating += dr
	ul.Rating -= dr
	return
}

const (
	scoreWin = 1.0
	kDefault = 20.0
)

func eloDelta(a, b, k, score float64) float64 {
	x := 1.0 / (1.0 + math.Pow(10.0, (b-a)/400.0))
	return k * (score - x)
}
