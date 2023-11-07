package game

import (
	"math/rand"
	"nechego/elo"
)

func (u *User) Strength() float64 {
	const base, mult = 1, 2
	c := base + u.Activity + u.Luck() + u.Modifiers().Sum()
	return mult * c
}

func Fight(u1, u2 *User) (winner, loser *User, rating float64) {
	if u1 == u2 {
		panic("user cannot be an opponent to themself")
	}

	s1 := u1.Strength() * rand.Float64()
	s2 := u2.Strength() * rand.Float64()
	if s1 > s2 {
		winner, loser = u1, u2
	} else {
		winner, loser = u2, u1
	}

	rating = elo.EloDelta(winner.Rating, loser.Rating, elo.KDefault, elo.ScoreWin)
	winner.Rating += rating
	loser.Rating -= rating
	return
}
