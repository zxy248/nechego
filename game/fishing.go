package game

import (
	"math/rand"
	"nechego/fishing"
	"nechego/item"
)

func (u *User) FishingRod() (r *fishing.Rod, ok bool) {
	for _, v := range u.Inventory.Normal() {
		switch r := v.Value.(type) {
		case *fishing.Rod:
			return r, true
		}
	}
	return nil, false
}

func (u *User) Fish(r *fishing.Rod) *item.Item {
	r.Durability -= 0.01
	if rand.Float64() < 0.08 {
		return item.Random()
	}

	quality := 1.0 + 0.5*r.Quality
	luck := 0.9 + 0.2*u.Luck()
	total := quality * luck

	f := fishing.RandomFish()
	f.Length *= total
	f.Weight *= total
	return &item.Item{Type: item.TypeFish, Transferable: true, Value: f}
}
