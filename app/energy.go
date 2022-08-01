package app

import (
	"nechego/service"
	"time"

	tele "gopkg.in/telebot.v3"
)

const (
	notEnoughEnergy      = UserError("Недостаточно энергии.")
	energyCooldownFormat = Response("⏰ До восстановления энергии: <code>%d минут %d секунд</code>.")
)

// !стамина, !энергия
func (a *App) energyHandler() tele.HandlerFunc {
	a.service.RestoreEnergy()
	next := service.PeriodicallyRun(a.service.RestoreEnergy, a.pref.EnergyPeriod)
	return func(c tele.Context) error {
		return respond(c, energyCooldownResponse(time.Until(next()), getUser(c).Energy))
	}
}

var handleEnergy tele.HandlerFunc

func energyCooldownResponse(t time.Duration, energy int) Response {
	mins := int(t.Minutes())
	secs := int(t.Seconds()) % 60
	return Response(joinSections(
		string(energyCooldownFormat.Fill(mins, secs)),
		string(energyRemaining(energy)),
	))
}
