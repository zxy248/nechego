package fishing

import (
	"fmt"
	"math/rand"
	"nechego/modifier"
)

// Rod represents a fishing rod.
type Rod struct {
	Level      int     // Level of fishing effectivity.
	Durability float64 // Durability in the range [0, 1].
}

// NewRod returns a new Rod with random quality and random durability.
func NewRod() *Rod {
	return &Rod{1, randomDurability()}
}

func (r Rod) String() string {
	return fmt.Sprintf("🎣 Удочка (%d ур., %.f%%)",
		r.Level, 100*r.Durability)
}

// randomDurability returns a random value in the range [0.8, 1.0).
func randomDurability() float64 {
	return 0.8 + 0.2*rand.Float64()
}

// Mod implements the modifier.Moder interface.
func (r *Rod) Mod() (m *modifier.Mod, ok bool) {
	return &modifier.Mod{
		Emoji:       "🎣",
		Multiplier:  0.02 * float64(r.Level),
		Description: "Вы можете рыбачить.",
	}, true
}

// Broken returns true if the rod's durability is below zero.
func (r *Rod) Broken() bool {
	return r.Durability < 0
}
