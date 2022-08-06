package app

import (
	"errors"
	"nechego/service"

	tele "gopkg.in/telebot.v3"
)

const (
	userBlocked          = Response("Пользователь заблокирован 🚫")
	userUnblocked        = Response("Пользователь разблокирован ✅")
	userAlreadyBlocked   = UserError("Пользователь уже заблокирован.")
	userAlreadyUnblocked = UserError("Пользователь не заблокирован.")
)

// !бан
func (a *App) handleBan(c tele.Context) error {
	user := getReplyUser(c)
	if err := a.service.Ban(user); err != nil {
		if errors.Is(err, service.ErrAlreadyBanned) {
			return respondUserError(c, userAlreadyBlocked)
		}
		return respondInternalError(c, err)
	}
	return respond(c, userBlocked)
}

// !разбан
func (a *App) handleUnban(c tele.Context) error {
	user := getReplyUser(c)
	if err := a.service.Unban(user); err != nil {
		if errors.Is(err, service.ErrAlreadyUnbanned) {
			return respondUserError(c, userAlreadyUnblocked)
		}
		return respondInternalError(c, err)
	}
	return respond(c, userUnblocked)
}
