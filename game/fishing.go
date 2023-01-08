package game

import (
	"fmt"
	"math/rand"
)

type FishingRodType int

const (
	Spinning FishingRodType = iota
	Floating
)

type FishingRod struct {
	Type       FishingRodType
	Quality    float64 // from 0 to 1
	Durability float64 // from 0 to 1
}

func (f FishingRod) String() string {
	quality := []string{"i", "ii", "iii", "iv", "v", "I", "II", "III", "IV", "V"}
	q := int(f.Quality * float64(len(quality)))
	d := f.Durability * 100
	var t string
	switch x := f.Type; x {
	case Spinning:
		t = "спиннинг"
	case Floating:
		t = "поплавочная"
	default:
		panic(fmt.Errorf("unexpected fishing rod type %v", x))
	}
	return fmt.Sprintf("🎣 Удочка (%s) [%s, %.f%%]", t, quality[q], d)
}

func NewFishingRod(t FishingRodType) *FishingRod {
	f := &FishingRod{
		Type:       t,
		Quality:    rand.NormFloat64()*0.2 + 0.5,
		Durability: rand.Float64()*0.2 + 0.8,
	}
	if f.Quality < 0 || f.Quality > 1 {
		return NewFishingRod(t)
	}
	return f
}

func (u *User) FishingRod() (f *FishingRod, ok bool) {
	for _, v := range u.Items() {
		switch f := v.Value.(type) {
		case *FishingRod:
			return f, true
		}
	}
	return nil, false
}
