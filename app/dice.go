package app

import (
	"errors"
	"fmt"
	"nechego/input"
	"nechego/model"
	"sync"
	"time"

	tele "gopkg.in/telebot.v3"
)

type diceGame struct {
	id     time.Time
	uid    int64
	amount uint
	roll   int
}

var diceGames = &sync.Map{}

const (
	handleDiceTemplate = "🎲 %s играет на `%s 💰`\nУ вас `%v секунд` на то, чтобы кинуть кости\\!"
	secondsForRoll     = 25
	minimalBet         = 1
)

var handleDiceMutex = &sync.Mutex{}

// handleDice rolls a dice.
func (a *App) handleDice(c tele.Context) error {
	handleDiceMutex.Lock()
	defer handleDiceMutex.Unlock()

	gid := c.Chat().ID
	uid := c.Sender().ID
	if currentDiceGame(gid) != nil {
		return c.Send(makeError("Игра уже идет"))
	}

	arg, err := getMessage(c).Dynamic()
	if err != nil {
		if errors.Is(err, input.ErrSpecifyAmount) {
			return c.Send(makeError(input.ErrSpecifyAmount.Error()))
		}
		return err
	}
	amount := arg.(uint)
	if amount < minimalBet {
		return c.Send(makeError("Поставьте больше средств"))
	}
	if err := a.model.Economy.Update(gid, uid, int(-amount)); err != nil {
		if errors.Is(err, model.ErrNotEnoughMoney) {
			return c.Send(makeError("Недостаточно средств"))
		}
		return err
	}

	dice := &tele.Dice{Type: tele.Cube.Type}
	msg, err := dice.Send(c.Bot(), c.Chat(), &tele.SendOptions{})
	if err != nil {
		return err
	}
	roll := msg.Dice.Value

	a.newDiceGame(gid, uid, amount, roll)

	member, err := a.chatMember(gid, uid)
	if err != nil {
		return err
	}
	out := fmt.Sprintf(handleDiceTemplate,
		mentionName(uid, markdownEscaper.Replace(chatMemberName(member))),
		formatAmount(int(amount)),
		secondsForRoll)
	return c.Send(out, tele.ModeMarkdownV2)
}

func (a *App) handleRoll(c tele.Context) error {
	gid := c.Chat().ID
	uid := c.Sender().ID
	game := currentDiceGame(gid)
	if game == nil {
		return nil
	}
	if game.uid != uid {
		return nil
	}
	finishDiceGame(gid)
	roll := c.Message().Dice.Value
	if roll > game.roll {
		if err := a.model.Economy.Update(gid, uid, int(game.amount)*2); err != nil {
			return err
		}
		return c.Send(fmt.Sprintf("💥 Вы выиграли `%v 💰`", formatAmount(int(game.amount)*2)), tele.ModeMarkdownV2)
	}
	if roll == game.roll {
		if err := a.model.Economy.Update(gid, uid, int(game.amount)); err != nil {
			return err
		}
		return c.Send("Ничья.")
	}
	return c.Send("Вы проиграли.")
}

func currentDiceGame(gid int64) *diceGame {
	game, ok := diceGames.Load(gid)
	if !ok {
		return nil
	}
	return game.(*diceGame)
}

func (a *App) newDiceGame(gid, uid int64, amount uint, roll int) error {
	id := time.Now()
	_, loaded := diceGames.LoadOrStore(gid, &diceGame{id, uid, amount, roll})
	if loaded {
		return errors.New("game already going")
	}
	time.AfterFunc(time.Second*secondsForRoll, func() { a.deleteGame(id, gid) })
	return nil
}

func finishDiceGame(gid int64) {
	diceGames.Delete(gid)
}

func (a *App) deleteGame(id time.Time, gid int64) error {
	group, err := a.bot.ChatByID(gid)
	if err != nil {
		return err
	}
	value, loaded := diceGames.Load(gid)
	if loaded {
		game := value.(*diceGame)
		if id == game.id {
			finishDiceGame(gid)
			_, err := a.bot.Send(group, fmt.Sprintf("Время вышло: вы потеряли `%s`\\.",
				formatAmount(int(game.amount))),
				tele.ModeMarkdownV2)
			return err
		}
	}
	return nil
}
