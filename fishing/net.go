package fishing

import (
	"fmt"
	"math/rand"
)

// Net represents a fishing net.
type Net struct {
	Fish       []*Fish // Fish caught in the net.
	Durability float64 // Durability in the range [0, 1].
	Capacity   int     // Maximum number of fish in the net.
}

// NewNet creates an empty fishing net.
func NewNet() *Net {
	return &Net{
		Fish:       []*Fish{},
		Durability: 1.0,
		Capacity:   5 + rand.Intn(11), // [5, 15]
	}
}

// Count returns the number of fish caught in the fishing net.
func (n *Net) Count() int {
	return len(n.Fish)
}

// Broken returns true if the fishing net's durability is below 0.
func (n *Net) Broken() bool {
	return n.Durability < 0
}

// NetDurabilityDecrement is substracted from the fishing net's
// durability every time it fills with fish.
const NetDurabilityDecrement = 0.01

// Fill adds a random fish to the fishing net.
func (n *Net) Fill() {
	if len(n.Fish) < n.Capacity {
		n.Durability -= 0.01
		n.Fish = append(n.Fish, RandomFish())
	}
}

// Unload returns the list of caught fish and empties the fishing net.
func (n *Net) Unload() (catch []*Fish) {
	catch = n.Fish
	n.Fish = []*Fish{}
	return
}

// String returns the textual representation of the fishing net.
func (n *Net) String() string {
	return fmt.Sprintf("🕸 Рыболовная сеть (%d/%d, %.f%%)",
		n.Count(), n.Capacity, 100*n.Durability)
}
