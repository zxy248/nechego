package app

import (
	"nechego/service"
	"time"

	tele "gopkg.in/telebot.v3"
)

const notEnoughEnergy = UserError("Недостаточно энергии.")

// !стамина, !энергия
func (a *App) energyHandler() tele.HandlerFunc {
	a.service.RestoreEnergy()
	next := service.PeriodicallyRun(a.service.RestoreEnergy, a.pref.EnergyPeriod)
	return func(c tele.Context) error {
		return respond(c, energyCooldownResponse(time.Until(next()), getUser(c).Energy))
	}
}

var handleEnergy tele.HandlerFunc

func energyCooldownResponse(d time.Duration, energy int) Response {
	return Response(joinSections(formatEnergyCooldown(d), formatEnergyRemaining(energy)))
}
