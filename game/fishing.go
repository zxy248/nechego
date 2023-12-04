package game

import (
	"math/rand"
	"nechego/fishing"
	"nechego/item"
)

// FishCatchProb returns the user's chance to catch fish.
func (u *User) FishCatchProb() float64 {
	p := 0.5
	p += -0.04 + 0.08*u.Luck()
	return p
}

// Fish returns a new random item to be added to the user's inventory
// and decreases durability of the fishing rod r. If the user was not
// able to catch any fish, returns (nil, false).
func (u *User) Fish(r *fishing.Rod) (i *item.Item, ok bool) {
	r.Durability -= 0.01
	if rand.Float64() > u.FishCatchProb() {
		return nil, false
	}
	if rand.Float64() < 0.08 {
		return item.Random(), true
	}
	f := fishing.RandomFish()
	quality := 1 + 0.1*float64(r.Level)
	luck := 0.9 + 0.2*u.Luck()
	multiplier := quality * luck
	f.Length *= multiplier
	f.Weight *= multiplier
	return item.New(f), true
}
