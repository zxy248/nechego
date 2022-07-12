package app

import (
	"errors"
	"fmt"
	"nechego/input"
	"nechego/model"

	"golang.org/x/exp/slices"
	tele "gopkg.in/telebot.v3"
)

const (
	parliamentMembers    = 4
	impeachmentThreshold = 3
	parliamentTemplate   = `
📰 На сегодняшний день в парламенте заседают:
%s`
	notInParliament    = "Вы не состоите в парламенте."
	alreadyImpeached   = "Вы уже проголосовали за импичмент."
	impeachmentSuccess = "Вынесен импичмент. Вы можете выбрать нового администратора."
	impeachmentPartial = "До вынесения импичмента осталось %d голосов."
	impeachedToday     = "Сегодня уже был вынесен импичмент."
)

func (a *App) handleParliament(c tele.Context) error {
	group := getGroup(c)
	parliament, err := a.model.Parliament(group, parliamentMembers)
	if err != nil {
		return internalError(c, err)
	}
	out := fmt.Sprintf(parliamentTemplate, a.itemizeUsers(parliament...))
	return c.Send(out, tele.ModeMarkdownV2)

}

func (a *App) handleImpeachment(c tele.Context) error {
	group := getGroup(c)
	user := getUser(c)

	count, err := a.model.Impeachment(group, user, impeachmentThreshold)
	if err != nil {
		if errors.Is(err, model.ErrNotInParliament) {
			return userError(c, notInParliament)
		}
		if errors.Is(err, model.ErrAlreadyImpeached) {
			return userError(c, alreadyImpeached)
		}
		if errors.Is(err, model.ErrImpeachedToday) {
			return userError(c, impeachedToday)
		}
		return internalError(c, err)
	}
	if count == impeachmentThreshold {
		return c.Send(impeachmentSuccess)
	}
	return c.Send(fmt.Sprintf(impeachmentPartial, impeachmentThreshold-count))
}

func (a *App) isCommandForbidden(g model.Group, c input.Command) (bool, error) {
	forbidden, err := a.model.ForbiddenCommands(g)
	if err != nil {
		return false, err
	}
	if slices.Contains(forbidden, c) {
		return true, nil
	}
	return false, nil
}
