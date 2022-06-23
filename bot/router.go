package bot

import (
	"nechego/input"

	tele "gopkg.in/telebot.v3"
)

// route processes the input message and handles it to the
// corresponding function. Ignores the message if the chat type is not
// a group. Caches the group ID and the user ID.
func (b *Bot) check(next tele.HandlerFunc) tele.HandlerFunc {
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
		b.cacheGroupMember(gid, uid)

		ok, err = b.status.Active(gid)
		if err != nil {
			return err
		}
		if !ok {
			if message.Command == input.CommandTurnOn {
				return b.handleTurnOn(c)
			}
			return nil
		}

		c = addMessage(c, message)
		c = addCommandForbidden(c, commandForbidden)
		c = addAdminAuthorized(c, adminAuthorized)
		return next(c)
	}
}

func (b *Bot) getHandler(c input.Command) tele.HandlerFunc {
	switch c {
	case input.CommandProbability:
		return b.handleProbability
	case input.CommandWho:
		return b.handleWho
	case input.CommandCat:
		return b.handleCat
	case input.CommandTitle:
		return b.handleTitle
	case input.CommandAnime:
		return b.handleAnime
	case input.CommandFurry:
		return b.handleFurry
	case input.CommandFlag:
		return b.handleFlag
	case input.CommandPerson:
		return b.handlePerson
	case input.CommandHorse:
		return b.handleHorse
	case input.CommandArt:
		return b.handleArt
	case input.CommandCar:
		return b.handleCar
	case input.CommandPair:
		return b.handlePair
	case input.CommandEblan:
		return b.handleEblan
	case input.CommandAdmin:
		return b.handleAdmin
	case input.CommandMasyunya:
		return b.handleMasyunya
	case input.CommandPoppy:
		return b.handlePoppy
	case input.CommandHello:
		return b.handleHello
	case input.CommandMouse:
		return b.handleMouse
	case input.CommandWeather:
		return b.handleWeather
	case input.CommandTikTok:
		return b.handleTikTok
	case input.CommandList:
		return b.handleList
	case input.CommandTop:
		return b.handleTop
	case input.CommandBasili:
		return b.handleBasili
	case input.CommandCasper:
		return b.handleCasper
	case input.CommandZeus:
		return b.handleZeus
	case input.CommandPic:
		return b.handlePic
	case input.CommandDice:
		return b.handleDice
	case input.CommandGame:
		return b.handleGame
	case input.CommandKeyboardOpen:
		return b.handleKeyboardOpen
	case input.CommandKeyboardClose:
		return b.handleKeyboardClose
	case input.CommandTurnOff:
		return b.handleTurnOff
	case input.CommandBan:
		return requireReply(requireAdminAccess(b.handleBan))
	case input.CommandUnban:
		return requireReply(requireAdminAccess(b.handleUnban))
	case input.CommandInfo:
		return b.handleInfo
	case input.CommandHelp:
		return b.handleHelp
	case input.CommandForbid:
		return requireAdminAccess(b.handleForbid)
	case input.CommandPermit:
		return requireAdminAccess(b.handlePermit)
	}
	return handleNothing
}

// routeMessage routes the input message to the appropriate handler.
func (b *Bot) route(c tele.Context) error {
	if err := b.handleRandomPhoto(c); err != nil {
		return err
	}
	return b.getHandler(getMessage(c).Command)(c)
}

// isGroup returns true if the chat type is one of the group types; false otherwise.
func isGroup(t tele.ChatType) bool {
	return t == tele.ChatGroup || t == tele.ChatSuperGroup
}

func (b *Bot) cacheGroupMember(gid, uid int64) error {
	ok, err := b.users.Exists(gid, uid)
	if err != nil {
		return err
	}
	if !ok {
		if err := b.users.Insert(gid, uid); err != nil {
			return err
		}
	}
	return nil
}

const messageKey = "nechegoParsedMessage"

func addMessage(c tele.Context, m *input.Message) tele.Context {
	c.Set(messageKey, m)
	return c
}

func getMessage(c tele.Context) *input.Message {
	return c.Get(messageKey).(*input.Message)
}

const commandForbiddenKey = "commandForbidden"

func addCommandForbidden(c tele.Context, v bool) tele.Context {
	c.Set(commandForbiddenKey, v)
	return c
}

func isCommandForbidden(c tele.Context) bool {
	return c.Get(commandForbiddenKey).(bool)
}

const adminAuthorizedKey = "adminAuthorized"

func addAdminAuthorized(c tele.Context, v bool) tele.Context {
	c.Set(adminAuthorizedKey, v)
	return c
}

func isAdminAuthorized(c tele.Context) bool {
	return c.Get(adminAuthorizedKey).(bool)
}

const accessRestricted = "Доступ ограничен 🔒"

func requireAdminAccess(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if isAdminAuthorized(c) {
			return next(c)
		}
		return c.Send(accessRestricted)
	}
}

const replyRequired = "Перешлите сообщение ↩️"

func requireReply(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		if c.Message().IsReply() {
			return next(c)
		}
		return c.Send(replyRequired)
	}
}
