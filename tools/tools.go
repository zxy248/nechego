package tools

import (
	"fmt"
	"math/rand"
)

// Knife is used for crafting.
type Knife struct {
	Durability float64
}

// NewKnife returns a Knife of random durability.
func NewKnife() *Knife {
	return &Knife{Durability: 0.8 + 0.2*rand.Float64()}
}

func (k Knife) String() string {
	return fmt.Sprintf("🔪 Нож (%.f%%)", k.Durability*100)
}
