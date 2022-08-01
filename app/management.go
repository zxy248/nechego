package app

import (
	"errors"
	"nechego/input"
	"nechego/model"
	"nechego/service"

	tele "gopkg.in/telebot.v3"
)

// !открыть
func (a *App) handleKeyboardOpen(c tele.Context) error {
	return openKeyboard(c)
}

// !закрыть
func (a *App) handleKeyboardClose(c tele.Context) error {
	return closeKeyboard(c)
}

const (
	botTurnedOn         = Response("Бот включен %s")
	botAlreadyTurnedOn  = UserError("Бот уже включен.")
	botTurnedOff        = Response("Бот выключен %s")
	botAlreadyTurnedOff = UserError("Бот уже выключен.")
)

// !включить
func (a *App) handleTurnOn(c tele.Context) error {
	if err := a.service.TurnOn(getGroup(c)); err != nil {
		if errors.Is(err, service.ErrAlreadyTurnedOn) {
			return respondUserError(c, botAlreadyTurnedOn)
		}
		return respondInternalError(c, err)
	}
	return respond(c, botTurnedOn.Fill(activeEmoji()))
}

// !выключить
func (a *App) handleTurnOff(c tele.Context) error {
	if err := a.service.TurnOff(getGroup(c)); err != nil {
		if errors.Is(err, service.ErrAlreadyTurnedOff) {
			return respondUserError(c, botAlreadyTurnedOff)
		}
		return respondInternalError(c, err)
	}
	if err := closeKeyboard(c); err != nil {
		return respondInternalError(c, err)
	}
	return respond(c, botTurnedOff.Fill(inactiveEmoji()))
}

const info = Response("ℹ️ <b>Информация</b> 📌\n\n%s")

// !инфо
func (a *App) handleInfo(c tele.Context) error {
	group := getGroup(c)
	admins, err := a.service.Admins(group)
	if err != nil {
		return respondInternalError(c, err)
	}
	bans, err := a.service.Bans(group)
	if err != nil {
		return respondInternalError(c, err)
	}
	commands, err := a.service.ForbiddenCommands(group)
	if err != nil {
		return respondInternalError(c, err)
	}
	return respond(c, info.Fill(HTML(joinSections(
		string(a.adminSection(admins)),
		string(a.bansSection(bans)),
		string(a.forbiddenCommandsSection(commands))),
	)))
}

const adminListHeader = "👤 <i>Администрация</i>"

func (a *App) adminSection(u []model.User) HTML {
	return HTML(joinLines(adminListHeader, string(a.itemizeUsers(u...))))
}

const bansHeader = "🛑 <i>Черный список</i>"

func (a *App) bansSection(u []model.User) HTML {
	return HTML(joinLines(bansHeader, string(a.itemizeUsers(u...))))
}

const forbiddenCommandsHeader = "🔒 <i>Запрещенные команды</i>"

func (a *App) forbiddenCommandsSection(c []input.Command) HTML {
	return HTML(joinLines(forbiddenCommandsHeader, string(itemizeCommands(c...))))
}
