package app

import (
	"errors"
	"nechego/input"
	"nechego/model"
	"nechego/service"
	"time"

	tele "gopkg.in/telebot.v3"
)

func (a *App) pipeline(next tele.HandlerFunc) tele.HandlerFunc {
	line := []tele.MiddlewareFunc{
		ignoreForwarded,
		ignoreNotGroup,
		injectMessage,
		a.logMessage,
		a.injectGroup,
		requireGroupWhitelisted,
		requireStatusActive,
		a.injectUser,
		a.raiseLimit,
		requireUserUnbanned,
		a.requireCommandPermitted,
		a.incrementMessageCount,
	}
	for i := len(line) - 1; i >= 0; i-- {
		next = line[i](next)
	}
	return next
}

// BUG: don't work in groups
func ignoreForwarded(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if c.Message().IsForwarded() {
			return nil
		}
		return next(c)
	}
}

// isGroup returns true if the chat type is a group type, false otherwise.
func isGroup(t tele.ChatType) bool {
	return t == tele.ChatGroup || t == tele.ChatSuperGroup
}

func ignoreNotGroup(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if isGroup(c.Chat().Type) {
			return next(c)
		}
		return nil
	}
}

func injectMessage(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		return next(addMessage(c, input.ParseMessage(c.Text())))
	}
}

func (a *App) injectGroup(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		g, err := a.service.Group(model.Group{GID: c.Chat().ID})
		if err != nil {
			return err
		}
		return next(addGroup(c, g))
	}
}

func (a *App) injectUser(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		u, err := a.service.User(model.User{GID: c.Chat().ID, UID: c.Sender().ID})
		if err != nil {
			return err
		}
		return next(addUser(c, u))
	}
}

func (a *App) raiseLimit(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		u := getUser(c)
		u.DebtLimit = a.service.RaiseLimit(u, u.Summary())
		return next(addUser(c, u))
	}
}

const userNotFound = UserError("Пользователь не найден.")

func (a *App) injectReplyUser(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if !c.Message().IsReply() {
			return next(c)
		}
		u, err := a.service.FindUser(model.User{GID: c.Chat().ID, UID: c.Message().ReplyTo.Sender.ID})
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				return respondUserError(c, userNotFound)
			}
			return respondInternalError(c, err)
		}
		return next(addReplyUser(c, u))
	}
}

const accessRestricted = UserError("Недостаточно прав.")

func requireAdmin(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if getUser(c).Admin {
			return next(c)
		}
		return respondUserError(c, accessRestricted)
	}
}

const replyRequired = UserError("Перешлите сообщение пользователя.")

func requireReply(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if !c.Message().IsReply() || c.Message().ReplyTo.Sender.IsBot {
			return respondUserError(c, replyRequired)
		}
		return next(c)
	}
}

func requireGroupWhitelisted(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if getGroup(c).Whitelisted {
			return next(c)
		}
		return nil
	}
}

func requireUserUnbanned(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if input.IsImmune(getMessage(c).Command) {
			return next(c)
		}
		if getUser(c).Banned {
			return nil
		}
		return next(c)
	}
}

const commandForbidden = UserError("Эта команда запрещена администратором.")

func (a *App) requireCommandPermitted(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		command := getMessage(c).Command
		if input.IsImmune(command) {
			return next(c)
		}
		forbidden, err := a.service.IsCommandForbidden(getGroup(c), command)
		if err != nil {
			return err
		}
		if forbidden {
			return respondUserError(c, commandForbidden)
		}
		return next(c)
	}
}

func requireStatusActive(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		command := getMessage(c).Command
		if command == input.CommandTurnOn || command == input.CommandTurnOff {
			return next(c)
		}
		if getGroup(c).Status {
			return next(c)
		}
		return nil
	}
}

const (
	debtor    = UserError("Вы не можете использовать эту функцию, пока у вас есть непогашенные кредиты.")
	notDebtor = UserError("У вас нет непогашенных кредитов.")
)

func requireNonDebtor(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if getUser(c).Debtor() {
			return respondUserError(c, debtor)
		}
		return next(c)
	}
}

func requireDebtor(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if !getUser(c).Debtor() {
			return respondUserError(c, notDebtor)
		}
		return next(c)
	}
}

func technicalMaintenance(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		return respondUserError(c, UserError("Технические работы."))
	}
}

func (a *App) logMessage(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		t0 := time.Now()
		msg := getMessage(c)
		err := next(c)
		a.log.Sugar().Infow(msg.Raw,
			"cmd", msg.Command,
			"uid", c.Sender().ID,
			"gid", c.Chat().ID,
			"time", time.Since(t0))
		return err
	}
}

func (a *App) incrementMessageCount(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		a.service.IncrementMessages(getUser(c))
		return next(c)
	}
}
