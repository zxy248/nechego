package app

import (
	tele "gopkg.in/telebot.v3"
)

const dailyPair = Response("<b>✨ Пара дня</b> — %s 💘 %s")

// !пара дня
func (a *App) handlePair(c tele.Context) error {
	x, y, err := a.service.DailyPair(getGroup(c))
	if err != nil {
		return respondInternalError(c, err)
	}
	return respond(c, dailyPair.Fill(a.mustMentionUser(x), a.mustMentionUser(y)))
}

const dailyEblan = Response("<b>Еблан дня</b> — %s 😸")

// !еблан дня
func (a *App) handleEblan(c tele.Context) error {
	user, err := a.service.DailyEblan(getGroup(c))
	if err != nil {
		return respondInternalError(c, err)
	}
	return respond(c, dailyEblan.Fill(a.mustMentionUser(user)))
}

const dailyAdmin = Response("<b>Админ дня</b> — %s 👑")

// !админ дня
func (a *App) handleAdmin(c tele.Context) error {
	user, err := a.service.DailyAdmin(getGroup(c))
	if err != nil {
		return respondInternalError(c, err)
	}
	return respond(c, dailyAdmin.Fill(a.mustMentionUser(user)))
}
