package app

import (
	"errors"
	"math/rand"
	"nechego/input"
	"nechego/service"
	"regexp"

	tele "gopkg.in/telebot.v3"
)

const randomPhotoChance = 0.01

func (a *App) handleRandomPhoto(c tele.Context) error {
	if rand.Float64() < randomPhotoChance {
		if rand.Float64() < 0.5 {
			if a, ok := loadAvatar(getUser(c).UID); ok {
				return c.Send(a)
			}
		}
		return sendLargeProfilePhoto(c)
	}
	return nil
}

func (a *App) handleJoin(c tele.Context) error {
	u := c.Message().UserJoined
	m, err := c.Bot().ChatMemberOf(c.Chat(), u)
	if err != nil {
		return err
	}
	if !isTeleAdmin(m) {
		if err := promoteAdmin(c, m); err != nil {
			return err
		}
	}
	return c.Send(a.stickers.Hello.Random())
}

func isTeleAdmin(m *tele.ChatMember) bool {
	return m.Role == tele.Administrator || m.Role == tele.Creator
}

func promoteAdmin(c tele.Context, m *tele.ChatMember) error {
	m.Rights.CanBeEdited = true
	m.Rights.CanManageChat = true
	return c.Bot().Promote(c.Chat(), m)
}

const (
	allCommandsPermitted    = Response("Все команды разрешены ✅")
	commandForbiddenSuccess = Response("Команда запрещена 🚫")
	commandPermittedSuccess = Response("Команда разрешена ✅")
	commandAlreadyForbidden = UserError("Команда уже запрещена.")
	commandAlreadyPermitted = UserError("Команда уже разрешена.")
	commandNotForbiddable   = UserError("Эту команду нельзя запретить.")
)

// !запретить
func (a *App) handleForbid(c tele.Context) error {
	return actOnCommand(c, func(command input.Command) error {
		if input.IsImmune(command) {
			return respondUserError(c, commandNotForbiddable)
		}
		if err := a.service.Forbid(getGroup(c), command); err != nil {
			if errors.Is(err, service.ErrAlreadyForbidden) {
				return respondUserError(c, commandAlreadyForbidden)
			}
			return respondInternalError(c, err)
		}
		return respond(c, commandForbiddenSuccess)
	})
}

var regexpAll = regexp.MustCompile("вс[её]")

// !разрешить
func (a *App) handlePermit(c tele.Context) error {
	if regexpAll.MatchString(getMessage(c).Argument()) {
		if err := a.service.PermitAll(getGroup(c)); err != nil {
			return respondInternalError(c, err)
		}
		return respond(c, allCommandsPermitted)
	}
	return actOnCommand(c, func(command input.Command) error {
		if err := a.service.Permit(getGroup(c), command); err != nil {
			if errors.Is(err, service.ErrAlreadyPermitted) {
				return respondUserError(c, commandAlreadyPermitted)
			}
			return respondInternalError(c, err)
		}
		return respond(c, commandPermittedSuccess)
	})
}

const (
	specifyCommand = UserError("Укажите команду.")
	unknownCommand = UserError("Неизвестная команда.")
)

func actOnCommand(c tele.Context, f func(input.Command) error) error {
	command, err := getMessage(c).CommandActionArgument()
	if err != nil {
		if errors.Is(err, input.ErrNoCommand) {
			return respondUserError(c, specifyCommand)
		}
		if errors.Is(err, input.ErrUnknownCommand) {
			return respondUserError(c, unknownCommand)
		}
		return respondInternalError(c, err)
	}
	return f(command)
}
