package app

import (
	"nechego/input"

	tele "gopkg.in/telebot.v3"
)

// route routes an input message to a corresponding handler.
func (a *App) route(c tele.Context) error {
	if err := a.handleRandomPhoto(c); err != nil {
		return err
	}
	f := a.commandHandler(getMessage(c).Command)
	if f == nil {
		return nil
	}
	return f(c)
}

// commandHandler returns a corresponding handler for a command.
func (a *App) commandHandler(c input.Command) tele.HandlerFunc {
	switch c {
	case input.CommandProbability:
		return a.handleProbability
	case input.CommandWho:
		return a.handleWho
	case input.CommandCat:
		return a.handleCat
	case input.CommandTitle:
		return a.handleTitle
	case input.CommandAnime:
		return a.handleAnime
	case input.CommandFurry:
		return a.handleFurry
	case input.CommandFlag:
		return a.handleFlag
	case input.CommandPerson:
		return a.handlePerson
	case input.CommandHorse:
		return a.handleHorse
	case input.CommandArt:
		return a.handleArt
	case input.CommandCar:
		return a.handleCar
	case input.CommandPair:
		return a.handlePair
	case input.CommandEblan:
		return a.handleEblan
	case input.CommandAdmin:
		return a.handleAdmin
	case input.CommandFight:
		return requireReply(a.handleFight)
	case input.CommandBalance:
		return a.handleBalance
	case input.CommandTransfer:
		return requireReply(a.handleTransfer)
	case input.CommandProfile:
		return a.handleProfile
	case input.CommandTopRich:
		return a.handleTopRich
	case input.CommandTopPoor:
		return a.handleTopPoor
	case input.CommandMasyunya:
		return a.masyunyaHandler()
	case input.CommandPoppy:
		return a.poppyHandler()
	case input.CommandHello:
		return a.handleHello
	case input.CommandMouse:
		return a.handleMouse
	case input.CommandWeather:
		return a.handleWeather
	case input.CommandTikTok:
		return a.handleTikTok
	case input.CommandList:
		return a.handleList
	case input.CommandTop:
		return a.handleTop
	case input.CommandBasili:
		return a.handleBasili
	case input.CommandCasper:
		return a.handleCasper
	case input.CommandZeus:
		return a.handleZeus
	case input.CommandPic:
		return a.handlePic
	case input.CommandDice:
		return a.handleDice
	case input.CommandGame:
		return a.handleGame
	case input.CommandKeyboardOpen:
		return a.handleKeyboardOpen
	case input.CommandKeyboardClose:
		return a.handleKeyboardClose
	case input.CommandTurnOn:
		return a.handleTurnOn
	case input.CommandTurnOff:
		return a.handleTurnOff
	case input.CommandBan:
		return requireAdminAccess(requireReply(a.handleBan))
	case input.CommandUnban:
		return requireAdminAccess(requireReply(a.handleUnban))
	case input.CommandInfo:
		return a.handleInfo
	case input.CommandHelp:
		return a.handleHelp
	case input.CommandForbid:
		return requireAdminAccess(a.handleForbid)
	case input.CommandPermit:
		return requireAdminAccess(a.handlePermit)
	}
	return nil
}
