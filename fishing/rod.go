package fishing

import (
	"fmt"
	"math/rand"
)

type Quality float64

func RandomQuality() Quality {
	q := 0.5 + 0.2*rand.NormFloat64()
	if q < 0 || q > 1 {
		return RandomQuality()
	}
	return Quality(q)
}

func (q Quality) String() string {
	return fmt.Sprintf("%.1f", 10*q)
}

type Durability float64

func RandomDurability() Durability {
	return Durability(0.8 + 0.2*rand.Float64())
}

func (d Durability) String() string {
	return fmt.Sprintf("%.f%%", 100*d)
}

type Rod struct {
	Quality    Quality
	Durability Durability
}

func NewRod() *Rod {
	return &Rod{RandomQuality(), RandomDurability()}
}

func (r Rod) String() string {
	return fmt.Sprintf("🎣 Удочка (%v, %v)", r.Quality, r.Durability)
}
