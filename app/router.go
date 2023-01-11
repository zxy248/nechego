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
	case input.CommandBalance:
		return a.handleBalance
	case input.CommandTransfer:
		return requireNonDebtor(requireReply(a.injectReplyUser(a.handleTransfer)))
	case input.CommandTopRich:
		return a.handleTopRich
	case input.CommandTopPoor:
		return a.handleTopPoor
	case input.CommandCapital:
		return a.handleCapital
	case input.CommandStrength:
		return a.handleStrength
	case input.CommandRating:
		return a.handleTopElo
	case input.CommandTopStrong:
		return a.handleTopStrong
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
	case input.CommandDice:
		return a.handleDice
	case input.CommandKick:
		return a.injectReplyUser(a.handleKick)
	case input.CommandBan:
		return requireAdmin(requireReply(a.injectReplyUser(a.handleBan)))
	case input.CommandUnban:
		return requireAdmin(requireReply(a.injectReplyUser(a.handleUnban)))
	}
	return nil
}
