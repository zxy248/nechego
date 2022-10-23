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
		return requireReply(a.injectReplyUser(a.handleFight))
	case input.CommandBalance:
		return a.handleBalance
	case input.CommandTransfer:
		return requireNonDebtor(requireReply(a.injectReplyUser(a.handleTransfer)))
	case input.CommandProfile:
		return a.handleProfile
	case input.CommandTopRich:
		return a.handleTopRich
	case input.CommandTopPoor:
		return a.handleTopPoor
	case input.CommandCapital:
		return a.handleCapital
	case input.CommandStrength:
		return a.handleStrength
	case input.CommandEnergy:
		return handleEnergy
	case input.CommandRating:
		return a.handleTopElo
	case input.CommandFishingRod:
		return a.handleFishingRod
	case input.CommandFishing:
		return a.handleFishing
	case input.CommandTopStrong:
		return a.handleTopStrong
	case input.CommandEatFish:
		return a.handleEatFood
	case input.CommandDeposit:
		return a.handleDeposit
	case input.CommandWithdraw:
		return requireNonDebtor(a.handleWithdraw)
	case input.CommandBank:
		return a.handleBank
	case input.CommandDebt:
		return requireNonDebtor(a.handleDebt)
	case input.CommandRepay:
		return requireDebtor(a.handleRepay)
	case input.CommandTopWeak:
		return a.handleTopWeak
	case input.CommandParliament:
		return a.handleParliament
	case input.CommandImpeachment:
		return a.handleImpeachment
	case input.CommandFishList:
		return a.handleFish
	case input.CommandFreezeFish:
		return a.handleFreeze
	case input.CommandFreezer:
		return a.handleFreezer
	case input.CommandUnfreezeFish:
		return a.handleUnfreeze
	case input.CommandSellFish:
		return a.handleSellFish
	case input.CommandMasyunya:
		return a.handleMasyunya
	case input.CommandPoppy:
		return a.handlePoppy
	case input.CommandSima:
		return a.handleSima
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
		return requireAdmin(requireReply(a.injectReplyUser(a.handleBan)))
	case input.CommandUnban:
		return requireAdmin(requireReply(a.injectReplyUser(a.handleUnban)))
	case input.CommandInfo:
		return a.handleInfo
	case input.CommandHelp:
		return a.handleHelp
	case input.CommandForbid:
		return requireAdmin(a.handleForbid)
	case input.CommandPermit:
		return requireAdmin(a.handlePermit)
	case input.CommandPet:
		return a.handlePet
	case input.CommandBuyPet:
		return requireNonDebtor(a.handleBuyPet)
	case input.CommandNamePet:
		return a.handleNamePet
	case input.CommandDropPet:
		return a.handleDropPet
	}
	return nil
}
