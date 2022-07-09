package app

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

const pairOfTheDayFormat = "Пара дня ✨\n%s 💘 %s"

// !пара дня
func (a *App) handlePair(c tele.Context) error {
	u1, u2, err := a.model.GetDailyPair(getGroup(c))
	if err != nil {
		return err
	}
	return c.Send(fmt.Sprintf(pairOfTheDayFormat,
		a.mustMentionUser(u1), a.mustMentionUser(u2)),
		tele.ModeMarkdownV2)
}

const eblanOfTheDayFormat = "Еблан дня: %s 😸"

// !еблан дня
func (a *App) handleEblan(c tele.Context) error {
	u, err := a.model.GetDailyEblan(getGroup(c))
	if err != nil {
		return err
	}
	return c.Send(fmt.Sprintf(eblanOfTheDayFormat, a.mustMentionUser(u)),
		tele.ModeMarkdownV2)
}

const adminOfTheDayFormat = "Админ дня: %s 👑"

func (a *App) handleAdmin(c tele.Context) error {
	u, err := a.model.GetDailyAdmin(getGroup(c))
	if err != nil {
		return err
	}
	return c.Send(fmt.Sprintf(adminOfTheDayFormat, a.mustMentionUser(u)),
		tele.ModeMarkdownV2)
}
