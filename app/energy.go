package app

import (
	"fmt"
	"nechego/model"
	"sync"
	"time"

	tele "gopkg.in/telebot.v3"
)

const (
	notEnoughEnergy       = "Недостаточно энергии."
	restoreEnergyCooldown = time.Minute * 20
	energyDelta           = 1
	energyCap             = 6
	energyLimit           = 1_000_000
)

type muTime struct {
	sync.RWMutex
	time time.Time
}

func (e *muTime) when() time.Time {
	e.RLock()
	defer e.RUnlock()
	return e.time
}

func (e *muTime) set(t time.Time) {
	e.Lock()
	defer e.Unlock()
	e.time = t
}

var energyCooldown *muTime

func (a *App) restoreEnergyEvery(d time.Duration) {
	a.model.RestoreEnergy(energyDelta, energyCap)
	energyCooldown = &muTime{time: time.Now().Add(d)}
	for t := range time.Tick(d) {
		energyCooldown.set(t.Add(d))
		a.model.RestoreEnergy(energyDelta, energyCap)
	}
}

func hasMuchEnergy(u model.User) bool {
	return u.Energy > energyCap
}

func hasFullEnergy(u model.User) bool {
	return u.Energy >= energyCap
}

func hasNoEnergy(u model.User) bool {
	return u.Energy == 0
}

// !стамина, !энергия
func (a *App) handleEnergy(c tele.Context) error {
	t := energyCooldown.when().Sub(time.Now())
	e := getUser(c).Energy
	return respondHTML(c, energyCooldownResponse(t, e))
}

const energyCooldownFormat = "⏰ До восстановления энергии: <code>%d минут %d секунд</code>."

func energyCooldownResponse(t time.Duration, energy int) string {
	mins := int(t.Minutes())
	secs := int(t.Seconds()) % 60
	return joinSections(
		fmt.Sprintf(energyCooldownFormat, mins, secs),
		energyRemaining(energy))
}
