package game

import (
	"fmt"
	"math/rand"
	"nechego/elo"
	"nechego/modifier"
	"nechego/pets"
)

func (u *User) Strength() float64 {
	return 10 * (1.0 + u.Modset().Sum())
}

func (u *User) Fight(opponent *User) (win, lose *User, r float64) {
	if u == opponent {
		panic("user cannot be an opponent to themself")
	}
	if u.power() > opponent.power() {
		win, lose = u, opponent
	} else {
		win, lose = opponent, u
	}
	r = elo.EloDelta(win.Rating, lose.Rating, elo.KDefault, elo.ScoreWin)
	win.Rating += r
	lose.Rating -= r
	return
}

func (u *User) power() float64 {
	return (5*u.Luck() + u.Strength()) * rand.Float64()
}

func (u *User) Modset() modifier.Set {
	set := modifier.Set{}
	if u.Admin() {
		set.Add(modifier.Admin)
	}
	if u.Eblan() {
		set.Add(modifier.Eblan)
	}
	if u.Energy == 0 {
		set.Add(modifier.NoEnergy)
	}
	if u.Energy == EnergyCap {
		set.Add(modifier.FullEnergy)
	}
	if u.Energy > EnergyCap {
		set.Add(modifier.MuchEnergy)
	}
	if u.Rich() {
		set.Add(modifier.Rich)
	}
	if u.Poor() {
		set.Add(modifier.Poor)
	}
	if u.InDebt() {
		set.Add(modifier.Debtor)
	}
	if u.Inventory.Count() > InventorySize {
		set.Add(modifier.Heavy)
	}
	if l, ok := luckModifier(u.Luck()); ok {
		set.Add(l)
	}
	if _, ok := u.FishingRod(); ok {
		set.Add(modifier.Fisher)
	}
	if _, ok := u.Phone(); ok {
		set.Add(modifier.Phone)
	}
	if p, ok := u.Pet(); ok {
		q := 0.05
		switch p.Species.Quality() {
		case pets.Rare:
			q = 0.10
		case pets.Exotic:
			q = 0.15
		case pets.Legendary:
			q = 0.20
		}
		r := ""
		if p.Species.Quality() != pets.Common {
			r = fmt.Sprintf("%s ", p.Species.Quality())
		}
		set.Add(&modifier.Mod{
			Emoji:       p.Species.Emoji(),
			Multiplier:  q,
			Description: fmt.Sprintf("У вас есть %sпитомец: <code>%s</code>", r, p),
		})
	}
	return set
}
