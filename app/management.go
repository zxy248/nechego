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
	return respond(c, botTurnedOff.Fill(inactiveEmoji()), tele.RemoveKeyboard)
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
	return respond(c, info.Fill(joinSections(
		a.formatAdminList(admins),
		a.formatBlackList(bans),
		a.formatForbiddenCommandList(commands)),
	))
}

const (
	usersNotDeleted = Response("♻️ Некого удалить.")
	usersDeleted    = Response("♻️ Пользователи удалены:\n%s")
)

// !очистка
func (a *App) handleClean(c tele.Context) error {
	absent := []model.User{}
	if err := a.service.DeleteUsers(getGroup(c), func(u model.User) bool {
		memb, err := a.chatMember(u)
		if err != nil {
			return false
		}
		if chatMemberAbsent(memb) {
			absent = append(absent, u)
			return true
		}
		return false
	}); err != nil {
		return respondInternalError(c, err)
	}
	if len(absent) == 0 {
		return respond(c, usersNotDeleted)
	}
	return respond(c, usersDeleted.Fill(a.itemizeUsers(absent...)))
}

const adminListHeader = "👤 <i>Администрация</i>"

func (a *App) formatAdminList(u []model.User) string {
	return joinLines(adminListHeader, a.itemizeUsers(u...))
}

const blackListHeader = "🛑 <i>Черный список</i>"

func (a *App) formatBlackList(u []model.User) string {
	return joinLines(blackListHeader, a.itemizeUsers(u...))
}

const forbiddenCommandListHeader = "🔒 <i>Запрещенные команды</i>"

func (a *App) formatForbiddenCommandList(c []input.Command) string {
	return joinLines(forbiddenCommandListHeader, itemizeCommands(c...))
}
