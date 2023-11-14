package game

import (
	"fmt"
)

// Energy represents the user's energy level.
// It must be in the range [0, 1].
type Energy float64

// Spend subtracts x energy and returns true on success.
// Returns false if the energy level would drop below zero.
func (e *Energy) Spend(x Energy) bool {
	if x < 0 {
		panic(fmt.Sprintf("cannot spend %v energy", x))
	}
	if *e < x {
		return false
	}
	*e -= x
	return true
}

// Add restores x energy.
// Clamps at the upper bound of the range.
func (e *Energy) Add(x Energy) {
	if x < 0 {
		panic(fmt.Sprintf("cannot add %v energy", x))
	}
	*e += x
	if *e > 1 {
		*e = 1
	}
}

// Low returns true if the energy level is close to 0.
func (e *Energy) Low() bool {
	if *e < 0.1 {
		return true
	}
	return false
}

// Full returns true if the energy level is close to 1.
func (e *Energy) Full() bool {
	if *e > 0.97 {
		return true
	}
	return false
}
