package farm

import (
	"fmt"
	"math/rand"
)

// Fertilizer can be applied to the farm.
type Fertilizer struct {
	Volume int
}

// NewFertilizer returns a new fertilizer of a non-zero volume.
func NewFertilizer() *Fertilizer {
	return &Fertilizer{1 + rand.Intn(100)}
}

func (f Fertilizer) String() string {
	return fmt.Sprintf("🛢 Удобрение (%d л.)", f.Volume)
}
