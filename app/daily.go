package app

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

const dailyPair = "Пара дня ✨\n%s 💘 %s"

// !пара дня
func (a *App) handlePair(c tele.Context) error {
	u1, u2, err := a.model.GetDailyPair(getGroup(c))
	if err != nil {
		return internalError(c, err)
	}
	return c.Send(fmt.Sprintf(dailyPair, a.mustMentionUser(u1), a.mustMentionUser(u2)),
		tele.ModeMarkdownV2)
}

const dailyEblan = "Еблан дня — %s 😸"

// !еблан дня
func (a *App) handleEblan(c tele.Context) error {
	u, err := a.model.GetDailyEblan(getGroup(c))
	if err != nil {
		return internalError(c, err)
	}
	return c.Send(fmt.Sprintf(dailyEblan, a.mustMentionUser(u)),
		tele.ModeMarkdownV2)
}

const dailyAdmin = "Админ дня — %s 👑"

// !админ дня
func (a *App) handleAdmin(c tele.Context) error {
	u, err := a.model.GetDailyAdmin(getGroup(c))
	if err != nil {
		return internalError(c, err)
	}
	return c.Send(fmt.Sprintf(dailyAdmin, a.mustMentionUser(u)),
		tele.ModeMarkdownV2)
}
