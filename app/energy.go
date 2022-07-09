package app

import (
	"time"
)

const (
	restoreEnergyCooldown = time.Minute * 10
	energyDelta           = 1
	energyCap             = 3
)

func (a *App) restoreEnergyEvery(d time.Duration) {
	for range time.Tick(d) {
		a.model.RestoreEnergy(energyDelta, energyCap)
	}
}
