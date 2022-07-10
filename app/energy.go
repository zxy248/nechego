package app

import (
	"nechego/model"
	"time"
)

const (
	restoreEnergyCooldown = time.Minute * 15
	energyDelta           = 1
	energyCap             = 4
)

func (a *App) restoreEnergyEvery(d time.Duration) {
	for range time.Tick(d) {
		a.model.RestoreEnergy(energyDelta, energyCap)
	}
}

func hasFullEnergy(u model.User) bool {
	return u.Energy == energyCap
}

func hasNoEnergy(u model.User) bool {
	return u.Energy == 0
}
