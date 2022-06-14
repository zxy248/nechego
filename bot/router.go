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
		uid := c.Sender().ID

		ok, err := b.whitelist.Allow(gid)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		text := c.Text()
		message := input.Parse(text)

		banned, err := b.bans.Banned(uid)
		if err != nil {
			return err
		}
		if banned && !(uid == b.config.Owner && message.Command == input.CommandUnban) {
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

		return next(addMessage(c, message))
	}
}

// routeMessage routes the input message to the appropriate handler.
func (b *Bot) route(c tele.Context) error {
	if err := b.handleRandomPhoto(c); err != nil {
		return err
	}
	switch getMessage(c).Command {
	case input.CommandProbability:
		return b.handleProbability(c)
	case input.CommandWho:
		return b.handleWho(c)
	case input.CommandCat:
		return b.handleCat(c)
	case input.CommandTitle:
		return b.handleTitle(c)
	case input.CommandAnime:
		return b.handleAnime(c)
	case input.CommandFurry:
		return b.handleFurry(c)
	case input.CommandFlag:
		return b.handleFlag(c)
	case input.CommandPerson:
		return b.handlePerson(c)
	case input.CommandHorse:
		return b.handleHorse(c)
	case input.CommandArt:
		return b.handleArt(c)
	case input.CommandCar:
		return b.handleCar(c)
	case input.CommandPair:
		return b.handlePair(c)
	case input.CommandEblan:
		return b.handleEblan(c)
	case input.CommandMasyunya:
		return b.handleMasyunya(c)
	case input.CommandHello:
		return b.handleHello(c)
	case input.CommandMouse:
		return b.handleMouse(c)
	case input.CommandWeather:
		return b.handleWeather(c)
	case input.CommandTikTok:
		return b.handleTikTok(c)
	case input.CommandList:
		return b.handleList(c)
	case input.CommandTop:
		return b.handleTop(c)
	case input.CommandBasili:
		return b.handleBasili(c)
	case input.CommandCasper:
		return b.handleCasper(c)
	case input.CommandZeus:
		return b.handleZeus(c)
	case input.CommandPic:
		return b.handlePic(c)
	case input.CommandDice:
		return b.handleDice(c)
	case input.CommandGame:
		return b.handleGame(c)
	case input.CommandKeyboardOpen:
		return b.handleKeyboardOpen(c)
	case input.CommandKeyboardClose:
		return b.handleKeyboardClose(c)
	case input.CommandTurnOff:
		return b.handleTurnOff(c)
	case input.CommandBan:
		return b.handleBan(c)
	case input.CommandUnban:
		return b.handleUnban(c)
	case input.CommandInfo:
		return b.handleInfo(c)
	case input.CommandHelp:
		return b.handleHelp(c)
	}
	return nil
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

const messageKey = `nechegoParsedMessage`

func addMessage(c tele.Context, m *input.Message) tele.Context {
	c.Set(messageKey, m)
	return c
}

func getMessage(c tele.Context) *input.Message {
	return c.Get(messageKey).(*input.Message)
}
