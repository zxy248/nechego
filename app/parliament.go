package app

import (
	"errors"
	"nechego/service"

	tele "gopkg.in/telebot.v3"
)

const (
	parliamentMembers   = 4
	parliamentMajority  = 3
	parliamentarians    = Response("📰 В парламенте заседают:\n%s")
	impeachmentSuccess  = Response("Вынесен импичмент. Вы можете выбрать нового администратора.")
	impeachmentPartial  = Response("До вынесения импичмента осталось %d голосов.")
	notParliamentMember = UserError("Вы не состоите в парламенте.")
	alreadyVoted        = UserError("Вы уже проголосовали за импичмент.")
	alreadyImpeached    = UserError("Сегодня уже был вынесен импичмент.")
)

// !парламент
func (a *App) handleParliament(c tele.Context) error {
	parliament, err := a.service.Parliament(getGroup(c))
	if err != nil {
		return respondInternalError(c, err)
	}
	return respond(c, parliamentarians.Fill(a.itemizeUsers(parliament...)))
}

// !импичмент
func (a *App) handleImpeachment(c tele.Context) error {
	votesLeft, err := a.service.Impeachment(getGroup(c), getUser(c))
	if err != nil {
		if errors.Is(err, service.ErrNotParliamentMember) {
			return respondUserError(c, notParliamentMember)
		}
		if errors.Is(err, service.ErrAlreadyVoted) {
			return respondUserError(c, alreadyVoted)
		}
		if errors.Is(err, service.ErrAlreadyImpeached) {
			return respondUserError(c, alreadyImpeached)
		}
		return respondInternalError(c, err)
	}
	if votesLeft == 0 {
		return respond(c, impeachmentSuccess)
	}
	return respond(c, impeachmentPartial.Fill(votesLeft))
}
