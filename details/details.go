package details

import (
	"fmt"
	"math"
	"math/rand"
)

// Details are used for crafting and repairing items.
type Details struct {
	Count int
}

func (d *Details) String() string {
	return fmt.Sprintf("⚙️ Детали (%d)", d.Count)
}

// Spend subtracts n details from the count.
func (d *Details) Spend(n int) bool {
	if d.Count < n {
		return false
	}
	d.Count -= n
	return true
}

// Random returns a random amount of details.
func Random() *Details {
	c := int(1 + 30*math.Abs(rand.NormFloat64()))
	return &Details{c}
}
