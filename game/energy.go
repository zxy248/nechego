package game

const EnergyCap = 100

func (u *User) SpendEnergy(e int) bool {
	if u.Energy < e {
		return false
	}
	u.Energy -= e
	return true
}

func (u *User) RestoreEnergy(e int) {
	u.Energy += e
	if u.Energy > EnergyCap {
		u.Energy = EnergyCap
	}
}
