package game

import "fmt"

const EnergyCap = 100

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

func (u *User) RestoreEnergy(e int) {
	if e < 0 {
		panic(fmt.Errorf("cannot restore %v energy", e))
	}
	u.Energy += e
	if u.Energy > EnergyCap {
		u.Energy = EnergyCap
	}
}
