package app

import (
	"nechego/input"

	tele "gopkg.in/telebot.v3"
)

// route routes an input message to a corresponding handler.
func (b *App) route(c tele.Context) error {
	if err := b.handleRandomPhoto(c); err != nil {
		return err
	}
	f := b.commandHandler(getMessage(c).Command)
	if f == nil {
		return nil
	}
	return f(c)
}

// commandHandler returns a corresponding handler for a command.
func (b *App) commandHandler(c input.Command) tele.HandlerFunc {
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
		return b.masyunyaHandler()
	case input.CommandPoppy:
		return b.poppyHandler()
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
	case input.CommandTurnOn:
		return b.handleTurnOn
	case input.CommandTurnOff:
		return b.handleTurnOff
	case input.CommandBan:
		return requireAdminAccess(requireReply(b.handleBan))
	case input.CommandUnban:
		return requireAdminAccess(requireReply(b.handleUnban))
	case input.CommandInfo:
		return b.handleInfo
	case input.CommandHelp:
		return b.handleHelp
	case input.CommandForbid:
		return requireAdminAccess(b.handleForbid)
	case input.CommandPermit:
		return requireAdminAccess(b.handlePermit)
	}
	return nil
}
