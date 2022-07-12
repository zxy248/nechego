package app

import (
	"errors"
	"nechego/input"
	"nechego/model"
	"time"

	tele "gopkg.in/telebot.v3"
)

func (a *App) pipeline(next tele.HandlerFunc) tele.HandlerFunc {
	line := []tele.MiddlewareFunc{
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
		gid := c.Chat().ID
		g, err := a.model.GetGroup(model.Group{
			GID: gid,
		})
		if errors.Is(err, model.ErrGroupNotFound) {
			g = model.Group{
				GID:         gid,
				Whitelisted: false,
				Status:      true,
			}
			a.model.InsertGroup(g)
		} else if err != nil {
			return err
		}
		return next(addGroup(c, g))
	}
}

func (a *App) injectUser(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		gid := c.Chat().ID
		uid := c.Sender().ID
		u, err := a.model.GetUser(model.User{
			GID: gid,
			UID: uid,
		})
		if errors.Is(err, model.ErrUserNotFound) {
			u = model.User{
				GID:     gid,
				UID:     uid,
				Energy:  energyCap,
				Account: initialBalance,
			}
			a.model.InsertUser(u)
		} else if err != nil {
			return err
		}
		return next(addUser(c, u))
	}
}

func (a *App) raiseLimit(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		u := getUser(c)
		sum := u.Summary()
		if u.DebtLimit < sum {
			u.DebtLimit = sum
			a.model.RaiseLimit(u, sum)
		}
		return next(addUser(c, u))
	}
}

const userNotFound = "Пользователь не найден 🔎"

func (a *App) injectReplyUser(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		gid := c.Chat().ID
		uid := c.Message().ReplyTo.Sender.ID
		u, err := a.model.GetUser(model.User{
			GID: gid,
			UID: uid,
		})
		if err != nil {
			return c.Send(makeError(userNotFound))
		}
		return next(addReplyUser(c, u))
	}
}

const accessRestricted = "Доступ ограничен 🔒"

func requireAdmin(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if getUser(c).Admin {
			return next(c)
		}
		return c.Send(makeError(accessRestricted))
	}
}

const (
	replyRequired     = "Перешлите сообщение ↩️"
	userReplyRequired = "Перешлите сообщение пользователя ↩️"
)

func requireReply(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if c.Message().IsReply() {
			if c.Message().ReplyTo.Sender.IsBot {
				return c.Send(makeError(userReplyRequired))
			}
			return next(c)
		}
		return c.Send(makeError(replyRequired))
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
		command := getMessage(c).Command
		if input.IsImmune(command) {
			return next(c)
		}

		user := getUser(c)
		if user.Banned {
			return nil
		}
		return next(c)
	}
}

func (a *App) requireCommandPermitted(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		command := getMessage(c).Command
		if input.IsImmune(command) {
			return next(c)
		}

		forbidden, err := a.isCommandForbidden(getGroup(c), command)
		if err != nil {
			return err
		}
		if forbidden {
			return nil
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

func (a *App) logMessage(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		t0 := time.Now()
		msg := getMessage(c)
		err := next(c)
		a.SugarLog().Infow(msg.Raw,
			"cmd", msg.Command,
			"uid", c.Sender().ID,
			"gid", c.Chat().ID,
			"time", time.Since(t0))
		return err
	}
}

func (a *App) incrementMessageCount(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		a.model.IncrementMessages(getUser(c))
		return next(c)
	}
}
