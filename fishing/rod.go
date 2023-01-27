package fishing

import (
	"fmt"
	"math/rand"
	"nechego/modifier"
)

// Rod represents a fishing rod.
type Rod struct {
	Quality    float64
	Durability float64
}

// NewRod returns a Rod with random quality and random durability.
func NewRod() *Rod {
	return &Rod{randomQuality(), randomDurability()}
}

func (r Rod) String() string {
	return fmt.Sprintf("🎣 Удочка (%.1f, %.f%%)",
		10*r.Quality, 100*r.Durability)
}

// randomQuality returns a normally distributed value
// (mean = 0.5, stddev = 0.2) in the range [0, 1].
func randomQuality() float64 {
	q := 0.5 + 0.2*rand.NormFloat64()
	if q < 0 || q > 1 {
		return randomQuality()
	}
	return q
}

// randomDurability returns a random value in the range [0.8, 1.0).
func randomDurability() float64 {
	return 0.8 + 0.2*rand.Float64()
}

func (r *Rod) Mod() (m *modifier.Mod, ok bool) {
	return &modifier.Mod{
		Emoji:       "🎣",
		Multiplier:  +0.05,
		Description: "Вы можете рыбачить.",
	}, true
}
