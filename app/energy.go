package app

import (
	"fmt"
	"nechego/model"
	"sync"
	"time"

	tele "gopkg.in/telebot.v3"
)

const (
	restoreEnergyCooldown = time.Minute * 12
	energyDelta           = 1
	energyCap             = 5
	energyTrueCap         = 1000
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

const (
	energyRestoreAfter = "Осталось энергии: %s\n\n" +
		"_Энергия будет восстановлена через %d минут %d секунд\\._"
	notEnoughEnergy = "Недостаточно энергии."
)

// !стамина, !энергия
func (a *App) handleEnergy(c tele.Context) error {
	t := energyCooldown.when()
	now := time.Now()
	mins := int(t.Sub(now).Minutes())
	secs := int(t.Sub(now).Seconds()) % 60
	out := fmt.Sprintf(energyRestoreAfter, formatEnergy(getUser(c).Energy), mins, secs)
	return c.Send(out, tele.ModeMarkdownV2)
}
