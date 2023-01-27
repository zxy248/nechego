package game

import "fmt"

const (
	EnergyCap = 100
	LowEnergy = EnergyCap / 10
)

// SpendEnergy subtracts e energy and returns true.
// If the user's energy would become negative, do nothing and return false.
func (u *User) SpendEnergy(e int) bool {
	if e < 0 {
		panic(fmt.Errorf("cannot spend %v energy", e))
	}
	if u.Energy < e {
		return false
	}
	u.Energy -= e
	return true
}

// RestoreEnergy adds e energy, clamping at EnergyCap.
func (u *User) RestoreEnergy(e int) {
	if e < 0 {
		panic(fmt.Errorf("cannot restore %v energy", e))
	}
	u.Energy += e
	if u.Energy > EnergyCap {
		u.Energy = EnergyCap
	}
}

// LowEnergy is true if the user's energy is below LowEnergy.
func (u *User) LowEnergy() bool { return u.Energy < LowEnergy }

// FullEnergy is true if the user's energy is equal to EnergyCap.
func (u *User) FullEnergy() bool { return u.Energy == EnergyCap }
