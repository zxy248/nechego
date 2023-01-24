package game

import (
	"math/rand"
	"nechego/elo"
)

func (u *User) Strength(w *World) float64 {
	const (
		base = 2
		mult = 5
	)
	return mult * (base + u.Modset(w).Sum() + u.Luck())
}

func (w *World) Fight(a, b *User) (winner, loser *User, rating float64) {
	if a == b {
		panic("user cannot be an opponent to themself")
	}

	f := func(u *User) float64 { return u.Strength(w) * rand.Float64() }
	if f(a) > f(b) {
		winner, loser = a, b
	} else {
		winner, loser = b, a
	}

	rating = elo.EloDelta(winner.Rating, loser.Rating, elo.KDefault, elo.ScoreWin)
	winner.Rating += rating
	loser.Rating -= rating
	return
}
