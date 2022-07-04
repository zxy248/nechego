package app

import (
	"time"
)

const (
	restoreEnergyCooldown = time.Minute * 10
	energyCap = 3
)

func (a *App) restoreEnergyEvery(d time.Duration) {
	for range time.Tick(d) {
		users, err := a.model.Users.All()
		if err != nil {
			a.sugar().Errorw("can't get users", "err", err)
		}
		for _, u := range users {
			if u.Energy < energyCap {
				if err := a.model.Energy.Update(u.GID, u.UID, 1); err != nil {
					a.sugar().Errorw("can't restore the energy",
						"err", err,
						"uid", u.UID,
						"gid", u.GID,
						"energy", u.Energy)
				}
			}
		}
	}
}
