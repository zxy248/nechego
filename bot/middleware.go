package bot

import (
	"nechego/input"

	tele "gopkg.in/telebot.v3"
)

// preprocess parses an input message, ignores it on certain conditions, caches
// a group member, saves necessary data in the context.
func (b *Bot) preprocess(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if !isGroup(c.Chat().Type) {
			return nil
		}
		gid := c.Chat().ID
		ok, err := b.whitelist.Allow(gid)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		uid := c.Sender().ID
		text := c.Text()
		message := input.ParseMessage(text)

		userBanned, err := b.bans.Banned(uid)
		if err != nil {
			return err
		}
		commandForbidden, err := b.forbid.Forbidden(gid, message.Command)
		if err != nil {
			return err
		}
		adminAuthorized, err := b.admins.Authorize(gid, uid)
		if err != nil {
			return err
		}
		if (userBanned || commandForbidden) &&
			!(adminAuthorized && input.IsManagementCommand(message.Command)) {
			return nil
		}
		active, err := b.status.Active(gid)
		if err != nil {
			return err
		}
		if !active && message.Command != input.CommandTurnOn {
			return nil
		}

		b.cacheGroupMember(gid, uid)
		c = addMessage(c, message)
		c = addCommandForbidden(c, commandForbidden)
		c = addAdminAuthorized(c, adminAuthorized)
		return next(c)
	}
}

// isGroup returns true if the chat type is a group type, false otherwise.
func isGroup(t tele.ChatType) bool {
	return t == tele.ChatGroup || t == tele.ChatSuperGroup
}

// cacheGroupMember adds a user to the users table if it is not there already.
func (b *Bot) cacheGroupMember(gid, uid int64) error {
	exists, err := b.users.Exists(gid, uid)
	if err != nil {
		return err
	}
	if !exists {
		if err := b.users.Insert(gid, uid); err != nil {
			return err
		}
	}
	return nil
}

const messageKey = "nechegoParsedMessage"

// addMessage adds a message to the context.
func addMessage(c tele.Context, m *input.Message) tele.Context {
	c.Set(messageKey, m)
	return c
}

// getMessage gets a message from the context.
func getMessage(c tele.Context) *input.Message {
	return c.Get(messageKey).(*input.Message)
}

const commandForbiddenKey = "commandForbidden"

// addCommandForbidden adds a command forbidden flag to the context.
func addCommandForbidden(c tele.Context, v bool) tele.Context {
	c.Set(commandForbiddenKey, v)
	return c
}

// isCommandForbidden gets a command forbidden flag from the context.
func isCommandForbidden(c tele.Context) bool {
	return c.Get(commandForbiddenKey).(bool)
}

const adminAuthorizedKey = "adminAuthorized"

// addAdminAuthorized adds an admin authorized flag to the context.
func addAdminAuthorized(c tele.Context, v bool) tele.Context {
	c.Set(adminAuthorizedKey, v)
	return c
}

// isAdminAuthorized gets an admin authorized flag from the context.
func isAdminAuthorized(c tele.Context) bool {
	return c.Get(adminAuthorizedKey).(bool)
}

const accessRestricted = "Доступ ограничен 🔒"

// requireAdminAccess returns a closure that requires an admin status to proceed.
func requireAdminAccess(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if isAdminAuthorized(c) {
			return next(c)
		}
		return c.Send(accessRestricted)
	}
}

const replyRequired = "Перешлите сообщение ↩️"

// requireReply returns a closure that requires the message to be a reply to proceed.
func requireReply(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if c.Message().IsReply() {
			return next(c)
		}
		return c.Send(replyRequired)
	}
}
