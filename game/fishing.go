package game

import (
	"fmt"
	"math/rand"
	"nechego/fishing"
)

type FishingRod struct {
	Quality    float64 // from 0 to 1
	Durability float64 // from 0 to 1
}

func (f FishingRod) String() string {
	lvls := [...]string{"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X"}
	lvl := lvls[int(f.Quality*float64(len(lvls)))]
	dur := f.Durability * 100
	return fmt.Sprintf("🎣 Удочка (%s, %.f%%)", lvl, dur)
}

func NewFishingRod() *FishingRod {
	f := &FishingRod{
		Quality:    rand.NormFloat64()*0.2 + 0.5,
		Durability: rand.Float64()*0.2 + 0.8,
	}
	if f.Quality < 0 || f.Quality > 1 {
		return NewFishingRod()
	}
	return f
}

func (u *User) FishingRod() (f *FishingRod, ok bool) {
	for _, v := range u.Inventory.normalize() {
		switch f := v.Value.(type) {
		case *FishingRod:
			return f, true
		}
	}
	return nil, false
}

func (u *User) Fish(rod *FishingRod) *Item {
	rod.Durability -= 0.01

	if rand.Float64() < 0.08 {
		return randomItem()
	}
	f := fishing.RandomFish()
	q := rod.Quality*0.5 + 1.0
	f.Length *= q
	f.Weight *= q
	return &Item{Type: ItemTypeFish, Transferable: true, Value: f}
}
