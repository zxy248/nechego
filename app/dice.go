package app

import (
	"errors"
	"nechego/dice"
	"nechego/input"
	"nechego/service"
	"time"

	tele "gopkg.in/telebot.v3"
)

const (
	diceStart      = Response("🎲 %s играет на %s\nУ вас <code>%d секунд</code> на то, чтобы кинуть кости!")
	diceWin        = Response("💥 Вы выиграли %s")
	diceDraw       = Response("Ничья.")
	diceLose       = Response("Вы проиграли.")
	diceBonus      = Response("<i>🎰 %s получает бонус за риск: %s</i>")
	diceTimeout    = Response("<i>Время вышло: вы потеряли %s</i>")
	diceInProgress = UserError("Игра уже идет.")
	betTooLow      = UserError("Поставьте больше средств.")
	tired          = UserError("<i>Вы устали от азартных игр.</i>")
)

func (a *App) handleDice(c tele.Context) error {
	user := getUser(c)
	group := getGroup(c)
	bet, err := getMessage(c).MoneyArgument()
	if errors.Is(err, input.ErrAllIn) {
		bet = user.Balance
	} else if err != nil {
		return respondUserError(c, specifyAmount)
	}
	act := dice.Actions{
		Throw: func() (int, error) {
			message, err := tele.Cube.Send(c.Bot(), c.Chat(), &tele.SendOptions{})
			if err != nil {
				return 0, err
			}
			return message.Dice.Value, nil
		},
		Timeout: func() {
			respond(c, diceTimeout.Fill(formatMoney(bet)))
		},
	}
	if err := a.service.Dice(group, user, bet, act); err != nil {
		if errors.Is(err, service.ErrNotEnoughEnergy) {
			return respondUserError(c, notEnoughEnergy)
		}
		if errors.Is(err, service.ErrNotEnoughMoney) {
			return respondUserError(c, notEnoughMoney)
		}
		if errors.Is(err, service.ErrBetTooLow) {
			return respondUserError(c, betTooLow)
		}
		if errors.Is(err, dice.ErrGameInProgress) {
			return respondUserError(c, diceInProgress)
		}
		return respondInternalError(c, err)
	}
	return respond(c, diceStart.Fill(
		a.mustMentionUser(user),
		formatMoney(bet),
		a.service.Casino.Settings.RollTime/time.Second,
	))
}

func (a *App) handleRoll(c tele.Context) error {
	group := getGroup(c)
	user := getUser(c)
	result, err := a.service.Roll(group, user, c.Message().Dice.Value)
	if err != nil {
		if errors.Is(err, dice.ErrNoGame) {
			return nil
		}
		if errors.Is(err, dice.ErrWrongUser) {
			return nil
		}
		return respondInternalError(c, err)
	}
	switch result.Outcome {
	case dice.Win:
		return respond(c, diceWin.Fill(formatMoney(result.Bet*2)))
	case dice.Draw:
		return respond(c, diceDraw)
	case dice.Lose:
		return respond(c, diceLose)
	}
	return nil
}
